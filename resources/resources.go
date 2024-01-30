package resources

import "time"

type (
	// Geolocation is the estimation of the real-world geographic location
	Geolocation struct {
		Longitude float64 `bson:"longitude" json:"longitude"`
		Latitude  float64 `bson:"latitude" json:"latitude"`
	}

	// Garage represents a garage resource
	Garage struct {
		ID          string      `bson:"id" json:"id"`
		Name        string      `bson:"name" json:"name"`
		City        string      `bson:"city" json:"city"`
		Address     string      `bson:"address" json:"address"`
		Geolocation Geolocation `bson:"geolocation" json:"geolocation"`
		Sections    []Section   `bson:"sections" json:"sections"`
	}

	// GarageRespObj is a JSON response object representing a garage
	GarageRespObj struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		City        string      `json:"city"`
		Address     string      `json:"address"`
		Geolocation Geolocation `json:"geolocation"`
		FreeSpots   int         `json:"free_spots"`
	}

	// Section represents a garage section resource
	Section struct {
		Name        string `bson:"name" json:"name"`
		Level       string `bson:"level" json:"level"`
		Description string `bson:"description" json:"description"`
		TotalSpots  int    `bson:"total_spots" json:"total_spots"`
		FreeSpots   int
		Spots       []Spot
	}

	// SectionRespObj is a JSON response object representing a section
	SectionRespObj struct {
		Name        string `json:"name"`
		Level       string `json:"level"`
		Description string `json:"description"`
		TotalSpots  int    `json:"total_spots"`
		FreeSpots   int    `json:"free_spots"`
	}

	// Spot represents a parking spot
	Spot struct {
		Label      string
		Online     bool
		Taken      bool
		LastUpdate time.Time
	}
)
