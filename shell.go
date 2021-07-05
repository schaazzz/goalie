package main

import (
	"bufio"
	"os"
	"sync"
	"time"
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

func (this *Shell) start(join chan bool, cmd chan string, args chan []string) {
	usrInput := make(chan string)

	if this.shell == Local || this.shell == Both {
		this.wg.Add(1)
		go this.startLocalShell(usrInput)
	}

	if this.shell == Remote || this.shell == Both {
		this.wg.Add(1)
		go this.startRemoteShell(usrInput)
	}

	exit := make(chan bool)
	go this.processCmdStr(usrInput, cmd, args, exit)

	this.wg.Wait()
	exit <- true
	join <- true
}

func (this *Shell) processCmdStr(usrInput chan string, cmd chan string, args chan []string, exit chan bool) {
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

func (this *Shell) startLocalShell(usrInput chan string) {
	defer this.wg.Done()

	writer := bufio.NewWriter(os.Stdout)
	reader := bufio.NewReader(os.Stdin)

	for {
		writer.WriteString("#> ")
		writer.Flush()

		cmdStr, _ := reader.ReadString('\n')
		usrInput <- cmdStr
		time.Sleep(250 * time.Millisecond)
	}
}

func (this *Shell) startRemoteShell(usrInput chan string) {
	defer this.wg.Done()
}
