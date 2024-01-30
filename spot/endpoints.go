package spot

import (
	"net/http"

	"github.com/cicovic-andrija/spot/api"
	"github.com/gorilla/mux"
)

func (s *server) setupEndpoints() http.Handler {
	s.router = mux.NewRouter()

	s.router.HandleFunc(
		api.Path(api.V1, api.CollectionGarages),
		s.httpGarages,
	)

	s.router.HandleFunc(
		api.Path(api.V1, api.CollectionGarages, api.ObjectGarage),
		s.httpGarage,
	)

	s.router.HandleFunc(
		api.Path(api.V1, api.CollectionGarages, api.ObjectGarage,
			api.CollectionSections),
		s.httpSections,
	)

	s.router.HandleFunc(
		api.Path(api.V1, api.CollectionGarages, api.ObjectGarage,
			api.CollectionSections, api.ObjectSection),
		s.httpSection,
	)

	s.router.HandleFunc(
		api.Path(api.V1, api.CollectionGarages, api.ObjectGarage,
			api.CollectionSections, api.ObjectSection, api.Actions),
		s.httpSpots,
	)

	s.router.HandleFunc(
		api.Path(api.V1, api.Control),
		s.httpControl,
	)

	return s.router
}
