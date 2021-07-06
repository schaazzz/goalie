package main

import "fmt"

// ...
type CommandBaseType struct {
	name        string
	description string
}

// ...
type Command struct {
	top         CommandBaseType
	subCommands []CommandBaseType
}

// ...
type CommandParser struct {
	commands []*Command
}

// ...
type ParsedCommand struct {
	command    string
	subCommand string
	args       []string
}

func (this *Command) init(name string, description string) {
	this.top.name = name
	this.top.description = description
	this.subCommands = make([]CommandBaseType, 0)
}

func (this *Command) addSubCommand(name string, description string) {
	this.subCommands = append(
		this.subCommands,
		CommandBaseType{
			name:        name,
			description: description,
		})
}

func (this *CommandParser) init() {
	this.commands = make([]*Command, 0)
}

func (this *CommandParser) addCommand(command *Command) {
	this.commands = append(this.commands, command)
}

func (this *CommandParser) printHelp() {
	println("\nAvailable commands:")
	for _, command := range this.commands {
		fmt.Printf("\n\t%s\t%s\n", command.top.name, command.top.description)

		if len(command.subCommands) > 0 {
			fmt.Printf("\tSubcommands:\n")
		}

		for _, subCommand := range command.subCommands {
			fmt.Printf("\t\t%s\t%s\n", subCommand.name, subCommand.description)
		}
	}

	println()
}

func (this *CommandParser) parseCommand() (parsedCommand *ParsedCommand) {
	return &ParsedCommand{
		command:    "",
		subCommand: "",
		args:       make([]string, 0),
	}
}
