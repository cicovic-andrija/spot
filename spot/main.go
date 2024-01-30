package spot

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cicovic-andrija/spot/config"
)

var (
	cfg config.Config
)

func readconfig() {
	var (
		conffile string
		err      error
	)

	flag.StringVar(&conffile, "config", "", "config file")
	flag.Parse()
	if conffile == "" {
		fmt.Fprintln(os.Stderr, "Error: Config file missing")
		os.Exit(1)
	}

	cfg, err = config.ReadConfig(conffile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

func init() {
	// set seed value for RDG
	rand.Seed(time.Now().UTC().UnixNano())
}

func systemsetup() *server {
	readconfig()
	return &server{
		addr: fmt.Sprintf("%s:%d", cfg.DevAddr, cfg.DevPort),
	}
}

// Run initialies and runs the service
func Run() {
	srvr := systemsetup()
	srvr.run()
}
