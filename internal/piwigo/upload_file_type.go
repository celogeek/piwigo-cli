package piwigo

import (
	"encoding/json"
	"strings"
)

type UploadFileType map[string]bool

func (uft *UploadFileType) UnmarshalJSON(data []byte) error {
	var r string
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*uft = UploadFileType{}
	for _, v := range strings.Split(r, ",") {
		(*uft)[v] = true
	}
	return nil
}

func (uft UploadFileType) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(uft))
	for k := range uft {
		keys = append(keys, k)
	}
	return json.Marshal(keys)
}

func (uft UploadFileType) String() string {
	keys := make([]string, 0, len(uft))
	for k := range uft {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}
