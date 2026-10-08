// Package config ...
package config

import (
	"encoding/json/v2"
	"net/netip"
)

type Config struct {
	Bind       netip.Addr `json:"bind"`
	Port       Port       `json:"port"`
	Replicaof  string     `json:"replicaof"`
	Debug      bool       `json:"-"`
	Dir        string     `json:"dir"`
	DBFilename string     `json:"dbfilename"`

	AppendOnly     string `json:"appendonly"`
	AppendDirName  string `json:"appenddirname"`
	AppendFileName string `json:"appendfilename"`
	AppendFSync    string `json:"appendfsync"`
}

func (c *Config) ToMap() (map[string]string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
