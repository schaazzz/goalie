package main

import (
    "os"
    _ "time"
    "log"
    tcp "github.com/schaazzz/golibs/network/tcp"
)

type echoConfig struct {
    address         string
}

type Echo struct {
    serverConfig    * echoConfig
    clientConfig    * echoConfig
    server          * tcp.Connection
    client          * tcp.Connection
    join            chan bool
}

var echoLogger * log.Logger

func (e * Echo) Init(echoJSON []EchoJSON) {
    for _, element := range echoJSON {
        config := &echoConfig {element.Address}
        if element.Role == "server" {
            e.serverConfig = config
        } else {
            e.clientConfig = config
        }
    }

    echoLogger = log.New(os.Stdout, "[ECHO MAIN] ", log.Lmicroseconds)
    
    if e.serverConfig != nil {
        e.server = &tcp.Connection {
                    Channels: tcp.Channels {
                        Ctrl        : make(chan string),
                        DataIn      : make(chan * tcp.DataChunk),
                        DataOut     : make(chan * tcp.DataChunk),
                        Done        : make(chan bool),
                        Connected   : make(chan bool),
                        Panic       : make(chan bool),
                    },
                    Server  : true,
                    Address : e.serverConfig.address,
                    Name    : "ECHO SERVER",
                }
    }

    if e.clientConfig != nil {
        e.client = &tcp.Connection {
                    Channels: tcp.Channels {
                        Ctrl        : make(chan string),
                        DataIn      : make(chan * tcp.DataChunk),
                        DataOut     : make(chan * tcp.DataChunk),
                        Done        : make(chan bool),
                        Connected   : make(chan bool),
                        Panic       : make(chan bool),
                    },
                    Server  : false,
                    Address : e.clientConfig.address,
                    Name    : "ECHO CLIENT",
                }
    }
}

func (e * Echo) Start(join chan bool) {
    e.join = make(chan bool)

    if e.serverConfig != nil {
        go handleConnection(e.server, echoLogger, e.join)
    }

    if e.clientConfig != nil {
        go handleConnection(e.client, echoLogger, e.join)
    }
    
    go func() {
        echoLogger.Println("Starting piping goroutine...")
        for {
            select {
            case serverDataIn, ok := <- e.server.DataIn:
                if ok {
                    e.server.DataOut <- serverDataIn
                } else {
                    echoLogger.Println("=====================")
                    e.server.DataIn = nil
                    e.server.DataOut = nil
                }
            case clientDataIn, ok := <- e.client.DataIn:
                if ok {
                    e.client.DataOut <- clientDataIn
                } else {
                    echoLogger.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$")
                    e.client.DataIn = nil
                    e.client.DataOut = nil
                }
            }
        }
    } ()

    gone := 0
    for gone < 3 {
        select {
        case <- e.join:
            gone++
        }
    }

    join <- true
}