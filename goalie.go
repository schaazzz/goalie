package main

import (
    "fmt"
    "flag"
    "io/ioutil"
)



func main() {
    config := flag.String("config", "", "Config file for the selected mode")
    mode := flag.String("mode", "echo", "Mode of operation: input, proxy, pipe or echo - [default = echo]")
    join := make(chan bool)

    flag.Parse()
    fmt.Println(*mode, *config)

    configJSON, _ := ioutil.ReadFile(*config)

    if *mode == "pipe" {
        pipe := &Pipe{}
        pipe.Init(parsePipeConfigJSON(configJSON)[0:])
        go pipe.Start(join)
    }

    gone := 0
    for gone < 1 {
        select {
        case <- join:
            gone++
        }
    }
}
