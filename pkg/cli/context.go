package cli

type Context struct {
	Name     string
	keys     map[string]interface{}
	Commands []string
}

type Args struct {
	Values []string
}

func (ctx *Context) ResetCommands() {
	ctx.Commands = ctx.Commands[:0]
}

func (ctx *Context) PushCommand(name string) {
	ctx.Commands = append(ctx.Commands, name)
}
