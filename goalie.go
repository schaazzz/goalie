package main
import (
    "fmt"
    "flag"
    "log"
    "time"
    "io/ioutil"
    tcp "github.com/schaazzz/golibs/network/tcp"
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
    } else if *mode == "echo" {
        echo := &Echo{}
        echo.Init(parseEchoConfigJSON(configJSON)[0:])
        go echo.Start(join)
    }

    gone := 0
    for gone < 1 {
        select {
        case <- join:
            gone++
        }
    }
}

func handleConnection(c * tcp.Connection, logger * log.Logger, join chan bool) {
    reset: go c.Start()

    forever: for {
        select {
        case <- c.Panic:
            logger.Println(c.Name, "panicked, resetting in 3 seconds!")
            time.Sleep(3 * time.Second)
            goto reset
        case serverConnectionState := <- c.Connected:
            if serverConnectionState {
                c.Ctrl <- "start"
            } else {
                break forever
            }
        case <- c.Done:
            break forever
        }
    }

    join <- true
}