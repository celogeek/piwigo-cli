/*
	Debug tools

	debug.Dump(myStruct)
*/
package debug

import (
	"encoding/json"
)

/*
	Dump an interface to the stdout
*/
func Dump(v interface{}) string {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(result)
}
