package main

type MethodGroup struct {
	List    MethodListCommand    `command:"list" description:"List of available methods"`
	Details MethodDetailsCommand `command:"details" description:"Details of a method"`
	Try     MethodTryCommand     `command:"try" description:"Test a method. Parameters after the command as k=v, can be repeated like k=v1 k=v2."`
}

var methodGroup MethodGroup

func init() {
	_, err := parser.AddCommand("method", "Reflexion management", "", &methodGroup)
	if err != nil {
		panic(err)
	}
}
