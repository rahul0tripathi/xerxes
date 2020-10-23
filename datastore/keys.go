package datastore

const (
	XerxesPubSub  = "XERXES:*"
)

var (
	XerxesEvents = map[string]func() string{
		"reloadAll": func() string {
			return "XERXES:reloadAll"
		},
	}
)
