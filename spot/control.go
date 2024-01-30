package spot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *server) httpControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpErrorResp(w, r, http.StatusBadRequest, "invalid request")
		return
	}

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
	case "shutdown":
		s.shutdown()
		w.WriteHeader(http.StatusOK)
	default:
		errMsg := fmt.Sprintf("action '%s' not supported", actionMsg.Action)
		httpErrorResp(w, r, http.StatusNotFound, errMsg)
	}

}
