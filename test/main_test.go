package main

import (
	"testing"

	"github.com/k1nky/cli/pkg/cli"
)

/*********************************************************************************
*     File Name           :     test/main_test.go
*     Created By          :     jonesax
*     Creation Date       :     [2017-06-26 18:34]
*     Last Modified       :     [2017-06-26 18:34]
*     Description         :
**********************************************************************************/
func TestAddCommand(t *testing.T) {

	c := cli.NewCli()

	c.AddCommand(cli.Command{})

	if len(c.Commands) != 1 {
		t.Error("Incorrect arg count")
	}
}
