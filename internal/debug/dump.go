/*
	Debug tools

	debug.Dump(myStruct)
*/
package debug

import (
	"encoding/json"
	"os"
)

/*
	Dump an interface to the stdout
*/
func Dump(v interface{}) error {
	d := json.NewEncoder(os.Stdout)
	d.SetIndent("", "  ")
	return d.Encode(v)
}
