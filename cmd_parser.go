package main

var supportedCmds = [...]string{"ls", "start", "stop", "plugin", "exit"}

type CmdParser struct {
	temp int
}
