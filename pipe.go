package main

import (
    "os"
    "time"
    "log"
    tcp "github.com/schaazzz/golibs/network/tcp"
)

type pipeConfig struct {
    address         string
    redirectDelay   int
}

type Pipe struct {
    serverConfig    * pipeConfig
    clientConfig    * pipeConfig
    server          * tcp.Connection
    client          * tcp.Connection
    join            chan bool
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

func (p * Pipe) Init(pipeJSON []PipeJSON) {

    for _, element := range pipeJSON {
        config := &pipeConfig {element.Address, element.RedirectDelay}
        if element.Role == "server" {
            p.serverConfig = config
        } else {
            p.clientConfig = config
        }
    }

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
                Address : p.serverConfig.address,
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
                Address : p.clientConfig.address,
                Name    : "PIPE CLIENT",
            }
}

func (p * Pipe) Start(join chan bool) {
    p.join = make(chan bool)

    go handleConnection(p.server, p.join)
    go handleConnection(p.client, p.join)
    
    go func() {
        logger.Println("Starting piping goroutine...")
        for {
            select {
            case serverDataIn := <- p.server.DataIn:
                time.Sleep(time.Duration(p.serverConfig.redirectDelay) * time.Millisecond)
                p.client.DataOut <- serverDataIn
            case clientDataIn := <- p.client.DataIn:
                time.Sleep(time.Duration(p.clientConfig.redirectDelay) * time.Millisecond)
                p.server.DataOut <- clientDataIn
            }
        }
    } ()

    gone := 0
    for gone < 3 {
        select {
        case <- p.join:
            gone++
        }
    }

    join <- true
}
