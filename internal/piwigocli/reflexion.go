package piwigocli

type ReflexionGroup struct {
	MethodList    ReflexionMethodListCommand    `command:"method-list" description:"List of available methods"`
	MethodDetails ReflexionMethodDetailsCommand `command:"method-details" description:"Details of a method"`
}

var reflexionGroup ReflexionGroup

func init() {
	parser.AddCommand("reflexion", "Reflexion management", "", &reflexionGroup)
}
