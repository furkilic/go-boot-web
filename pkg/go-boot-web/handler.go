package gobootweb

import (
	"github.com/gorilla/mux"
	"sync"
)

var _getRouterCh = make(chan chan *mux.Router)
var _startServerCh = make(chan chan chan error)
var _stopServerCh = make(chan chan error)

var once sync.Once

func startWebHandler() {
	once.Do(func() {
		retrieveConf()
		go _startWebHandler()
	})
}

func _startWebHandler() {
	for {
		select {
		case callback := <-_getRouterCh:
			callback <- _getRouter()
			break
		case callback := <-_startServerCh:
			callback <- _startServer()
			break
		case callback := <-_stopServerCh:
			callback <- _stopServer()
			break
		}
	}
}
