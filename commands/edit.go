package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

func EditCommand(c *cli.Context) error {
	isEdit := c.Command.FullName() == "edit"
	editor := c.String("editor")
	if editor == "" {
		editor = os.Getenv("VISUAL")
		if editor == "" {
			editor = os.Getenv("EDITOR")
			if editor == "" {
				return fmt.Errorf("editor option ('-E') not set and neither VISUAL nor EDITOR environment variables are set")
			}
		}
	}
	pwd, err := passwd.GetPassword(c.String("password"))
	if err != nil {
		return err
	}
	var outputUri string
	parsedInputData := make(map[string]interface{})
	if isEdit {
		if len(c.String("file")) > 0 {
			outputUri = c.String("file")
			err = uri.ReadMapFromURI(outputUri, uri.ReadDataFromEncryptedStream(pwd), parsedInputData)
			if err != nil {
				return err
			}
		} else {
			outputUri = c.String("output")
			sources := c.StringSlice("encrypted-input")
			if len(sources) > 0 {
				for _, src := range sources {
					err = uri.ReadMapFromURI(src, uri.ReadDataFromEncryptedStream(pwd), parsedInputData)
					if err != nil {
						return err
					}
				}
			} else {
				err = uri.ReadMapFromURI("", uri.ReadDataFromEncryptedStream(pwd), parsedInputData)
				if err != nil {
					return err
				}
			}
		}
	}
	tmpData, err := yaml.Marshal(&parsedInputData)
	if err != nil {
		return err
	}
	tempFile, err := ioutil.TempFile("", "nickelcase-")
	defer os.Remove(tempFile.Name())
	tempFile.Write(tmpData)
	tempFile.Close()

	cmd := exec.Command(editor)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	switch filepath.Base(editor) {
	case "vi", "vim", "nvi", "ex":
		cmd.Env = append(os.Environ(),
			"VIMINIT=set secure|set noswapfile|set nobackup",
			"EXINIT=set secure|set noswapfile|set nobackup")
	case "nano":
		cmd.Args = append(cmd.Args, "-w")
	}
	cmd.Args = append(cmd.Args, tempFile.Name())
	err = cmd.Run()
	if err != nil {
		return err
	}

	tempFile, err = os.OpenFile(tempFile.Name(), os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	tmpData, err = ioutil.ReadAll(tempFile)
	tempFile.Close()

	parsedOutputData := make(map[string]interface{})
	err = yaml.Unmarshal(tmpData, &parsedOutputData)
	if err != nil {
		return err
	}
	return uri.WriteMapToURI(outputUri, parsedOutputData, uri.WriteDataToEncryptedStream(pwd))
}
