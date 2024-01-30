package spot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/cicovic-andrija/spot/api"
	"github.com/cicovic-andrija/spot/resources"
	"github.com/gorilla/mux"
)

const (
	sectionNamePattern = `^[a-zA-Z0-9]+$`
)

var (
	nameRegex = regexp.MustCompile(sectionNamePattern)
)

func (s *server) httpSections(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	garageID := urlVars["garage-id"]

	switch r.Method {
	case http.MethodGet:
		s.getSections(w, r, garageID)
	case http.MethodPost:
		s.postSections(w, r, garageID)
	default:
		errMsg := fmt.Sprintf("invalid request for resource '%s'", api.CollectionSections)
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
	}
}

func (s *server) httpSection(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	garageID := urlVars["garage-id"]
	sectionName := urlVars["section-name"]

	switch r.Method {
	case http.MethodGet:
		s.getSection(w, r, garageID, sectionName)
	case http.MethodPut, http.MethodPatch:
		s.putSection(w, r, garageID, sectionName)
	case http.MethodDelete:
		s.deleteSection(w, r, garageID, sectionName)
	default:
		errMsg := fmt.Sprintf("invalid request for resource '%s/%s'", api.CollectionSections, sectionName)
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
	}
}

func (s *server) getSections(w http.ResponseWriter, r *http.Request, garageID string) {
	respArray, found := s.garages.getSections(garageID)
	if !found {
		errMsg := fmt.Sprintf("resource '%s/%s' not found", api.CollectionGarages, garageID)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}

	resp, err := json.Marshal(respArray)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(resp)
}

func (s *server) postSections(w http.ResponseWriter, r *http.Request, garageID string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}

	section := &resources.Section{}
	err = json.Unmarshal(body, section)
	if err != nil {
		errMsg := "failed to unmarshal JSON object: " + err.Error()
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	if !nameRegex.MatchString(section.Name) {
		errMsg := "section name in wrong format, use pattern: " + sectionNamePattern
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	if section.TotalSpots < 1 {
		errMsg := "illegal value for total number of spots, must be at least 1"
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	garageFound, sectionExists, err := s.garages.addSection(garageID, section)
	if !garageFound {
		errMsg := fmt.Sprintf("resource '%s/%s' not found", api.CollectionGarages, garageID)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}
	if sectionExists {
		errMsg := "section already exists"
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}
	if err != nil {
		err = errors.New("DB error: failed to create section: " + err.Error())
		httpInternalError(w, r, err)
		return
	}

	resp, err := json.Marshal(
		resources.SectionRespObj{
			Name:        section.Name,
			Level:       section.Level,
			Description: section.Description,
			TotalSpots:  section.TotalSpots,
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

func (s *server) getSection(w http.ResponseWriter, r *http.Request, garageID string, sectionName string) {
	respObj, found := s.garages.getSection(garageID, sectionName)
	if !found {
		errMsg := fmt.Sprintf(
			"resource '%s/%s/%s/%s' not found",
			api.CollectionGarages,
			garageID,
			api.CollectionSections,
			sectionName,
		)
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

func (s *server) putSection(w http.ResponseWriter, r *http.Request, garageID string, sectionName string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}

	update := &resources.Section{}
	err = json.Unmarshal(body, update)
	if err != nil {
		errMsg := "failed to unmarshal JSON object: " + err.Error()
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	if update.Name != "" && !nameRegex.MatchString(update.Name) {
		errMsg := "section name in wrong format, use pattern: " + sectionNamePattern
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	// 0 is ignored while updating section information (it is 0 if not provided in request)
	if update.TotalSpots < 0 {
		errMsg := "illegal value for total number of spots, must be at least 1"
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	found, respObj, err := s.garages.updateSection(garageID, sectionName, update)
	if !found {
		errMsg := fmt.Sprintf(
			"resource '%s/%s/%s/%s' not found",
			api.CollectionGarages,
			garageID,
			api.CollectionSections,
			sectionName,
		)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}
	if err != nil {
		err = errors.New("DB error: failed to update section: " + err.Error())
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

func (s *server) deleteSection(w http.ResponseWriter, r *http.Request, garageID string, sectionName string) {
	found, err := s.garages.deleteSection(garageID, sectionName)
	if !found {
		errMsg := fmt.Sprintf(
			"resource '%s/%s/%s/%s' not found",
			api.CollectionGarages,
			garageID,
			api.CollectionSections,
			sectionName,
		)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
		return
	}
	if err != nil {
		err = errors.New("DB error: failed to delete section: " + err.Error())
		httpInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
