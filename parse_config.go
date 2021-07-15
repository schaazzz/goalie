package main

import (
	"encoding/json"
)

type Service struct {
	Name      string   `json:"name"`
	CmdList   []string `json:"cmd-list"`
	Path      string   `json:"path"`
	Autostart bool     `json:"autostart"`
}

type Config struct {
	Services []Service `json:"services"`

	TcpIpSockets []struct {
		Role    string `json:"role"`
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"tcp-ip"`
}

func parseConfigJSON(jsonData []byte) Config {
	var config Config
	_ = json.Unmarshal([]byte(jsonData), &config)

	// if err != nil ||
	// 	config[0].Role == config[1].Role ||
	// 	!common.CheckAgainst(config[0].Role, "server", "client") ||
	// 	!common.CheckAgainst(config[1].Role, "server", "client") {
	// 	panic("There was an error while trying to parse the configuration file!")
	// }
	return config
}

// func parseTcpIpSocketsJSON(jsonData []byte) []TcpIpSocket {
// 	var sockets []TcpIpSocket
// 	err := json.Unmarshal([]byte(jsonData), &sockets)

// 	fmt.Println("...", err, sockets)

// 	for _, element := range sockets {
// 		fmt.Println("---", element)
// 		if err != nil { //||
// 			//!common.CheckAgainst(element.Role, "server", "client") {
// 			//panic("There was an error while trying to parse the configuration file!")
// 			println("@@@ err:", err)
// 		}
// 	}
// 	return sockets[0:]
// }
