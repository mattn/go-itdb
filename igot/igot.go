package main

import (
	"flag"
	"fmt"
	"github.com/gonuts/commander"
	"github.com/mattn/go-itdb"
	"log"
	"os"
	"runtime"
)

func defaultIPodPath() string {
	p := os.Getenv("IPOD_MOUNT_DIR")
	if p != "" {
		return p
	}
	if runtime.GOOS == "windows" {
		return `F:\`
	}
	return `/mnt/ipod`
}

var p = flag.String("p", defaultIPodPath(), "iPod mount directory")

func main() {
	iPod, err := itdb.New(*p)
	if err != nil {
		log.Fatal(err)
	}
	command := &commander.Command{
		UsageLine: os.Args[0],
		Short:     "iPod commander",
	}
	command.Subcommands = []*commander.Command{
		make_cmd_ls(iPod),
		make_cmd_cp(iPod),
	}

	err = command.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
