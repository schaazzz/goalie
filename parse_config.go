package main

import (
    _ "fmt"
    "encoding/json"
    common "github.com/schaazzz/golibs/common"
)

type PipeJSON struct {
    Role    string      `json:role`
    Address string      `json:address`
    RedirectDelay int   `json:redirectDelay`
}

func parsePipeConfigJSON(jsonData []byte) []PipeJSON {
    var config [2]PipeJSON
    err := json.Unmarshal([]byte(jsonData), &config)

    if err != nil ||
        config[0].Role == config[1].Role ||
        !common.CheckAgainst(config[0].Role, "server", "client") ||
        !common.CheckAgainst(config[1].Role, "server", "client") {
            panic("There was an error while trying to parse the configuration file!")
        }
    return config[0:]
}