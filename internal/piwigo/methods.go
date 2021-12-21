package piwigo

import (
	"encoding/json"
)

type Methods []string

type MethodParams []MethodParam

type MethodParam struct {
	Name         string      `json:"name"`
	Optional     bool        `json:"optional"`
	Type         string      `json:"type"`
	AcceptArray  bool        `json:"acceptArray"`
	DefaultValue interface{} `json:"defaultValue"`
	MaxValue     interface{} `json:"maxValue"`
	Info         string      `json:"info"`
}

type MethodOptions struct {
	Admin    bool `json:"admin_only"`
	PostOnly bool `json:"post_only"`
}

func (j *MethodOptions) UnmarshalJSON(data []byte) error {
	var r interface{}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	switch r := r.(type) {
	case map[string]interface{}:
		j.Admin, _ = r["admin_only"].(bool)
		j.PostOnly, _ = r["post_only"].(bool)
	}
	return nil
}

type MethodDetails struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Options     MethodOptions `json:"options"`
	Parameters  MethodParams  `json:"params"`
}
