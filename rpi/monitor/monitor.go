package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type (
	// State of the parking spot
	State int

	// Mode of device execution
	Mode int

	actionMsgParams struct {
		Number int    `json:"number"`
		Label  string `json:"label"`
		Taken  bool   `json:"taken"`
	}

	actionMsg struct {
		Action string            `json:"action"`
		Params []actionMsgParams `json:"params"`
	}
)

const (
	bcmPinLEDRed        = 25
	bcmPinLEDGreen      = 12
	bcmPinLEDBlue       = 16
	bcmPinHCSR04Trigger = 17
	bcmPinHCSR04Echo    = 4
	sTaken              = State(0)
	sFree               = State(1)
	mNormal             = Mode(0)
	mPanic              = Mode(1)
)

var (
	redLED            rpio.Pin
	greenLED          rpio.Pin
	blueLED           rpio.Pin
	hcsr04Trigger     rpio.Pin
	hcsr04Echo        rpio.Pin
	maxNoupdateIntrvl int32
	panicCh           = make(chan struct{})
	nopanicCh         = make(chan struct{})
	httpReqCh         = make(chan State, 16)
)

func initGPIO() {
	redLED = rpio.Pin(bcmPinLEDRed)
	redLED.Output()

	greenLED = rpio.Pin(bcmPinLEDGreen)
	greenLED.Output()

	blueLED = rpio.Pin(bcmPinLEDBlue)
	blueLED.Output()

	hcsr04Trigger = rpio.Pin(bcmPinHCSR04Trigger)
	hcsr04Trigger.Output()

	hcsr04Echo = rpio.Pin(bcmPinHCSR04Echo)
	hcsr04Echo.Input()

	hcsr04Trigger.Low()
	time.Sleep(2 * time.Second)
}

func cmdistance() float64 {
	var pulseStart, pulseEnd time.Time
	hcsr04Trigger.High()
	time.Sleep(100 * time.Microsecond)
	hcsr04Trigger.Low()
	time.Now()
	for hcsr04Echo.Read() == rpio.Low {
		pulseStart = time.Now()
	}
	for hcsr04Echo.Read() == rpio.High {
		pulseEnd = time.Now()
	}
	duration := pulseEnd.Sub(pulseStart)
	return duration.Seconds() * 17150.0
}

func modeController(quit chan struct{}, heartbeatIntrvl, discoveryIntrvl int) {
	blueLED.High()
	mode := mNormal
	atomic.StoreInt32(&maxNoupdateIntrvl, int32(heartbeatIntrvl))
	wg := sync.WaitGroup{}
	ch := make(chan struct{})
	for {
		select {
		case <-panicCh:
			if mode == mPanic {
				continue
			}
			wg.Add(1)
			go func() {
				for {
					select {
					case <-ch:
						wg.Done()
						return
					case <-time.After(time.Second / 2):
						blueLED.Toggle()
					}
				}
			}()
			mode = mPanic
			atomic.StoreInt32(&maxNoupdateIntrvl, int32(discoveryIntrvl))
		case <-nopanicCh:
			if mode == mNormal {
				continue
			}
			ch <- struct{}{}
			wg.Wait()
			blueLED.High()
			mode = mNormal
			atomic.StoreInt32(&maxNoupdateIntrvl, int32(heartbeatIntrvl))
		case <-quit:
			if mode == mPanic {
				ch <- struct{}{}
				wg.Wait()
			}
			blueLED.Low()
			return
		}
	}
}

func httpRunner(quit chan struct{}, url string, number int, label string) {
	httpclient := &http.Client{}
	actionMsg := actionMsg{Action: "update", Params: []actionMsgParams{actionMsgParams{Number: number, Label: label}}}

	for {
		select {
		case state := <-httpReqCh:
			actionMsg.Params[0].Taken = state == sTaken
			reqBody, err := json.Marshal(actionMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				continue
			}
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				continue
			}
			resp, err := httpclient.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				panicCh <- struct{}{}
				continue
			}
			resp.Body.Close()
			nopanicCh <- struct{}{}
		case <-quit:
			return
		}
	}
}

func stateController(quit chan struct{}, maxdist float64, stateUpdateIntrvl int) {
	greenLED.High()
	state := sFree
	var newstate State
	sinceLastChange := 0
	sinceLastUpdate := 0
	httpReqCh <- state

	for {
		select {
		case <-time.After(time.Second):
			if dist := cmdistance(); dist > maxdist {
				newstate = sFree
				redLED.Low()
				greenLED.High()
			} else {
				newstate = sTaken
				redLED.High()
				greenLED.Low()
			}

			sinceLastUpdate++
			if state != newstate {
				sinceLastChange++
				if sinceLastChange == stateUpdateIntrvl {
					state = newstate
					sinceLastChange = 0
					sinceLastUpdate = 0
					httpReqCh <- state
					continue
				}
			} else {
				sinceLastChange = 0
			}

			if int32(sinceLastUpdate) >= atomic.LoadInt32(&maxNoupdateIntrvl) {
				sinceLastUpdate = 0
				httpReqCh <- state
			}
		case <-quit:
			redLED.Low()
			greenLED.Low()
			return
		}
	}
}

func disconnect(url string, number int) {
	httpclient := &http.Client{}
	actionMsg := actionMsg{Action: "disconnect", Params: []actionMsgParams{actionMsgParams{Number: number}}}
	reqBody, err := json.Marshal(actionMsg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to disconnect: %v\n", err.Error())
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to disconnect: %v\n", err.Error())
		return
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to disconnect: %v\n", err.Error())
		return
	}
	resp.Body.Close()
}

func main() {
	var (
		url               string
		label             string
		maxdist           float64
		number            int
		heartbeatIntrvl   int
		discoveryIntrvl   int
		stateUpdateIntrvl int
		wg                sync.WaitGroup
	)

	flag.StringVar(&url, "url", "http://localhost:8000/", "Spot service API endpoint")
	flag.IntVar(&number, "number", 0, "Parking spot number")
	flag.StringVar(&label, "label", "", "Parking spot label")
	flag.Float64Var(&maxdist, "maxdist", 200.0, "Maximal valid distance [cm]")
	flag.IntVar(&heartbeatIntrvl, "heartbeatintrvl", 900, "Heartbeat interval [s]")
	flag.IntVar(&discoveryIntrvl, "discoveryintrvl", 30, "Server discovery interval [s]")
	flag.IntVar(&stateUpdateIntrvl, "stateupdateintrvl", 5, "Change state interval [s]")
	flag.Parse()

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	if maxdist < 2.0 || maxdist > 400.0 {
		maxdist = 200.0
		fmt.Fprintf(os.Stderr, "maxdist not in range [2.0, 400.0], defaulting to 200.0cm")
	}

	if heartbeatIntrvl < 1 {
		heartbeatIntrvl = 900
		fmt.Fprintf(os.Stderr, "invalid heartbeatintrvl value, defaulting to 900s")
	}

	if discoveryIntrvl < 1 {
		discoveryIntrvl = 30
		fmt.Fprintf(os.Stderr, "invalid discoveryintrvl value, defaulting to 30s")
	}

	if stateUpdateIntrvl < 1 {
		stateUpdateIntrvl = 5
		fmt.Fprintf(os.Stderr, "invalid stateupdateintrvl value, defaulting to 5s")
	}

	err := rpio.Open()
	if err != nil {
		panic("failed to open GPIO: " + err.Error())
	}
	defer rpio.Close()
	initGPIO()

	wg.Add(1)
	quitM := make(chan struct{})
	go func() {
		modeController(quitM, heartbeatIntrvl, discoveryIntrvl)
		wg.Done()
	}()

	wg.Add(1)
	quitS := make(chan struct{})
	go func() {
		stateController(quitS, maxdist, stateUpdateIntrvl)
		wg.Done()
	}()

	wg.Add(1)
	quitH := make(chan struct{})
	go func() {
		httpRunner(quitH, url, number, label)
		wg.Done()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	quitM <- struct{}{}
	quitS <- struct{}{}
	quitH <- struct{}{}
	wg.Wait()

	disconnect(url, number)
}
