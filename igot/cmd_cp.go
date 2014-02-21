package main

import (
	"fmt"
	"github.com/gonuts/commander"
	"github.com/mattn/go-itdb"
	"log"
)

func make_cmd_cp(iPod *itdb.IPod) *commander.Command {
	cmd_cp := func(cmd *commander.Command, args []string) error {
		// TODO
		return nil
	}

	return &commander.Command{
		Run:       cmd_cp,
		UsageLine: "cp [file]",
		Short:     "copy track into iPod",
	}
}
