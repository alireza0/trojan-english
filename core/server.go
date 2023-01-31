package core

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
	"io/ioutil"
)

var configPath = "/usr/local/etc/trojan/config.json"

// ServerConfig structure
type ServerConfig struct {
	Config
	SSl   ServerSSL `json:"ssl"`
	Tcp   ServerTCP `json:"tcp"`
	Mysql Mysql     `json:"mysql"`
}

// ServerSSL structure
type ServerSSL struct {
	SSL
	Key                string `json:"key"`
	KeyPassword        string `json:"key_password"`
	PreferServerCipher bool   `json:"prefer_server_cipher"`
	SessionTimeout     int    `json:"session_timeout"`
	PlainHttpResponse  string `json:"plain_http_response"`
	Dhparam            string `json:"dhparam"`
}

// ServerTCP structure
type ServerTCP struct {
	TCP
	PreferIPv4 bool `json:"prefer_ipv4"`
}

// Load Load the server configuration File
func Load(path string) []byte {
	if path == "" {
		path = configPath
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

// Save Save the server configuration File
func Save(data []byte, path string) bool {
	if path == "" {
		path = configPath
	}
	if err := ioutil.WriteFile(path, pretty.Pretty(data), 0644); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetConfig Get config configuration
func GetConfig() *ServerConfig {
	data := Load("")
	config := ServerConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		return nil
	}
	return &config
}

// GetMysql Get mysql connection
func GetMysql() *Mysql {
	return &GetConfig().Mysql
}

// WriteMysql Write mysql configuration
func WriteMysql(mysql *Mysql) bool {
	mysql.Enabled = true
	data := Load("")
	result, _ := sjson.SetBytes(data, "mysql", mysql)
	return Save(result, "")
}

// WriteTls Write TLS configuration
func WriteTls(cert, key, domain string) bool {
	data := Load("")
	data, _ = sjson.SetBytes(data, "ssl.cert", cert)
	data, _ = sjson.SetBytes(data, "ssl.key", key)
	data, _ = sjson.SetBytes(data, "ssl.sni", domain)
	return Save(data, "")
}

// WriteDomain Writing domain name
func WriteDomain(domain string) bool {
	data := Load("")
	data, _ = sjson.SetBytes(data, "ssl.sni", domain)
	return Save(data, "")
}

// WritePassword Writing Password
func WritePassword(pass []string) bool {
	data := Load("")
	data, _ = sjson.SetBytes(data, "password", pass)
	return Save(data, "")
}

// WritePort Write the Trojan port
func WritePort(port int) bool {
	data := Load("")
	data, _ = sjson.SetBytes(data, "local_port", port)
	return Save(data, "")
}

// WriteLogLevel Write a log level
func WriteLogLevel(level int) bool {
	data := Load("")
	data, _ = sjson.SetBytes(data, "log_level", level)
	return Save(data, "")
}
