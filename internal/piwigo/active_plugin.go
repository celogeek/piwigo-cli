package piwigo

import (
	"encoding/json"
	"strings"
)

type ActivePlugin map[string]bool

func (uft *ActivePlugin) UnmarshalJSON(data []byte) error {
	var r []struct {
		Id    string `json:"id"`
		State string `json:"state"`
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	*uft = ActivePlugin{}
	for _, plugin := range r {
		if plugin.State == "active" {
			(*uft)[plugin.Id] = true
		}
	}

	return nil
}

func (uft ActivePlugin) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(uft))
	for k := range uft {
		keys = append(keys, k)
	}
	return json.Marshal(keys)
}

func (uft ActivePlugin) String() string {
	keys := make([]string, 0, len(uft))
	for k := range uft {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}
