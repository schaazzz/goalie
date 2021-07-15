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
		go this.startLocalShell(cmdProcessed, parsedCmd)
	}

	if this.shell == Remote || this.shell == Both {
		this.wg.Add(1)
		go this.startRemoteShell(cmdProcessed, parsedCmd)
	}

	this.wg.Wait()
	join <- true
}

func (this *Shell) startLocalShell(cmdProcessed chan bool, parsedCmd chan *ParsedCommand) {
	defer this.wg.Done()

	writer := bufio.NewWriter(os.Stdout)
	reader := bufio.NewReader(os.Stdin)

	for {
		writer.WriteString("#> ")
		writer.Flush()

		usrInput, _ := reader.ReadString('\n')
		parsedCmdLocal, err := this.cmdParser.parseCommand(usrInput)

		if err == nil {
			fmt.Printf("parsedCmd: %+v\n", parsedCmdLocal)
			if parsedCmdLocal.commandName == "exit" {
				return
			} else if parsedCmdLocal.commandName == "help" {
				this.cmdParser.printHelp(nil)
			} else {
				parsedCmd <- parsedCmdLocal
				<-cmdProcessed
				continue
			}
		} else {
			this.cmdParser.printHelp(&err)
		}

	}
}

func (this *Shell) startRemoteShell(cmdProcessed chan bool, parsedCmd chan *ParsedCommand) {
	defer this.wg.Done()
}
