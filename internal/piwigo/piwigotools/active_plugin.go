package piwigotools

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

func (uft ActivePlugin) keys() []string {
	keys := make([]string, 0, len(uft))
	for k := range uft {
		keys = append(keys, k)
	}
	return keys
}

func (uft ActivePlugin) MarshalJSON() ([]byte, error) {
	return json.Marshal(uft.keys())
}

func (uft ActivePlugin) String() string {
	return strings.Join(uft.keys(), ",")
}
