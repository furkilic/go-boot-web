package gobootweb

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/http2"
)

func Test_web(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	type req struct {
		url     string
		headers map[string]string
		https   bool
		http2   bool
	}
	type resp struct {
		proto    string
		status   int
		response string
		err      bool
	}
	defaultWebConf := GoWebConf{
		":0000",
		0,
		"",
		1 << 20,
		60000,
		15000,
		15000,
		5,
		Compression{false},
		NotFoundHandler{true},
		HTTP2{false},
		SSL{false, "", ""},
	}
	httpsConf := GoWebConf{
		":0000",
		0,
		"",
		1 << 20,
		60000,
		15000,
		15000,
		5,
		Compression{false},
		NotFoundHandler{true},
		HTTP2{true},
		SSL{true, "../../test/cert.pem", "../../test/key.pem"},
	}
	errorConf := GoWebConf{
		":0000",
		0,
		"",
		16,
		100,
		100,
		100,
		5,
		Compression{false},
		NotFoundHandler{true},
		HTTP2{false},
		SSL{false, "", ""},
	}
	tests := []struct {
		name    string
		conf    GoWebConf
		path    string
		handler func(w http.ResponseWriter, r *http.Request)
		req     req
		resp    resp
	}{
		{"None",
			defaultWebConf, "none",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("NoNo")
			},
			req{"http://localhost:0000/none", createMap(), false, false},
			resp{
				"HTTP/1.1",
				http.StatusOK,
				"\"NoNo\"",
				false,
			},
		},
		{"With BasePath",
			GoWebConf{
				":0000",
				0,
				"/base-path",
				1 << 20,
				60000,
				15000,
				15000,
				5,
				Compression{false},
				NotFoundHandler{true},
				HTTP2{false},
				SSL{false, "", ""},
			}, "basePath",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("BasePath")
			},
			req{"http://localhost:0000/base-path/basePath", createMap(), false, false},
			resp{
				"HTTP/1.1",
				http.StatusOK,
				"\"BasePath\"",
				false,
			},
		},
		{"With Compression",
			GoWebConf{
				":0000",
				0,
				"",
				1 << 20,
				60000,
				15000,
				15000,
				5,
				Compression{true},
				NotFoundHandler{true},
				HTTP2{false},
				SSL{false, "", ""},
			}, "compressed",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("Compressed")
			},
			req{"http://localhost:0000/compressed", createMap("Accept-Encoding", "gzip"), false, false},
			resp{
				"HTTP/1.1",
				http.StatusOK,
				"\"Compressed\"",
				false,
			},
		},
		{"NotFound Default Handler",
			defaultWebConf, "none",
			nil,
			req{"http://localhost:0000/not-found", createMap(), false, false},
			resp{
				"HTTP/1.1",
				http.StatusNotFound,
				"{\"status\":404,\"method\":\"GET\",\"uri\":\"/not-found\",\"message\":\"Not Found\"}",
				false,
			},
		},
		{"NotFound JSON Handler",
			defaultWebConf, "none",
			nil,
			req{"http://localhost:0000/not-found", createMap("Accept", "application/json"), false, false},
			resp{
				"HTTP/1.1",
				http.StatusNotFound,
				"{\"status\":404,\"method\":\"GET\",\"uri\":\"/not-found\",\"message\":\"Not Found\"}",
				false,
			},
		},
		{"NotFound XML Handler",
			defaultWebConf, "none",
			nil,
			req{"http://localhost:0000/not-found", createMap("Accept", "application/xml"), false, false},
			resp{
				"HTTP/1.1",
				http.StatusNotFound,
				"<customError><status>404</status><method>GET</method><uri>/not-found</uri><message>Not Found</message></customError>",
				false,
			},
		},
		{"NotFound HTML Handler",
			defaultWebConf, "none",
			nil,
			req{"http://localhost:0000/not-found", createMap("Accept", "text/html"), false, false},
			resp{
				"HTTP/1.1",
				http.StatusNotFound,
				"<html><head><title>404 - Not Found</title></head><body><h3>404 - Not Found</h3><p>Method: GET<br/>URI: <a href=\"/not-found\">/not-found</a></p></body></html>",
				false,
			},
		},
		{"With HTTPS",
			httpsConf, "https",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("HTTPS")
			},
			req{"https://localhost:0000/https", createMap(), true, false},
			resp{
				"HTTP/1.1",
				http.StatusOK,
				"\"HTTPS\"",
				false,
			},
		},
		{"With HTTP2",
			httpsConf, "http2",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("HTTP2")
			},
			req{"https://localhost:0000/http2", createMap(), true, true},
			resp{
				"HTTP/2.0",
				http.StatusOK,
				"\"HTTP2\"",
				false,
			},
		},
		{"Error Big Header",
			errorConf, "error",
			func(w http.ResponseWriter, r *http.Request) { time.Sleep(time.Second * 1) },
			req{"http://localhost:0000/error",
				func() map[string]string {
					m := make(map[string]string)
					for i := 0; i < 1000; i++ {
						m[fmt.Sprintf("HEADER-KEY-%d", i)] = fmt.Sprintf("HEADER-VALUE-%d", i)
					}
					return m
				}(),
				false, false},
			resp{
				"HTTP/1.1",
				http.StatusRequestHeaderFieldsTooLarge,
				"431 Request Header Fields Too Large",
				false,
			},
		},
		{"Error Write Timeout",
			errorConf, "error",
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Second * 1)
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("toto")
			},
			req{"http://localhost:0000/error", createMap(),
				false, false},
			resp{
				"",
				0,
				"",
				true,
			},
		},
	}

	for _, tt := range tests {
		port := rand.Intn(1010) + 8080
		router = mux.NewRouter().StrictSlash(true)
		goWebConf = tt.conf
		goWebConf.Address = fmt.Sprintf(":%d", port)
		tt.req.url = strings.ReplaceAll(tt.req.url, ":0000", goWebConf.Address)
		addDefaultValues()
		Start()
		time.Sleep(time.Millisecond * 100)
		Router().Methods("GET").Path("/" + tt.path).Name(strings.ToUpper(tt.path)).HandlerFunc(tt.handler)
		t.Run(tt.name, func(t *testing.T) {
			proto, status, response, err := call(tt.req.url, tt.req.headers, tt.req.https, tt.req.http2)
			if err != nil && !tt.resp.err {
				t.Errorf("web() expected no error but got : %v", err.Error())
				return
			} else if err == nil && tt.resp.err {
				t.Errorf("web() expected error but don't got")
				return
			} else if err != nil && tt.resp.err {
				return
			}
			r := resp{proto, status, response, false}

			if !reflect.DeepEqual(r, tt.resp) {
				t.Errorf("web() = %v, want %v", r, tt.resp)
			}
		})
		Stop()
	}
}

func call(url string, headers map[string]string, isHttps, isHttp2 bool) (string, int, string, error) {
	client := &http.Client{}
	if isHttps {
		caCert, err := ioutil.ReadFile("../../test/cert.pem")
		if err != nil {
			log.Fatalf("Reading server certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		if isHttp2 {
			client.Transport = &http2.Transport{
				TLSClientConfig: tlsConfig,
			}
		} else {
			client.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		}
	}
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, "", err
	}
	defer resp.Body.Close()
	r := make([]byte, 0)
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		r = append(r, scanner.Bytes()...)
	}
	reader, err := gzip.NewReader(bytes.NewReader(r))
	if err == nil {
		all, _ := ioutil.ReadAll(reader)
		r = all
	}
	return resp.Proto, resp.StatusCode, strings.TrimSpace(string(r)), nil
}

func createMap(keyValues ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i < len(keyValues); i += 2 {
		m[keyValues[i]] = keyValues[i+1]
	}
	return m
}
