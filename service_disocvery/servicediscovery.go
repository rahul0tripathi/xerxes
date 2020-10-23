package main

import (
	"github.com/gookit/color"
	"github.com/rahultripathidev/docker-utility/service_disocvery/health"
	"net/http"
	"time"
)

func startSchedulers() {
	go health.Scheduler()
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/reload", reloadHandler)
	startSchedulers()
	color.Style{color.FgCyan, color.OpBold}.Printf("[%s] Server Started \n",time.Now().String())
	err := http.ListenAndServe(":8937", nil)
	if err != nil {
		panic(err)
	}
}
func reloadHandler(w http.ResponseWriter, r *http.Request) {
	go health.StopTheWholeWorldAndReload()
	w.WriteHeader(201)
	w.Write([]byte{})
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
