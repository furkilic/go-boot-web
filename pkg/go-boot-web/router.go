package gobootweb

import (
	"github.com/gorilla/mux"
)

var router = mux.NewRouter().StrictSlash(true)

func _getRouter() *mux.Router {
	r := router.Name("Root")
	if goWebConf.BasePath != "" {
		r = r.PathPrefix(goWebConf.BasePath)
	}
	return r.Subrouter()
}
