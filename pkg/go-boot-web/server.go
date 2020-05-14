package gobootweb

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/handlers"
	"net/http"
	"time"
)

var srv *http.Server

func _startServer() chan error {
	errorChan := make(chan error)
	tlsConfig := &tls.Config{}
	if goWebConf.HTTP2.Enabled {
		tlsConfig.NextProtos = []string{"h2"}
	} else {
		tlsConfig.NextProtos = []string{"http/1.1"}
	}
	var s http.Handler
	if goWebConf.Compression.Enabled {
		s = handlers.CompressHandler(router)
	} else {
		s = router
	}
	srv = &http.Server{
		Addr:           goWebConf.Address,
		WriteTimeout:   time.Millisecond * time.Duration(goWebConf.WriteTimeout),
		ReadTimeout:    time.Millisecond * time.Duration(goWebConf.ReadTimeout),
		IdleTimeout:    time.Millisecond * time.Duration(goWebConf.IdleTimeout),
		Handler:        s,
		TLSConfig:      tlsConfig,
		MaxHeaderBytes: goWebConf.MaxHttpHeaderSize,
	}
	go func() {
		if goWebConf.SSL.Enabled {
			if err := srv.ListenAndServeTLS(goWebConf.SSL.CertFile, goWebConf.SSL.KeyFile); err != nil {
				errorChan <- err
			}
		} else {
			if err := srv.ListenAndServe(); err != nil {
				errorChan <- err
			}
		}
	}()
	return errorChan
}

func _stopServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(goWebConf.ShutdownTimeout))
	defer cancel()
	return srv.Shutdown(ctx)
}
