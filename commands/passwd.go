package commands

import (
	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

func PasswdCommand(c *cli.Context) error {
	oldPwd, newPwd, err := passwd.GetPasswordChange(c.String("password"), c.String("new-password"))
	if err != nil {
		return err
	}
	var outputUri string
	parsedData := make(map[string]interface{})
	if len(c.String("file")) > 0 {
		outputUri = c.String("file")
		err = uri.ReadMapFromURI(outputUri, uri.ReadDataFromEncryptedStream(oldPwd), parsedData)
		if err != nil {
			return err
		}
	} else {
		outputUri = c.String("output")
		sources := c.StringSlice("encrypted-input")
		if len(sources) > 0 {
			for _, src := range sources {
				err = uri.ReadMapFromURI(src, uri.ReadDataFromEncryptedStream(oldPwd), parsedData)
				if err != nil {
					return err
				}
			}
		} else {
			err = uri.ReadMapFromURI("", uri.ReadDataFromEncryptedStream(oldPwd), parsedData)
			if err != nil {
				return err
			}
		}
	}
	return uri.WriteMapToURI(outputUri, parsedData, uri.WriteDataToEncryptedStream(newPwd))
}
