package main

type SessionGroup struct {
	Login  SessionLoginCommand  `command:"login" description:"Initialize a connection to a piwigo instance"`
	Status SessionStatusCommand `command:"status" description:"Get the status of your session"`
}

var sessionGroup SessionGroup

func init() {
	_, err := parser.AddCommand("session", "Session management", "", &sessionGroup)
	if err != nil {
		panic(err)
	}
}
