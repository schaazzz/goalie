package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	tcp "github.com/schaazzz/golibs/network/tcp"
	//"rsc.io/getopt"
)

// Config is...
type Config struct {
	shell  string
	plugin map[string]string
}

func parseArgs(config *Config) error {
	flags := flag.NewFlagSet("default", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"goalie\":\n")
		fmt.Fprintf(flag.CommandLine.Output(), " -config string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Config file for the selected mode\n")
		fmt.Fprintf(flag.CommandLine.Output(), " -shell string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Interactive shell: none, local, remote, all (default \"none\")\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Remote shell serves by default @ 127.0.0.1:17231, use config file to customize\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), " -help")
		fmt.Fprintf(flag.CommandLine.Output(), "        Print this help menu")
		fmt.Fprintf(flag.CommandLine.Output(), "\n\n")
	}

	configFile := flags.String("config", "", "Config file for the selected mode")
	shell := flags.String("shell", "none", "Interactive shell: none, local, remote, all\n"+
		"Remote shell serves by default @ 127.0.0.1:17231, use config file to customize")

	help := flags.Bool("help", false, "Print this help menu")

	_, _ /*configJSON, _ :*/ = ioutil.ReadFile(*configFile)

	err := flags.Parse(os.Args[1:])

	if err != nil {
		return err
	}
	if *help {
		flags.Usage()
		return errors.New("---")
	}

	switch *shell {
	case "none", "local", "remote", "all":
		fmt.Printf("Shell was valid: %s\n", *shell)
		break
	default:
		flags.Usage()
		return errors.New("undefined shell option: \"" + *shell + "\"")
	}

	// if *mode == "pipe" {
	// 	pipe := &Pipe{}
	// 	pipe.Init(parsePipeConfigJSON(configJSON)[0:])
	// 	go pipe.Start(join)
	// } else if *mode == "echo" {
	// 	echo := &Echo{}
	// 	echo.Init(parseEchoConfigJSON(configJSON)[0:])
	// 	go echo.Start(join)
	// }

	return nil
}

func main() {
	join := make(chan bool)

	config := &Config{}
	if err := parseArgs(config); err != nil {
		println(err.Error())
		return
	}

	gone := 0
	for gone < 1 {
		select {
		case <-join:
			gone++
		}
	}
}

func handleConnection(c *tcp.Connection, logger *log.Logger, join chan bool) {
reset:
	go c.Start()

forever:
	for {
		select {
		case <-c.Panic:
			logger.Println(c.Name, "panicked, resetting in 3 seconds!")
			time.Sleep(3 * time.Second)
			goto reset
		case serverConnectionState := <-c.Connected:
			if serverConnectionState {
				c.Ctrl <- "start"
			} else {
				break forever
			}
		case <-c.Done:
			break forever
		}
	}

	join <- true
}
