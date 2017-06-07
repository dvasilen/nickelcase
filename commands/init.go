package commands

import (
	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	"github.com/giacomocariello/nickelcase/uri"
)

// InitCommand : initialize an empty nickelcase
func InitCommand(c *cli.Context) error {
	pwd, err := passwd.GetNewPassword(c, c.String("password"))
	if err != nil {
		return err
	}
	parsedData := make(map[string]interface{})
	return uri.WriteMapToURI(c, c.Args().Get(0), parsedData, uri.WriteDataToEncryptedStream(pwd))
}
