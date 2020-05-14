package gobootweb

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	startWebHandler()
	callback := make(chan *mux.Router)
	_getRouterCh <- callback
	return <-callback
}

func Start() chan error {
	startWebHandler()
	callback := make(chan chan error)
	_startServerCh <- callback
	return <-callback
}

func Stop() error {
	startWebHandler()
	callback := make(chan error)
	_stopServerCh <- callback
	return <-callback
}
