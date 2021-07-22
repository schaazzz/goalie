package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ...
type CommandBaseType struct {
	numArgs     int
	usage       string
	description string
}

// ...
type Command struct {
	base        CommandBaseType
	subCommands map[string]CommandBaseType
}

// ...
type CommandParser struct {
	commands map[string]*Command
}

// ...
type ParsedCommand struct {
	CommandName    string
	SubCommandName string
	Args           []string
}

func (this *Command) Init(numArgs int, usage string, description string) {
	this.base.numArgs = numArgs
	this.base.usage = usage
	this.base.description = description
	this.subCommands = make(map[string]CommandBaseType)
}

func (this *Command) AddSubCommand(numArgs int, name string, usage string, description string) {
	this.subCommands[name] = CommandBaseType{
		numArgs:     numArgs,
		usage:       usage,
		description: description,
	}
}

func (this *CommandParser) Init() {
	this.commands = make(map[string]*Command)

	helpCmd := &Command{}
	helpCmd.Init(0, "help", "Print this help menu")
	this.AddCommand("help", helpCmd)

	exitCmd := &Command{}
	exitCmd.Init(0, "exit", "Exit this shell")
	this.AddCommand("exit", exitCmd)
}

func (this *CommandParser) AddCommand(name string, cmd *Command) {
	this.commands[name] = cmd
}

func (this *CommandParser) PrintHelp(err *error) {
	if err != nil {
		fmt.Printf("Error: %v\n", *err)
	}

	println("\nAvailable commands:")
	for _, command := range this.commands {
		fmt.Printf("\n\t%-30s     %s\n", command.base.usage, command.base.description)

		if len(command.subCommands) > 0 {
			fmt.Printf("\tSubcommands:\n")
		}

		for _, subCommand := range command.subCommands {
			fmt.Printf("\t\t%-30s     %s\n", subCommand.usage, subCommand.description)
		}
	}

	println()
}

func (this *Command) extractValidateSubCommand(tokens []string) (string, []string, error) {
	var err error = nil
	var cmdName string = ""

	fmt.Printf("==> tokens: %+v\n", tokens)

	if subCommand, ok := this.subCommands[tokens[0]]; ok {
		if len(tokens[1:]) >= subCommand.numArgs {
			cmdName = tokens[0]
		} else {
			err = errors.New("incorrect number of arguments")
		}
	} else {
		err = errors.New("unknown subcommand: " + tokens[0])
	}

	fmt.Printf("<== tokens: %+v\n", tokens[1:])
	return cmdName, tokens[1:], err
}

func (this *CommandParser) extractValidateCommand(tokens []string) (string, []string, error) {
	var err error = nil
	var cmdName string = ""

	if command, ok := this.commands[tokens[0]]; ok {
		if len(tokens[1:]) >= command.base.numArgs {
			cmdName = tokens[0]
		} else {
			err = errors.New("incorrect number of arguments")
		}
	} else {
		err = errors.New("unknown command: " + tokens[0])
	}

	return cmdName, tokens[1:], err
}

func (this *CommandParser) ParseCommand(cmdStr string) (*ParsedCommand, error) {
	re := regexp.MustCompile(`\s+`)
	cmdStr = re.ReplaceAllString(cmdStr, " ")
	tokens := strings.Split(cmdStr, " ")
	var err error = errors.New("unknown command")

	cmdName, tokens, err := this.extractValidateCommand(tokens)

	var subCmdName = ""
	var args []string

	if err == nil {
		cmd := this.commands[cmdName]

		if len(cmd.subCommands) > 0 {
			subCmdName, args, err = cmd.extractValidateSubCommand(tokens)
		}
	}

	if err == nil {
		return &ParsedCommand{
			CommandName:    cmdName,
			SubCommandName: subCmdName,
			Args:           args,
		}, nil
	} else {
		return nil, err
	}
}
