package datastore

const (
	XerxesService = "XERXES:SERVICES:"
	XerxesNodes   = "XERXES:NODES:"
	XerxesPubSub  = "XERXES:*"
)

var (
	XerxesEvents = map[string]func() string{
		"reloadAll": func() string {
			return "XERXES:reloadAll"
		},
	}
)
