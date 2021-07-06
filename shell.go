package main

import (
	"bufio"
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
	shell ShellType
	count int
	wg    sync.WaitGroup
}

func (this *Shell) init(shell ShellType) {
	this.shell = shell

	if this.shell == Both {
		this.count = 2
	}
}

func (this *Shell) start(join chan bool, parsedCommand chan *ParsedCommand, commandComplete chan bool) {
	usrInput := make(chan string)

	if this.shell == Local || this.shell == Both {
		this.wg.Add(1)
		go this.startLocalShell(usrInput, commandComplete)
	}

	if this.shell == Remote || this.shell == Both {
		this.wg.Add(1)
		go this.startRemoteShell(usrInput, commandComplete)
	}

	exit := make(chan bool)
	go this.processCmdStr(usrInput, parsedCommand, exit)

	this.wg.Wait()
	exit <- true
	join <- true
}

func (this *Shell) processCmdStr(usrInput chan string, parsedCommand chan *ParsedCommand, exit chan bool) {
loop:
	for {
		select {
		case cmdStr := <-usrInput:
			println("Received:", cmdStr)
		case <-exit:
			break loop
		}
	}
}

func (this *Shell) startLocalShell(usrInput chan string, commandComplete chan bool) {
	defer this.wg.Done()

	writer := bufio.NewWriter(os.Stdout)
	reader := bufio.NewReader(os.Stdin)

	for {
		writer.WriteString("#> ")
		writer.Flush()

		cmdStr, _ := reader.ReadString('\n')
		usrInput <- cmdStr
		<-commandComplete
	}
}

func (this *Shell) startRemoteShell(usrInput chan string, commandComplete chan bool) {
	defer this.wg.Done()
}
