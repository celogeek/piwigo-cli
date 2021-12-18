package piwigocli

type MethodGroup struct {
	List    MethodListCommand    `command:"list" description:"List of available methods"`
	Details MethodDetailsCommand `command:"details" description:"Details of a method"`
	Try     MethodTryCommand     `command:"try" description:"Test a method"`
}

var methodGroup MethodGroup

func init() {
	parser.AddCommand("method", "Reflexion management", "", &methodGroup)
}
