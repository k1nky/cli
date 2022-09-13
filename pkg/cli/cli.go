//Package cli is a simple package to help implement interactive command line interfaces in golang.
//One of the main reasons behind generate it is that there is a lack of subcommand support in other packages.
package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/k1nky/cli/pkg/parser"
)

//Cli structure contains configuration and commands
type Cli struct {
	Commands       []Command
	Current        *Context
	Contexts       map[string]*Context
	ReadlineConfig *readline.Config
	Scanner        *readline.Instance
	Parser         parser.Parser
	OnExit         func()
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

var completer = readline.NewPrefixCompleter()

//NewCli creates a new instance of Cli
//It returns a pointer to the Cli object
func NewCli() *Cli {
	c := &Cli{
		Parser: &parser.QuoteParser{},
		OnExit: func() {
			fmt.Println("bye")
		},
		Contexts: make(map[string]*Context),
	}
	c.SetContext("")
	l, err := readline.NewEx(&readline.Config{
		HistoryFile:     "/tmp/readline.tmp",
		Prompt:          "> ",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		//TODO some weird version error broke this
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	c.Scanner = l

	return c
}

func (cli *Cli) SetContext(name string) *Context {
	ctx, exists := cli.Contexts[name]
	if !exists {
		ctx = &Context{
			Name: name,
		}
	}
	cli.Current = ctx
	return ctx
}

// Close runs exit command
func (cli *Cli) Close() {
	cli.OnExit()
	os.Exit(0)
}

//AddCommand is a method on Cli takes Command as input
//This appends to the current command list to search through for input
func (cli *Cli) AddCommand(c Command) {
	cli.Commands = append(cli.Commands, c)

	//recusively add command names to completer
	pc := readline.PcItem(c.Name)
	cli.recurseCompletion(c.SubCommands, pc, 0)
	completer.Children = append(completer.Children, pc)
}

func (cli *Cli) peakChildren(c []Command, name string) *Command {
	for _, cmd := range c {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func (cli *Cli) recurseCompletion(c []Command, pc *readline.PrefixCompleter, i int) error {
	for _, cmd := range c {
		p := readline.PcItem(cmd.Name)
		pc.Children = append(pc.Children, p)

		if len(cmd.SubCommands) > 0 {
			cli.recurseCompletion(cmd.SubCommands, p, i+1)
		}
	}
	return nil
}

func (cli *Cli) recurseHelp(c []Command, rootCommands []string, offset int) {

	for _, cmd := range c {
		for i := 0; i < offset; i++ {
			fmt.Printf("\t")
		}
		for _, n := range rootCommands {
			if strings.Compare(n, cmd.Name) == 0 {
				offset = 0
			}
		}
		fmt.Printf("[%s]: %s\n", cmd.Name, cmd.Help)
		if len(cmd.SubCommands) > 0 {
			cli.recurseHelp(cmd.SubCommands, rootCommands, offset+1)
		}
	}
}

func (cli *Cli) parseSystemCommands(input []string) error {
	if input[0] == "exit" {
		cli.Close()
	}
	if input[0] == "clear" {
		print("\033[H\033[2J")
	}
	if input[0] == "help" {

		var rootCommands []string
		for _, r := range cli.Commands {
			rootCommands = append(rootCommands, r.Name)
		}
		cli.recurseHelp(cli.Commands, rootCommands, 0)
	}

	return nil
}

func (cli *Cli) recurse(ctx *Context, c []Command, args []string, i int) error {
	for _, cmd := range c {
		if i > len(args) {
			return nil
		}
		if cmd.Name == args[i] {
			ctx.PushCommand(cmd.Name)
			if len(args) > i+1 {
				if child := cli.peakChildren(cmd.SubCommands, args[i+1]); child != nil {
					cli.recurse(ctx, cmd.SubCommands, args, i+1)
				} else {
					cli.runCommand(ctx, cmd, args, i+1)
				}
			} else {
				cli.runCommand(ctx, cmd, args, i+1)
			}

		}
	}
	return nil
}

func (cli *Cli) runCommand(ctx *Context, cmd Command, args []string, i int) error {
	cmd.Func(ctx, Args{
		Values: args[i:],
	})
	fmt.Println()
	return nil
}

func (cli *Cli) findCommand(ctx *Context, input string) error {
	parsed := cli.Parser.Parse(input)
	ctx.ResetCommands()
	if len(parsed) == 0 {
		cli.Warning("No input detected")
		return nil
	}
	if systemCmd := cli.parseSystemCommands(parsed); systemCmd != nil {
		return nil
	}
	error := cli.recurse(ctx, cli.Commands, parsed, 0)
	if error != nil {
		return error
	}
	return nil
}

func (cli *Cli) readline() string {

	text, _ := cli.Scanner.Readline()
	cli.Scanner.SaveHistory(text)
	return text
}

//Run is the primary entrypoint to start blocking and reading user input
func (cli *Cli) Run() {

	if len(os.Args) > 1 && os.Args[1] == "unattended" {
		err := cli.findCommand(cli.Current, strings.Join(os.Args[2:], " "))
		if err != nil {
			color.Red(err.Error())
		}
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			cli.Close()
		}
	}()

	for {
		//Get user input
		fmt.Print(cli.Scanner.Config.Prompt)

		text := cli.readline()

		err := cli.findCommand(cli.Current, text)
		if err != nil {
			cli.Error(err.Error())
		}
	}
}

func (cli *Cli) Error(msg string) {
	fmt.Printf("%s: %s", color.RedString("%s", "error"), msg)
}

func (cli *Cli) Warning(msg string) {
	fmt.Printf("%s: %s", color.YellowString("%s", "warning"), msg)
}
