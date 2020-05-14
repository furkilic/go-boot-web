package gobootweb

import (
	"github.com/furkilic/go-boot-config/pkg/go-boot-config"
	"os"
	"reflect"
	"testing"
)

func Test_retrieveConf(t *testing.T) {
	tests := []struct {
		name    string
		cmdArgs []string
		want    GoWebConf
	}{
		{"None", []string{},
			GoWebConf{
				":8080",
				0,
				"",
				1 << 20,
				60000,
				15000,
				15000,
				15000,
				Compression{false},
				NotFoundHandler{true},
				HTTP2{false},
				SSL{false, "", ""},
			},
		},
		{"With Port", []string{"--server.port=9090"},
			GoWebConf{
				":9090",
				9090,
				"",
				1 << 20,
				60000,
				15000,
				15000,
				15000,
				Compression{false},
				NotFoundHandler{true},
				HTTP2{false},
				SSL{false, "", ""},
			},
		},
		{"With All Fix Values",
			[]string{
				"--server.address=my-server:9090", "--server.base-path=/test", "--server.max-http-header-size=2097152",
				"--server.idle-timeout=6000", "--server.read-timeout=1400", "--server.write-timeout=1500", "--server.shutdown-timeout=1300",
				"--server.compression.enabled", "--server.not-found-handler.enabled=false", "--server.http2.enabled",
				"--server.ssl.enabled",
			},
			GoWebConf{
				"my-server:9090",
				0,
				"/test",
				2 << 20,
				6000,
				1500,
				1400,
				1300,
				Compression{true},
				NotFoundHandler{false},
				HTTP2{true},
				SSL{false, "", ""},
			},
		},
		{"With SSL", []string{"--server.ssl.cert-file=../../test/cert.pem", "--server.ssl.key-file=../../test/key.pem"},
			GoWebConf{
				":8080",
				0,
				"",
				1 << 20,
				60000,
				15000,
				15000,
				15000,
				Compression{false},
				NotFoundHandler{true},
				HTTP2{false},
				SSL{true, "../../test/cert.pem", "../../test/key.pem"},
			},
		},
	}
	for _, tt := range tests {
		os.Args = append([]string{"cmd"}, tt.cmdArgs...)
		gobootconfig.Reload()
		goWebConf = GoWebConf{}
		t.Run(tt.name, func(t *testing.T) {
			retrieveConf()
			if !reflect.DeepEqual(goWebConf, tt.want) {
				t.Errorf("retrieveConf() = %v, want %v", goWebConf, tt.want)
			}
		})
	}
	os.Args = []string{}
	gobootconfig.Reload()
}
