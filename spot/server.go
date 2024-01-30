package spot

import (
	"context"
	"net/http"
	"time"

	"github.com/cicovic-andrija/spot/db"
	"github.com/cicovic-andrija/spot/log"

	"github.com/gorilla/mux"
)

type server struct {
	httpServer *http.Server
	router     *mux.Router
	addr       string
	garages    *garageManager
	runners    []backgroundRunner
}

func (s *server) startRunners(garages *garageManager) {
	s.runners = []backgroundRunner{&invalidationRunner{garages: garages}}
	for _, r := range s.runners {
		r.start()
	}
}

func (s *server) stopRunners() {
	for _, r := range s.runners {
		r.stop()
	}
}

func (s *server) run() {
	db, err := db.NewClient(cfg.DBConfig.ConnString, cfg.DBConfig.Database, cfg.DBConfig.Collection)
	if err != nil {
		log.Fatalf("DB: failed to connect to database: %s", err.Error())
	}

	s.garages, err = newGarageManager(db)
	if err != nil {
		log.Fatalf("DB: failed to get garages: %s", err.Error())
	}

	s.startRunners(s.garages)

	handler := s.setupEndpoints()
	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: handler,
	}
	err = s.httpServer.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatalf("Http server stopped unexpectedly")
		s.shutdown()
	} else {
		log.Infof("Http server stopped")
	}
}

func (s *server) shutdown() {
	if s.httpServer != nil {
		s.stopRunners()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			log.Fatalf("Failed to shutdown http server gracefully: %s", err.Error())
		} else {
			s.httpServer = nil
		}
	}
}

func httpErrorResp(w http.ResponseWriter, r *http.Request, status int, msg string) {
	errorMsg := http.StatusText(status) + ": " + r.Method + " " + r.URL.Path + ": " + msg
	http.Error(w, errorMsg, status)
}

func httpInternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Errorf("%s", r.Method+" "+r.URL.Path+": "+err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
