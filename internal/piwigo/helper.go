package piwigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func DumpResponse(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func ArgsToForm(args []string) (*url.Values, error) {
	params := &url.Values{}
	for _, arg := range args {
		r := strings.SplitN(arg, "=", 2)
		if len(r) != 2 {
			return nil, errors.New("args should be key=value")
		}
		params.Add(r[0], r[1])
	}
	return params, nil
}
