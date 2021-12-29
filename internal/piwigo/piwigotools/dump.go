package piwigotools

import (
	"encoding/json"
	"fmt"
)

func DumpResponse(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
