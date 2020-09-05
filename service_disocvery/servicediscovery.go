package main

import (
	"github.com/rahultripathidev/docker-utility/service_disocvery/health"
	"net/http"
)

func startSchedulers() {
	go health.Scheduler()
	go health.Reloader()
}
func main() {
	http.HandleFunc("/", handler)
	startSchedulers()
	err := http.ListenAndServe(":8937", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var _body []byte
	_body = health.GetTarget(r.URL.Path)
	if len(_body) > 0 {
		w.WriteHeader(200)
		w.Write(_body)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("NOT FOUND"))
	}

}
