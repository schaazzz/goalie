package main

import (
    "fmt"
    "os"
    "time"
    "log"
    "strconv"
    tcp "github.com/schaazzz/golibs/network/tcp"
)

var logger * log.Logger

func handleConnection(c * tcp.Connection, join chan bool) {
    reset: go c.Start()

    forever: for {
        select {
        case <- c.Panic:
            logger.Println(c.Name, "panicked, resetting in 3 seconds!")
            time.Sleep(3 * time.Second)
            goto reset
        case dataIn := <- c.DataIn:
            logger.Println(fmt.Sprintf("- %s >>", c.Name), string(dataIn.Bytes))
        case serverConnectionState := <- c.Connected:
            if serverConnectionState {
                c.Ctrl <- "start"
                c.DataOut <- &tcp.DataChunk{Length: len("oogabooga"), Bytes: []byte("oogabooga")}
            } else {
                break forever
            }
        case <- c.Done:
            break forever
        }
    }

    join <- true
}

func main() {
    logger = log.New(os.Stdout, "[MAIN] ", log.Lmicroseconds)
    server := &tcp.Connection {
                Channels: tcp.Channels {
                    Ctrl        : make(chan string),
                    DataIn      : make(chan * tcp.DataChunk),
                    DataOut     : make(chan * tcp.DataChunk),
                    Done        : make(chan bool),
                    Connected   : make(chan bool),
                    Panic       : make(chan bool),
                },
                Server  : true,
                Address : strconv.Itoa(17231),
                Name    : "SERVER 0",
            }

    client := &tcp.Connection {
                Channels: tcp.Channels {
                    Ctrl        : make(chan string),
                    DataIn      : make(chan * tcp.DataChunk),
                    DataOut     : make(chan * tcp.DataChunk),
                    Done        : make(chan bool),
                    Connected   : make(chan bool),
                    Panic       : make(chan bool),
                },
                Server  : false,
                Address : "127.0.0.1:17232",
                Name    : "CLIENT 0",
            }

    join := make(chan bool)

    go handleConnection(server, join)
    go handleConnection(client, join)

    gone := 0
    for gone < 2 {
        select {
        case <- join:
            gone++
        }
    }
}
