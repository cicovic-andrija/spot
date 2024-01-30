package spot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cicovic-andrija/spot/api"
	"github.com/cicovic-andrija/spot/resources"
	"github.com/gorilla/mux"
)

func (s *server) httpGarages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getGarages(w, r)
	case http.MethodPost:
		s.postGarages(w, r)
	default:
		errMsg := fmt.Sprintf("invalid request for resource '%s'", api.CollectionGarages)
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
	}
}

func (s *server) httpGarage(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	id := urlVars["garage-id"]
	switch r.Method {
	case http.MethodGet:
		s.getGarage(w, r, id)
	case http.MethodPut, http.MethodPatch:
		s.putGarage(w, r, id)
	case http.MethodDelete:
		s.deleteGarage(w, r, id)
	default:
		errMsg := fmt.Sprintf("invalid request for resource '%s/%s'", api.CollectionGarages, id)
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
	}
}

func (s *server) getGarages(w http.ResponseWriter, r *http.Request) {
	respArray := s.garages.getGarages()
	resp, err := json.Marshal(respArray)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(resp)
}

func (s *server) postGarages(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}

	garage := &resources.Garage{}
	err = json.Unmarshal(body, garage)
	if err != nil {
		errMsg := "failed to unmarshal JSON object: " + err.Error()
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	err = s.garages.addGarage(garage)
	if err != nil {
		err = errors.New("DB error: failed to insert garage: " + err.Error())
		httpInternalError(w, r, err)
		return
	}

	resp, err := json.Marshal(
		resources.GarageRespObj{
			ID:          garage.ID,
			Name:        garage.Name,
			City:        garage.City,
			Address:     garage.Address,
			Geolocation: garage.Geolocation,
		},
	)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (s *server) getGarage(w http.ResponseWriter, r *http.Request, id string) {
	respObj, found := s.garages.getGarage(id)
	if !found {
		errMsg := fmt.Sprintf("resource '%s/%s' not found", api.CollectionGarages, id)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *server) putGarage(w http.ResponseWriter, r *http.Request, id string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}

	update := &resources.Garage{}
	err = json.Unmarshal(body, update)
	if err != nil {
		errMsg := "failed to unmarshal JSON object: " + err.Error()
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	found, respObj, err := s.garages.updateGarage(id, update)
	if !found {
		errMsg := fmt.Sprintf("resource '%s/%s' not found", api.CollectionGarages, id)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}
	if err != nil {
		err = errors.New("DB error: failed to update garage: " + err.Error())
		httpInternalError(w, r, err)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *server) deleteGarage(w http.ResponseWriter, r *http.Request, id string) {
	found, err := s.garages.removeGarage(id)
	if !found {
		errMsg := fmt.Sprintf("resource '%s/%s' not found", api.CollectionGarages, id)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}
	if err != nil {
		err = errors.New("DB error: failed to delete garage: " + err.Error())
		httpInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
