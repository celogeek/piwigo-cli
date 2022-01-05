/*
	Debug tools

	fmt.Println(debug.Dump(myStruct))
*/
package debug

import (
	"encoding/json"
)

/*
	Dump an interface into a json string format
*/
func Dump(v interface{}) string {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(result)
}
