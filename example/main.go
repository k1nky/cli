package main

import (
	"fmt"

	"github.com/k1nky/cli/pkg/cli"
)

func AddCommands(c *cli.Cli) {
	c.AddCommand(cli.Command{
		Name: "github",
		Help: "github primary command interface",
		Func: func(ctx *cli.Context, args cli.Args) {
			c.Warning("I do nothing...")
		},
		SubCommands: []cli.Command{
			{
				Name: "login",
				Help: "access token to github",
				Func: func(ctx *cli.Context, args cli.Args) {
					if len(args.Values) == 0 {
						fmt.Println("Failed login")
						return
					}
					fmt.Printf("Logged in %s", args.Values[0])
				},
			},
			{
				Name: "logout",
				Help: "allows you to logout from github",
				Func: func(ctx *cli.Context, args cli.Args) {
					if len(args.Values) == 0 {
						fmt.Println("Failed logout")
						return
					}
					fmt.Printf("Logged out with username %s\n", args.Values[0])
				},
			},
		},
	})
	c.AddCommand(cli.Command{
		Name: "sql",
		Help: "sql primary command interface",
		Func: func(ctx *cli.Context, args cli.Args) {
			fmt.Println("I do nothing...")
		},
		SubCommands: []cli.Command{
			{
				Name: "login",
				Help: "access token to github",
				Func: func(ctx *cli.Context, args cli.Args) {
					if len(args.Values) == 0 {
						fmt.Println("Failed login")
						return
					}
					fmt.Println(ctx.Commands)
					fmt.Printf("Logged in %s", args.Values[0])
				},
			},
			{
				Name: "logout",
				Help: "allows you to logout from github",
				Func: func(ctx *cli.Context, args cli.Args) {
					if len(args.Values) == 0 {
						fmt.Println("Failed logout")
						return
					}
					fmt.Printf("Logged out with username %s\n", args.Values[0])
				},
				SubCommands: []cli.Command{
					{
						Name: "defer",
						Help: "Defer a logout",
						Func: nil,
					},
				},
			},
		}})
}
func main() {

	c := cli.NewCli()
	AddCommands(c)
	c.Run()
}
