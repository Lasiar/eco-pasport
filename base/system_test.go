package base

import (
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetConfig(t *testing.T) {
	_onceConfig.Do(func() {})
	// create and fill stub config
	config := new(Config)
	config.Port = "foo"
	config.PathTreeXML = "/test"
	// assign stub config to singleton
	_config = config
	tests := []struct {
		name string
		want *Config
	}{
		{name: "With stub config", want: config},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_load(t *testing.T) {
	type args struct {
		r io.Reader
	}
	stubConfJSON := `{"connect_string": "foobar","port":":80","xml_file":"Tree.xml"}`
	want := new(Config)
	want.Port = ":80"
	want.PathTreeXML = "Tree.xml"
	want.ConnStr = "foobar"
	want.Err = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	want.Warn = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime)
	want.Info = log.New(os.Stderr, "[INFO] ", log.Ldate|log.Ltime)
	tests := []struct {
		name string
		c    *Config
		args args
		want *Config
	}{
		{name: "With stub config", c: new(Config), args: args{strings.NewReader(stubConfJSON)}, want: want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.load(tt.args.r)
		})
		if !reflect.DeepEqual(tt.c, tt.want) {
			t.Errorf("load() = %v, want %v", tt.c, tt.want)
		}
	}
}
