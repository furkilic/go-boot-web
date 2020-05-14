package gobootweb

import (
	"fmt"
	"github.com/furkilic/go-boot-config/pkg/go-boot-config"
	"net/http"
)

type GoWebConf struct {
	Address           string
	Port              int
	BasePath          string
	MaxHttpHeaderSize int
	IdleTimeout       int64
	WriteTimeout      int64
	ReadTimeout       int64
	ShutdownTimeout   int64
	Compression       Compression
	NotFoundHandler   NotFoundHandler
	HTTP2             HTTP2
	SSL               SSL
}
type Compression struct {
	Enabled bool
}
type HTTP2 struct {
	Enabled bool
}

type NotFoundHandler struct {
	Enabled bool
}

type SSL struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

var goWebConf GoWebConf

func retrieveConf() {
	gobootconfig.Load()
	gobootconfig.GetObject("server", &goWebConf)
	addDefaultValues()
}

func addDefaultValues() {
	if goWebConf.Address == "" {
		if goWebConf.Port == 0 {
			goWebConf.Address = ":8080"
		} else {
			goWebConf.Address = fmt.Sprintf(":%d", goWebConf.Port)
		}
	}
	if goWebConf.IdleTimeout == 0 {
		goWebConf.IdleTimeout = 60000
	}
	if goWebConf.ReadTimeout == 0 {
		goWebConf.ReadTimeout = 15000
	}
	if goWebConf.WriteTimeout == 0 {
		goWebConf.WriteTimeout = 15000
	}
	if goWebConf.ShutdownTimeout == 0 {
		goWebConf.ShutdownTimeout = 15000
	}
	sslEnabled := gobootconfig.GetBoolWithDefault("server.ssl.enabled", true)
	if sslEnabled {
		if goWebConf.SSL.CertFile == "" || goWebConf.SSL.KeyFile == "" {
			goWebConf.SSL.Enabled = false
		} else {
			goWebConf.SSL.Enabled = true
		}
	}
	if goWebConf.MaxHttpHeaderSize == 0 {
		goWebConf.MaxHttpHeaderSize = http.DefaultMaxHeaderBytes
	}
	if gobootconfig.GetBoolWithDefault("server.not-found-handler.enabled", true) {
		goWebConf.NotFoundHandler.Enabled = true
		router.NotFoundHandler = http.HandlerFunc(myCustomHandler)
	}
}
