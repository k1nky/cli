package cli

// Command structure is composed of the function to route too plus information
// Subcommands can also be nested
type Command struct {
	Name        string
	Help        string
	Func        func(ctx *Context, args Args)
	SubCommands []Command
}

//NewCommand is the struct initialization for a command
//It should return *Command
func NewCommand() *Command {
	c := new(Command)
	c.SubCommands = []Command{}
	return c
}