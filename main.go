package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/commands"
)

var version = "master"

const (
	FLAG_DATA = iota
	FLAG_EDITOR
	FLAG_ENCRYPTED_DST
	FLAG_ENCRYPTED_SRC
	FLAG_ENV
	FLAG_ENV_PREFIX
	FLAG_FILE
	FLAG_FORMAT
	FLAG_NEW_PASSWORD_SRC
	FLAG_PASSWORD_SRC
	FLAG_PLAINTEXT_DST
	FLAG_PLAINTEXT_SRC
	FLAG_TEMPLATE_DATA
	FLAG_TEMPLATE_URI
)

const DESCRIPTION = `nickelcase allows the user to store, modify and retrieve secrets from
   encrypted files or streams with "Ansible Vault" file-format.`

func getFlagSlice(flagMap map[byte]cli.Flag, flagBytes []byte) []cli.Flag {
	flagSlice := make([]cli.Flag, len(flagBytes))
	for i := range flagBytes {
		flagSlice[i] = flagMap[flagBytes[i]]
	}
	return flagSlice
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Usage = "a command-line utility to manage secrets."
	app.Authors = []cli.Author{
		{
			Name:  "Giacomo Cariello",
			Email: "g.cariello@ieee.org",
		},
	}
	app.Description = DESCRIPTION
	app.Copyright = "Copyright 2017, Giacomo Cariello.\n   Licensed under the Apache License, Version 2.0."
	flags := map[byte]cli.Flag{
		FLAG_PASSWORD_SRC: cli.StringFlag{
			Name:  "p,password",
			Usage: "read password from `URI`",
		},
		FLAG_NEW_PASSWORD_SRC: cli.StringFlag{
			Name:  "P,new-password",
			Usage: "read new password from `URI`",
		},
		FLAG_ENCRYPTED_SRC: cli.StringSliceFlag{
			Name:  "i,encrypted-input",
			Usage: "read encrypted data from `URI`",
		},
		FLAG_PLAINTEXT_SRC: cli.StringSliceFlag{
			Name:  "I,plaintext-input",
			Usage: "read plaintext data from `URI`",
		},
		FLAG_ENCRYPTED_DST: cli.StringFlag{
			Name:  "o,output",
			Usage: "write encrypted data to `URI`",
		},
		FLAG_PLAINTEXT_DST: cli.StringFlag{
			Name:  "o,output",
			Usage: "write plaintext data to `URI`",
		},
		FLAG_DATA: cli.StringFlag{
			Name:  "d,data",
			Usage: "nickelcase archive from `DATA`",
		},
		FLAG_EDITOR: cli.StringFlag{
			Name:  "E,editor",
			Usage: "open with editor `PROG`",
		},
		FLAG_FILE: cli.StringFlag{
			Name:  "f,file",
			Usage: "modify encrypted `FILE` (same as -I FILE -O FILE)",
		},
		FLAG_ENV: cli.StringSliceFlag{
			Name:  "e,env",
			Usage: "inject `NAME` environment variable",
		},
		FLAG_ENV_PREFIX: cli.StringFlag{
			Name:  "P,env-prefix",
			Usage: "set prefix `PREFIX` for environment variable injections",
		},
		FLAG_TEMPLATE_URI: cli.StringFlag{
			Name:  "t,template",
			Usage: "read template from `URI`",
		},
		FLAG_TEMPLATE_DATA: cli.StringFlag{
			Name:  "T,template-data",
			Usage: "read template from `DATA` string",
		},
	}
	editFlags := getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_ENCRYPTED_DST, FLAG_FILE})
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "load configuration from `FILE`",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "init",
			Usage:     "Initialize an empty nickelcase",
			ArgsUsage: "[output uri]",
			Action:    commands.InitCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC}),
		},
		{
			Name:      "create",
			Usage:     "Create a nickelcase",
			ArgsUsage: "\000",
			Action:    commands.EditCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_PLAINTEXT_DST, FLAG_EDITOR}),
		},
		{
			Name:      "load",
			Usage:     "Load a nickelcase",
			ArgsUsage: "[output uri]",
			Action:    commands.CatCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_PLAINTEXT_SRC}),
		},
		{
			Name:      "dump",
			Usage:     "Dump a nickelcase",
			ArgsUsage: "[output uri]",
			Action:    commands.CatCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_PLAINTEXT_SRC}),
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List secrets contained in a nickelcase",
			ArgsUsage: "\000",
			Action:    commands.ListCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_PLAINTEXT_DST}),
		},
		{
			Name:      "get",
			Usage:     "Get secrets from a nickelcase",
			ArgsUsage: "<keys...>",
			Action:    commands.GetCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_PLAINTEXT_DST}),
		},
		{
			Name:      "set",
			Usage:     "Set a secret in a nickelcase",
			ArgsUsage: "<KEY> <VALUE>",
			Action:    commands.SetCommand,
			Flags:     editFlags,
		},
		{
			Name:      "remove",
			Aliases:   []string{"rm"},
			ArgsUsage: "<keys...>",
			Usage:     "Remove secrets from a nickelcase",
			Action:    commands.RemoveCommand,
			Flags:     editFlags,
		},
		{
			Name:      "edit",
			Usage:     "Edit a nickelcase",
			ArgsUsage: "\000",
			Action:    commands.EditCommand,
			Flags:     append(editFlags, flags[FLAG_EDITOR]),
		},
		{
			Name:   "exec",
			Usage:  "Run a command",
			Action: commands.ExecCommand,
			Flags:  getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_ENV, FLAG_ENV_PREFIX}),
		},
		{
			Name:      "template",
			Usage:     "Generate a file from a template",
			ArgsUsage: "[output uri]",
			Action:    commands.TemplateCommand,
			Flags:     getFlagSlice(flags, []byte{FLAG_PASSWORD_SRC, FLAG_ENCRYPTED_SRC, FLAG_PLAINTEXT_SRC, FLAG_TEMPLATE_URI, FLAG_TEMPLATE_DATA}),
		},
		{
			Name:      "passwd",
			Usage:     "Change nickelcase password",
			ArgsUsage: "\000",
			Action:    commands.PasswdCommand,
			Flags:     editFlags,
		},
	}
	app.Run(os.Args)
}
