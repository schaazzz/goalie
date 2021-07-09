package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type ShellType int

// ...
const (
	Local ShellType = iota
	Remote
	Both
)

// ...
type Shell struct {
	shell     ShellType
	cmdParser *CommandParser
	wg        sync.WaitGroup
}

func (this *Shell) init(shell ShellType, cmdParser *CommandParser) {
	this.shell = shell
	this.cmdParser = cmdParser
}

func (this *Shell) start(join chan bool, cmdProcessed chan bool, parsedCmd chan *ParsedCommand) {
	if this.shell == Local || this.shell == Both {
		this.wg.Add(1)
		go this.startLocalShell(cmdProcessed)
	}

	if this.shell == Remote || this.shell == Both {
		this.wg.Add(1)
		go this.startRemoteShell(cmdProcessed)
	}

	this.wg.Wait()
	join <- true
}

func (this *Shell) startLocalShell(cmdProcessed chan bool) {
	defer this.wg.Done()

	writer := bufio.NewWriter(os.Stdout)
	reader := bufio.NewReader(os.Stdin)

	for {
		writer.WriteString("#> ")
		writer.Flush()

		usrInput, _ := reader.ReadString('\n')
		parsedCmd, err := this.cmdParser.parseCommand(usrInput)

		if err == nil {
			fmt.Printf("parsedCmd %+v\n", parsedCmd)
			if parsedCmd.commandName == "exit" {
				continue
			} else if parsedCmd.commandName == "help" {
				this.cmdParser.printHelp(nil)
			} else {
				//<-cmdProcessed
				continue
			}
		} else {
			this.cmdParser.printHelp(&err)
		}

	}
}

func (this *Shell) startRemoteShell(cmdProcessed chan bool) {
	defer this.wg.Done()
}
