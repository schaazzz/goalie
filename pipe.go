package main

import (
    "fmt"
    "os"
    "time"
    "log"
    tcp "github.com/schaazzz/golibs/network/tcp"
)

type Pipe struct {
    ListenAddr  string
    ConnectAddr string
    Delay       int
    server      * tcp.Connection
    client      * tcp.Connection
    join        chan bool
}

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

func (p * Pipe) Init() {
    logger = log.New(os.Stdout, "[PIPE MAIN] ", log.Lmicroseconds)
    p.server = &tcp.Connection {
                Channels: tcp.Channels {
                    Ctrl        : make(chan string),
                    DataIn      : make(chan * tcp.DataChunk),
                    DataOut     : make(chan * tcp.DataChunk),
                    Done        : make(chan bool),
                    Connected   : make(chan bool),
                    Panic       : make(chan bool),
                },
                Server  : true,
                Address : p.ListenAddr,
                Name    : "PIPE SERVER",
            }

    p.client = &tcp.Connection {
                Channels: tcp.Channels {
                    Ctrl        : make(chan string),
                    DataIn      : make(chan * tcp.DataChunk),
                    DataOut     : make(chan * tcp.DataChunk),
                    Done        : make(chan bool),
                    Connected   : make(chan bool),
                    Panic       : make(chan bool),
                },
                Server  : false,
                Address : p.ConnectAddr,
                Name    : "PIPE CLIENT",
            }

}

func (p * Pipe) Start(join chan bool) {
    p.join = make(chan bool)
    go handleConnection(p.server, p.join)
    go handleConnection(p.client, p.join)

    gone := 0
    for gone < 2 {
        select {
        case <- join:
            gone++
        }
    }

    join <- true
}
