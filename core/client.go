package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"trojan/asset"
)

// ClientConfig structure
type ClientConfig struct {
	Config
	SSl ClientSSL `json:"ssl"`
	Tcp ClientTCP `json:"tcp"`
}

// ClientSSL structure
type ClientSSL struct {
	SSL
	Verify         bool `json:"verify"`
	VerifyHostname bool `json:"verify_hostname"`
}

// ClientTCP structure
type ClientTCP struct {
	TCP
}

// WriteClient generates client JSON
func WriteClient(port int, password, domain, writePath string) bool {
	data := asset.GetAsset("client.json")
	config := ClientConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		return false
	}
	config.RemoteAddr = domain
	config.RemotePort = port
	config.Password = []string{password}
	outData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println(err)
		return false
	}
	if err = ioutil.WriteFile(writePath, outData, 0644); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
