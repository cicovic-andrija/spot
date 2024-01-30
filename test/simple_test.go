package spot_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"github.com/cicovic-andrija/spot/resources"
)

const (
	testBaseURL           = "http://localhost:8000/"
	testGarageName        = "TestGarage"
	testSectionName       = "TestSection"
	testSectionTotalSpots = 10
)

func CreateGarage(client *http.Client, garageName string, expectedStatus int) (*resources.GarageRespObj, error) {
	garage := struct {
		Name string `json:"name"`
	}{
		Name: garageName,
	}

	reqBody, err := json.Marshal(garage)
	if err != nil {
		return nil, err
	}

	url := testBaseURL + path.Join("v1", "garages")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Unexpected POST status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := &resources.GarageRespObj{}
	err = json.Unmarshal(respBody, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func GetGarage(client *http.Client, garageID string, expectedStatus int) (*resources.GarageRespObj, error) {
	url := testBaseURL + path.Join("v1", "garages", garageID)
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Unexpected GET status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := &resources.GarageRespObj{}
	err = json.Unmarshal(respBody, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func DeleteGarage(client *http.Client, garageID string, expectedStatus int) error {
	url := testBaseURL + path.Join("v1", "garages", garageID)
	req, err := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("Unexpected DELETE status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	return nil
}

func CreateSection(client *http.Client, garageID string, sectionName string, totalSpots int, expectedStatus int) (*resources.SectionRespObj, error) {
	section := struct {
		Name       string `json:"name"`
		TotalSpots int    `json:"total_spots"`
	}{
		Name:       sectionName,
		TotalSpots: totalSpots,
	}

	reqBody, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	url := testBaseURL + path.Join("v1", "garages", garageID, "sections")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Unexpected POST status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := &resources.SectionRespObj{}
	err = json.Unmarshal(respBody, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func GetSections(client *http.Client, garageID string, expectedStatus int) ([]resources.SectionRespObj, error) {
	url := testBaseURL + path.Join("v1", "garages", garageID, "sections")
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Unexpected GET status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respArray := make([]resources.SectionRespObj, 0)
	err = json.Unmarshal(respBody, &respArray)
	if err != nil {
		return nil, err
	}

	return respArray, nil
}

func GetSection(client *http.Client, garageID string, sectionName string, expectedStatus int) (*resources.SectionRespObj, error) {
	url := testBaseURL + path.Join("v1", "garages", garageID, "sections", sectionName)
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Unexpected GET status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := &resources.SectionRespObj{}
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func DeleteSection(client *http.Client, garageID string, sectionName string, expectedStatus int) error {
	url := testBaseURL + path.Join("v1", "garages", garageID, "sections", sectionName)
	req, err := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("Unexpected DELETE status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	return nil
}

func UpdateStatus(client *http.Client, garageID string, sectionName string, spotNumber int, isTaken bool, expectedStatus int) error {
	type param struct {
		Number int  `json:"number"`
		Taken  bool `json:"taken"`
	}

	actionMsg := struct {
		Action string  `json:"action"`
		Params []param `json:"params"`
	}{
		Action: "update",
		Params: []param{param{Number: spotNumber, Taken: isTaken}},
	}

	reqBody, err := json.Marshal(actionMsg)
	if err != nil {
		return err
	}

	url := testBaseURL + path.Join("v1", "garages", garageID, "sections", sectionName, "actions")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("Unexpected POST status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	return nil
}

func Disconnect(client *http.Client, garageID string, sectionName string, spotNumber int, expectedStatus int) error {
	type param struct {
		Number int `json:"number"`
	}

	actionMsg := struct {
		Action string  `json:"action"`
		Params []param `json:"params"`
	}{
		Action: "disconnect",
		Params: []param{param{Number: spotNumber}},
	}

	reqBody, err := json.Marshal(actionMsg)
	if err != nil {
		return err
	}

	url := testBaseURL + path.Join("v1", "garages", garageID, "sections", sectionName, "actions")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("Unexpected POST status: %d. Expected: %d", resp.StatusCode, expectedStatus)
	}

	return nil
}

func TestCreateGarage(t *testing.T) {
	c := &http.Client{}

	postRespObj, err := CreateGarage(c, testGarageName, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	_, err = GetGarage(c, postRespObj.ID, http.StatusOK)
	if err != nil {
		t.Error(err)
	}

	err = DeleteGarage(c, postRespObj.ID, http.StatusNoContent)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateSection(t *testing.T) {
	c := &http.Client{}

	garageRespObj, err := CreateGarage(c, testGarageName, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	_, err = CreateSection(c, garageRespObj.ID, testSectionName, testSectionTotalSpots, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	sections, err := GetSections(c, garageRespObj.ID, http.StatusOK)
	if err != nil {
		t.Error(err)
	}

	if len := len(sections); len != 1 {
		t.Errorf("Just one section expected. Found %d", len)
	} else if sections[0].Name != testSectionName {
		t.Errorf("Unexpected section name: %s. Expected: %s", sections[0].Name, testSectionName)
	} else if sections[0].TotalSpots != testSectionTotalSpots {
		t.Errorf("Unexpected section's total spot number: %d. Expected: %d", sections[0].TotalSpots, testSectionTotalSpots)
	} else if sections[0].FreeSpots != 0 {
		t.Errorf("Unexpected section's free spot number: %d. Expected: 0", sections[0].FreeSpots)
	}

	err = DeleteSection(c, garageRespObj.ID, testSectionName, http.StatusNoContent)
	if err != nil {
		t.Error(err)
	}

	err = DeleteGarage(c, garageRespObj.ID, http.StatusNoContent)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateSpots(t *testing.T) {
	c := &http.Client{}

	garageRespObj, err := CreateGarage(c, testGarageName, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	_, err = CreateSection(c, garageRespObj.ID, testSectionName, testSectionTotalSpots, http.StatusCreated)
	if err != nil {
		t.Error(err)
	}

	t.Log("Updating all spots to status FREE")
	for i := 1; i <= testSectionTotalSpots; i++ {
		err = UpdateStatus(c, garageRespObj.ID, testSectionName, i, false, http.StatusOK)
		if err != nil {
			t.Error(err)
		}
	}

	sectionRespObj, err := GetSection(c, garageRespObj.ID, testSectionName, http.StatusOK)
	if err != nil {
		t.Error(err)
	} else if sectionRespObj.FreeSpots != testSectionTotalSpots {
		t.Errorf("Unexpected section's free spot number: %d. Expected: %d", sectionRespObj.FreeSpots, testSectionTotalSpots)
	} else {
		t.Logf("Expected free spot number: %d. Found: %d", testSectionTotalSpots, sectionRespObj.FreeSpots)
	}

	t.Log("Disconnecting all spots")
	for i := 1; i <= testSectionTotalSpots; i++ {
		err = Disconnect(c, garageRespObj.ID, testSectionName, i, http.StatusOK)
		if err != nil {
			t.Error(err)
		}
	}

	sectionRespObj, err = GetSection(c, garageRespObj.ID, testSectionName, http.StatusOK)
	if err != nil {
		t.Error(err)
	} else if sectionRespObj.FreeSpots != 0 {
		t.Errorf("Unexpected section's free spot number: %d. Expected: 0", sectionRespObj.FreeSpots)
	} else {
		t.Logf("Expected free spot number: 0. Found: %d", sectionRespObj.FreeSpots)
	}

	err = DeleteSection(c, garageRespObj.ID, testSectionName, http.StatusNoContent)
	if err != nil {
		t.Error(err)
	}

	err = DeleteGarage(c, garageRespObj.ID, http.StatusNoContent)
	if err != nil {
		t.Error(err)
	}
}
