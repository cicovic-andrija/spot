package spot_test

import (
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	virtualGarageName              = "VirtualGarage"
	virtualGarageSectionName       = "VirtualSection"
	virtualGarageSectionTotalSpots = 100
	virtualGarageCarCount          = 20
)

type (
	ReleaseMsg struct {
		ID     int
		Number int
	}

	Car struct {
		ID         int
		RequestCh  chan int
		ResponseCh chan int
		ReleaseCh  chan ReleaseMsg
		AckCh      chan struct{}
	}
)

func VirtualCar(sem chan struct{}, garageID string, car *Car, t *testing.T) {
	var err error

	// enter garage
	sem <- struct{}{}

	car.RequestCh <- car.ID
	number := <-car.ResponseCh

	c := &http.Client{}

	// look for parking
	time.Sleep(time.Duration(rand.Intn(10)+10) * time.Second)

	// park
	err = UpdateStatus(c, garageID, virtualGarageSectionName, number, true, http.StatusOK)
	if err != nil {
		t.Error(err)
	}

	// go shopping
	time.Sleep(time.Duration(rand.Intn(5)+15) * time.Second)

	// unpark
	err = UpdateStatus(c, garageID, virtualGarageSectionName, number, false, http.StatusOK)
	if err != nil {
		t.Error(err)
	}

	car.ReleaseCh <- ReleaseMsg{ID: car.ID, Number: number}
	<-car.AckCh

	// leave garage
	<-sem
}

func RunVirtualGarage(garageID string, t *testing.T) {
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, virtualGarageSectionTotalSpots)
	requestsCh := make(chan int)
	releaseCh := make(chan ReleaseMsg)
	section := make([]bool, virtualGarageSectionTotalSpots) // true == taken
	carsServed := 0

	// intialize spots to free
	c := &http.Client{}
	for i := 1; i <= virtualGarageSectionTotalSpots; i++ {
		err := UpdateStatus(c, garageID, virtualGarageSectionName, i, false, http.StatusOK)
		if err != nil {
			t.Error(err)
		}
	}

	// create cars
	cars := make([]*Car, virtualGarageCarCount)
	for i := range cars {
		cars[i] = &Car{
			ID:         i,
			RequestCh:  requestsCh,
			ResponseCh: make(chan int),
			ReleaseCh:  releaseCh,
			AckCh:      make(chan struct{}),
		}
	}

	// run cars
	for _, c := range cars {
		wg.Add(1)
		go func(car *Car) {
			VirtualCar(sem, garageID, car, t)
			wg.Done()
		}(c)
	}

	quit := false
	for !quit {
		select {
		case carId := <-requestsCh:
			for i, taken := range section {
				if !taken {
					section[i] = true
					cars[carId].ResponseCh <- i + 1
					break
				}
			}
		case releaseMsg := <-releaseCh:
			section[releaseMsg.Number-1] = false
			cars[releaseMsg.ID].AckCh <- struct{}{}
			carsServed++
			if carsServed == virtualGarageCarCount {
				quit = true
			}
		}
	}

	// wait for cars to finish
	wg.Wait()
}

func TestVirtualGarage(t *testing.T) {
	c := &http.Client{}

	garageRespObj, err := CreateGarage(c, virtualGarageName, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	_, err = CreateSection(c, garageRespObj.ID, virtualGarageSectionName, virtualGarageSectionTotalSpots, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	RunVirtualGarage(garageRespObj.ID, t)

	err = DeleteGarage(c, garageRespObj.ID, http.StatusNoContent)
	if err != nil {
		t.Error(err)
	}
}
