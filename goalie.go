package main

import (
    "fmt"
    "flag"
)



func main() {
    var mode = flag.String("mode", "echo", "Mode of operation: input, proxy, pipe or echo - [default = echo]")
    var delay = flag.Int("delay", 0, "Pipe mode: Relay data after the specified delay (milliseconds) - [default = 0]")


    flag.Parse()
    fmt.Println(*mode, *delay)

    join := make(chan bool)

    pipe := &Pipe {
        ListenAddr: "127.0.0.1:17231",
        ConnectAddr: "127.0.0.1:17232",
        Delay: 2000,
    }

    pipe.Init()
    go pipe.Start(join)

    gone := 0
    for gone < 1 {
        select {
        case <- join:
            gone++
        }
    }
}
