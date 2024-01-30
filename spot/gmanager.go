package spot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cicovic-andrija/spot/db"
	"github.com/cicovic-andrija/spot/log"
	"github.com/cicovic-andrija/spot/resources"
	"github.com/cicovic-andrija/spot/util"
)

type garageManager struct {
	db      *db.Client
	rw      *sync.RWMutex
	garages map[string]*resources.Garage
}

func newGarageManager(db *db.Client) (*garageManager, error) {
	var err error

	gm := &garageManager{
		db: db,
		rw: &sync.RWMutex{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	gm.garages, err = db.FindAllGarages(ctx)
	if err != nil {
		return nil, err
	}

	for _, g := range gm.garages {
		for i := range g.Sections {
			g.Sections[i].Spots = make([]resources.Spot, g.Sections[i].TotalSpots)
		}
	}

	return gm, nil
}

func (m *garageManager) uniqueID() string {
	m.rw.RLock()
	for {
		id, err := util.NewRandomID()
		if err != nil {
			log.Errorf("Failed to obtain a garage ID: %v", err)
			continue
		}

		if _, exists := m.garages[id]; !exists {
			m.rw.RUnlock()
			return id
		}

		log.Infof("Garage ID collision prevented. ID: '%s'", id)
	}
}

func (m *garageManager) getGarages() []resources.GarageRespObj {
	m.rw.RLock()
	respArray := []resources.GarageRespObj{}
	for _, v := range m.garages {
		respObj := resources.GarageRespObj{
			ID:          v.ID,
			Name:        v.Name,
			City:        v.City,
			Address:     v.Address,
			Geolocation: v.Geolocation,
		}
		for _, s := range v.Sections {
			respObj.FreeSpots += s.FreeSpots
		}
		respArray = append(respArray, respObj)
	}
	m.rw.RUnlock()
	return respArray
}

func (m *garageManager) getGarage(id string) (g resources.GarageRespObj, ok bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	garage, ok := m.garages[id]
	if !ok {
		return
	}
	g.ID = garage.ID
	g.Name = garage.Name
	g.City = garage.City
	g.Address = garage.Address
	g.Geolocation = garage.Geolocation
	for _, s := range garage.Sections {
		g.FreeSpots += s.FreeSpots
	}
	return
}

func (m *garageManager) addGarage(garage *resources.Garage) error {
	garage.ID = m.uniqueID()
	garage.Sections = []resources.Section{}

	m.rw.Lock()
	defer m.rw.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := m.db.InsertGarage(ctx, garage); err != nil {
		return err
	}
	m.garages[garage.ID] = garage
	return nil
}

func (m *garageManager) updateGarage(id string, update *resources.Garage) (found bool, respObj resources.GarageRespObj, err error) {
	// Note: Changing geolocation is not enabled

	m.rw.Lock()
	defer m.rw.Unlock()

	garage, found := m.garages[id]
	if !found {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = m.db.UpdateGarage(ctx, id, update.Name, update.City, update.Address); err != nil {
		return
	}

	if update.Name != "" {
		garage.Name = update.Name
	}
	if update.City != "" {
		garage.City = update.City
	}
	if update.Address != "" {
		garage.Address = update.Address
	}

	respObj.ID = id
	respObj.Name = garage.Name
	respObj.City = garage.City
	respObj.Address = garage.Address
	respObj.Geolocation = garage.Geolocation
	return
}

func (m *garageManager) removeGarage(id string) (found bool, err error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	if _, found = m.garages[id]; !found {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = m.db.DeleteGarage(ctx, id); err != nil {
		return
	}
	delete(m.garages, id)
	return
}

func (m *garageManager) sectionExists(garageID string, sectionName string) (bool, *resources.Garage, int) {
	// NOTE: This function is *not* thread-safe
	garage, ok := m.garages[garageID]
	if !ok {
		return false, nil, -1
	}
	for i, s := range garage.Sections {
		if s.Name == sectionName {
			return true, garage, i
		}
	}
	return false, nil, -1
}

func (m *garageManager) getSections(garageID string) (respArray []resources.SectionRespObj, found bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	garage, found := m.garages[garageID]
	if !found {
		return
	}

	respArray = []resources.SectionRespObj{}
	for _, s := range garage.Sections {
		respArray = append(
			respArray,
			resources.SectionRespObj{
				Name:        s.Name,
				Level:       s.Level,
				Description: s.Description,
				TotalSpots:  s.TotalSpots,
				FreeSpots:   s.FreeSpots,
			},
		)
	}
	return
}

func (m *garageManager) getSection(garageID string, sectionName string) (respObj resources.SectionRespObj, found bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	found, garage, i := m.sectionExists(garageID, sectionName)
	if !found {
		return
	}

	respObj.Name = garage.Sections[i].Name
	respObj.Level = garage.Sections[i].Level
	respObj.Description = garage.Sections[i].Description
	respObj.TotalSpots = garage.Sections[i].TotalSpots
	respObj.FreeSpots = garage.Sections[i].FreeSpots
	return
}

func (m *garageManager) addSection(garageID string, section *resources.Section) (garageFound bool, sectionExists bool, err error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	garage, garageFound := m.garages[garageID]
	if !garageFound {
		return
	}

	for _, s := range garage.Sections {
		if s.Name == section.Name {
			sectionExists = true
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = m.db.InsertSection(ctx, garageID, section); err != nil {
		return
	}

	section.Spots = make([]resources.Spot, section.TotalSpots)
	garage.Sections = append(garage.Sections, *section)
	return
}

func (m *garageManager) updateSection(garageID, sectionName string, update *resources.Section) (found bool, respObj resources.SectionRespObj, err error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	found, garage, i := m.sectionExists(garageID, sectionName)
	if !found {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = m.db.UpdateSection(
		ctx,
		garageID,
		sectionName,
		update.Name,
		update.Level,
		update.Description,
		update.TotalSpots,
	)
	if err != nil {
		return
	}

	section := &garage.Sections[i]
	if update.Name != "" {
		section.Name = update.Name
	}
	if update.Level != "" {
		section.Level = update.Level
	}
	if update.Description != "" {
		section.Description = update.Description
	}
	if update.TotalSpots > 0 {
		section.Spots = make([]resources.Spot, update.TotalSpots)
		section.TotalSpots = update.TotalSpots
		section.FreeSpots = 0
	}

	respObj.Name = section.Name
	respObj.Level = section.Level
	respObj.Description = section.Description
	respObj.TotalSpots = section.TotalSpots
	respObj.FreeSpots = section.FreeSpots
	return
}

func (m *garageManager) deleteSection(garageID string, sectionName string) (found bool, err error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	found, garage, i := m.sectionExists(garageID, sectionName)
	if !found {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = m.db.DeleteSection(ctx, garageID, sectionName)
	if err != nil {
		return
	}
	garage.Sections = append(garage.Sections[:i], garage.Sections[i+1:]...)
	return
}

func (m *garageManager) actionUpdate(garageID string, sectionName string, params []Params) error {
	var err error

	m.rw.Lock()
	defer m.rw.Unlock()

	exists, garage, i := m.sectionExists(garageID, sectionName)
	if !exists {
		return fmt.Errorf("section '%s', garage id %s not found", sectionName, garageID)
	}

	for _, param := range params {
		if param.Number < 1 || param.Number > garage.Sections[i].TotalSpots {
			if err == nil {
				err = fmt.Errorf(
					"%d is not a valid spot number for section '%s', garage '%s' (garage id %s)",
					param.Number,
					sectionName,
					garage.Name,
					garageID,
				)
			} else {
				err = fmt.Errorf(
					"%v\n%d is not a valid spot number for section '%s', garage '%s' (garage id %s)",
					err,
					param.Number,
					sectionName,
					garage.Name,
					garageID,
				)
			}
			continue
		}

		spot := &garage.Sections[i].Spots[param.Number-1]
		if !param.Taken {
			garage.Sections[i].FreeSpots++
		} else if spot.Online {
			garage.Sections[i].FreeSpots--
		}

		if param.Label != "" {
			spot.Label = param.Label
		}

		spot.Online = true
		spot.Taken = param.Taken
		spot.LastUpdate = time.Now()

		log.Infof(
			"Update: garage: '%s' (garage id %s); section: '%s'; spot #%d (label '%s'); taken: %v",
			garage.Name,
			garageID,
			sectionName,
			param.Number,
			spot.Label,
			param.Taken,
		)
	}

	return err
}

func (m *garageManager) actionDisconnect(garageID string, sectionName string, params []Params) error {
	var err error

	m.rw.Lock()
	defer m.rw.Unlock()

	exists, garage, i := m.sectionExists(garageID, sectionName)
	if !exists {
		return fmt.Errorf("section '%s', garage id %s not found", sectionName, garageID)
	}

	for _, param := range params {

		if param.Number < 1 || param.Number > garage.Sections[i].TotalSpots {
			if err == nil {
				err = fmt.Errorf(
					"%d is not a valid spot number for section '%s', garage '%s' (garage id %s)",
					param.Number,
					sectionName,
					garage.Name,
					garageID,
				)
			} else {
				err = fmt.Errorf(
					"%v\n%d is not a valid spot number for section '%s', garage '%s' (garage id %s)",
					err,
					param.Number,
					sectionName,
					garage.Name,
					garageID,
				)
			}
			continue
		}

		if spot := &garage.Sections[i].Spots[param.Number-1]; spot.Online {
			if !spot.Taken {
				garage.Sections[i].FreeSpots--
			}
			spot.Label, spot.Taken, spot.Online, spot.LastUpdate = "", false, false, time.Time{}

			log.Infof(
				"Update: garage: '%s' (garage id %s); section: '%s'; spot #%d (label '%s') disconnected",
				garage.Name,
				garageID,
				sectionName,
				param.Number,
				spot.Label,
			)
		}
	}

	return err
}

func (m *garageManager) invalidateOldUpdates() {
	m.rw.Lock()
	defer m.rw.Unlock()

	now := time.Now()

	for _, g := range m.garages {
		for i := range g.Sections {
			for j := range g.Sections[i].Spots {
				spot := &g.Sections[i].Spots[j]
				if spot.Online {
					if now.Sub(g.Sections[i].Spots[j].LastUpdate) > 20*time.Minute {
						if !spot.Taken {
							g.Sections[i].FreeSpots--
						}
						spot.Label = ""
						spot.Taken = false
						spot.Online = false
						spot.LastUpdate = time.Time{}
					}
				}
			}
		}
	}
}
