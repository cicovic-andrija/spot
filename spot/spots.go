package spot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cicovic-andrija/spot/api"
	"github.com/gorilla/mux"
)

// Params struct represents action parameters
type Params struct {
	Number int    `json:"number"`
	Label  string `json:"label"`
	Taken  bool   `json:"taken"`
}

// ActionMsg struct represents POST request data
type ActionMsg struct {
	Action string   `json:"action"`
	Params []Params `json:"params"`
}

func (s *server) httpSpots(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	garageID := urlVars["garage-id"]
	sectionName := urlVars["section-name"]

	switch r.Method {
	case http.MethodPost:
		s.postAction(w, r, garageID, sectionName)
	default:
		errMsg := fmt.Sprintf("invalid request for '%s'", api.Actions)
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
	}
}

func (s *server) postAction(w http.ResponseWriter, r *http.Request, garageID string, sectionName string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpInternalError(w, r, err)
		return
	}

	actionMsg := &ActionMsg{}
	err = json.Unmarshal(body, actionMsg)
	if err != nil {
		errMsg := "failed to unmarshal JSON object: " + err.Error()
		httpErrorResp(w, r, http.StatusBadRequest, errMsg)
		return
	}

	switch actionMsg.Action {
	case api.ActionUpdate:
		err = s.garages.actionUpdate(garageID, sectionName, actionMsg.Params)
	case api.ActionDisconnect:
		err = s.garages.actionDisconnect(garageID, sectionName, actionMsg.Params)
	default:
		errMsg := fmt.Sprintf("action '%s' not supported", actionMsg.Action)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
	}

	if err != nil {
		httpErrorResp(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
