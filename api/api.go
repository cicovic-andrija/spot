package api

import "path"

// API constans
const (
	V1                 = "v1"
	CollectionGarages  = "garages"
	ObjectGarage       = "{garage-id:" + patternID + "}"
	CollectionSections = "sections"
	ObjectSection      = "{section-name:" + patternSectionName + "}"
	Control            = "control"

	Actions          = "actions"
	ActionUpdate     = "update"
	ActionDisconnect = "disconnect"

	patternID          = `[0-9a-f]{8}`
	patternSectionName = `[0-9a-zA-Z]+`
)

// Path returns an API path constructed from given elements
func Path(elem ...string) string {
	return "/" + path.Join(elem...)
}
