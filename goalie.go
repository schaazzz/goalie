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
)

// Options is...
type Options struct {
	config Config
	shell  string
	plugin map[string]string
}

func parseArgs(options *Options) error {
	flags := flag.NewFlagSet("default", flag.ContinueOnError)

	flags.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"goalie\":\n")
		fmt.Fprintf(flag.CommandLine.Output(), " -config string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Configuration file\n")
		fmt.Fprintf(flag.CommandLine.Output(), " -shell string\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Interactive shell: none, local, remote, all (default \"none\")\n")
		fmt.Fprintf(flag.CommandLine.Output(), "        Remote shell serves by default @ 127.0.0.1:17231, use config file to customize\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), " -help")
		fmt.Fprintf(flag.CommandLine.Output(), "        Print this help menu")
		fmt.Fprintf(flag.CommandLine.Output(), "\n\n")
	}

	// Leaving out the usage string since we have a custom "Usage" function
	configFile := flags.String("config", "", "")
	shell := flags.String("shell", "none", "")
	help := flags.Bool("help", false, "")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("[Error] Argument parsing error: ", err.Error())
		return errors.New("---")
	}

	if *help {
		flags.Usage()
		return errors.New("---")
	}

	switch *shell {
	case "none", "local", "remote", "all":
		options.shell = *shell
		break
	default:
		flags.Usage()
		return errors.New("undefined shell option: \"" + *shell + "\"")
	}

	configJSON, err := ioutil.ReadFile(*configFile)
	_ = configJSON
	if err != nil {
		fmt.Printf("!!! %s\n", *configFile)
		log.Fatal("[Error] Configuration file - ", err.Error())
		return errors.New("---")
	}

	options.config = parseConfigJSON(configJSON)
	fmt.Printf("%+v\n", options.config)

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

func setupCmdParser(cmdParser *CommandParser) {
	serviceCmd := &Command{}
	serviceCmd.init(1, "service <subcommand>", "Service specific commands, a subcommand must be specified")
	serviceCmd.addSubCommand(0, "ls", "service ls", "Show list of all available services")
	serviceCmd.addSubCommand(1, "start", "service start <service>", "Start specified service")
	serviceCmd.addSubCommand(1, "stop", "service stop  <service>", "Stop specified service")
	serviceCmd.addSubCommand(2, "cmd", "service cmd   <service> <cmd>", "Forwards command to specifeid service")
	serviceCmd.addSubCommand(1, "help", "service help  <service>", "Print service specific command help")

	cmdParser.init()
	cmdParser.addCommand("service", serviceCmd)
}

func executeCmd(parsedCmd *ParsedCommand) {

}

func main() {
	join := make(chan bool)

	options := &Options{}
	if err := parseArgs(options); err != nil {
		println(err.Error())
		return
	}

	cmdParser := &CommandParser{}
	setupCmdParser(cmdParser)

	shell := &Shell{}
	shell.init(Local, cmdParser)

	parsedCmd := make(chan *ParsedCommand)
	cmdComplete := make(chan bool)

	if options.shell != "none" {
		switch options.shell {
		case "local":
			go shell.start(join, cmdComplete, parsedCmd)
		}
	}

	for {
		select {
		case parsedCmdLocal := <-parsedCmd:
			fmt.Printf("==>%+v\n", parsedCmdLocal)
			processCommand(parsedCmdLocal, &options.config)
			cmdComplete <- true
			break
		case <-join:
			return
		}
	}
}

func processCommand(parsedCmd *ParsedCommand, config *Config) {
	switch parsedCmd.commandName {
	case "service":
		processServiceCommand(parsedCmd.subCommandName, parsedCmd.args, config.Services)
	}
}

func processServiceCommand(subCmdName string, args []string, services []Service) {
	switch subCmdName {
	case "ls":
		for _, service := range services {
			fmt.Printf("%s\t", service.Name)
		}
		println()
		break
	case "start":
		// validate service name and start service
		break
	case "stop":
		// validate service name, run status and stop service
		break
	case "cmd":
		// validate service name, run status and forward command to the service
		// wait here afterwards
		break
	case "help":
		// validate service name, run status and forward command to the service
		// wait here afterwards
	default:
		log.Fatal("HTF did I end up here!?")
		break
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
