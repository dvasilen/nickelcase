package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

// ExecCommand : run a command with execve(2)
func ExecCommand(c *cli.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return fmt.Errorf("Invalid number of arguments")
	}
	cmd, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	pwd, err := passwd.GetPassword(c.String("password"))
	if err != nil {
		return err
	}
	parsedData := make(map[string]interface{})
	sources := c.StringSlice("encrypted-input")
	if len(sources) > 0 {
		for _, src := range sources {
			err = uri.ReadMapFromURI(c, src, uri.ReadDataFromEncryptedStream(pwd), parsedData)
			if err != nil {
				return err
			}
		}
	} else {
		err = uri.ReadMapFromURI(c, "", uri.ReadDataFromEncryptedStream(pwd), parsedData)
		if err != nil {
			return err
		}
	}
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		envMap[pair[0]] = pair[1]
	}
	envPrefix := c.String("env-prefix")
	for _, e := range c.StringSlice("env") {
		pair := strings.SplitN(e, ":", 2)
		key := pair[0]
		var envKey string
		if len(pair) == 2 {
			envKey = pair[1]
		} else {
			envKey = pair[0]
		}
		if envPrefix != "" {
			envKey = envPrefix + envKey
		}
		var ok bool
		envMap[envKey], ok = parsedData[key].(string)
		if !ok {
			return fmt.Errorf("Archive key \"%s\" is not a string", e)
		}
	}
	var env []string
	for k, v := range envMap {
		env = append(env, strings.Join([]string{k, v}, "="))
	}
	return syscall.Exec(cmd, args[1:], env)
}
