# go-boot-web 

go-boot-web let you easily start and configure a golang web server using gorilla/mux

Internally using go-boot-conf it let's you configure the server through properties files, YAML files, environment variables, and command-line arguments 

The idea behind go-boot-web is to standartized web configuration in Go


## Simple Usage 


To start your web server you only need to 
```go
gobootweb.Start()
defer gobootweb.Stop()
```

To Add a new Router to your web server you only need to 
```go
router := gobootweb.Router()
router.Methods("GET").Path("/hello").Name("Hello").HandlerFunc(hello)
```

Full Example
```go
package main

import (
	"encoding/json"
	"github.com/furkilic/go-boot-web/pkg/go-boot-web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errCh := gobootweb.Start()
	defer gobootweb.Stop()

	router := gobootweb.Router()
	router.Methods("GET").Path("/hello").Name("Hello").HandlerFunc(hello)

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		log.Println(err)
		break
	case err := <-termChan:
		log.Println(err)
		break
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Hello World !!")
}
```
To run your app
```sh
./my-app
```


## Server Properties 

Those properties can be set using properties files, YAML files, environment variables, and command-line arguments. See go-boot-conf for more information

Property Name | Default Value | Comment | Example Value
--- | --- | --- | ---
`server.address`| `:8080` | Address of the server | `my-server:8888` 
`server.port`| `8080` | Port of the server (useless if `server.address` is set)| `9090` 
`server.base-path`|  | Base Path of your Router | `/my-app`
`server.max-http-header-size` | `1048576` |  Max HTTP Header size (in byte) | `2097152`
`server.idle-timeout`| `60000` | Idle Timeout (in ms) | `1000` 
`server.write-timeout`| `15000` | Write Timeout (in ms) | `1000` 
`server.read-timeout`| `15000` | Read Timeout (in ms)| `1000` 
`server.shutdown-timeout`| `15000` | Server Shutdown Timeout (in ms)| `1000` 
`server.not-found-handler.enabled`| `true` | go-boot-web custom not found handler| `false`
`server.compression.enabled`| `false` | Enable Compression if request sends `Accept-Encoding: gzip`| `true` 
`server.http2.enabled`| `false` | Enable HTTP2 | `true` 
`server.ssl.enabled`| `true` if Cert And Key are provided else `false`| Enable SSL | `false` 
`server.ssl.cert-file`|  | SSL Certificate File Location| `assets/cert.pem` 
`server.ssl.key-file`|  | SSL Key File Location| `assets/key.pem` 
