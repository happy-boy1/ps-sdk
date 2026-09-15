package solarman

import (
	"ps-sdk/model"
	"time"
)

func (s *StationListItem) PowerStation() *model.PowerStation {
	gridConnectedAt := time.Unix(int64(s.StartOperatingTime), 0)
	var ps = model.PowerStation{
		PlatformID:      1,
		StationIDOrigin: s.ID.String(),
		StationName:     s.Name,
		Address:         s.LocationAddress,
		Longitude:       s.LocationLng.Float(),
		Latitude:        s.LocationLat.Float(),
		GridConnectedAt: &gridConnectedAt,
	}
	return &ps
}
