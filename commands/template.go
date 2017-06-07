package commands

import (
	"text/template"

	"github.com/urfave/cli"

	"github.com/giacomocariello/nickelcase/passwd"
	funcs "github.com/giacomocariello/nickelcase/template"
	"github.com/giacomocariello/nickelcase/uri"
)


// TemplateCommand : compute files from templates using nickelcase maps as context
func TemplateCommand(c *cli.Context) error {
	pwd, err := passwd.GetPassword(c.String("password"))
	if err != nil {
		return err
	}
	parsedData := make(map[string]interface{})
	encryptedSources := c.StringSlice("encrypted-input")
	plaintextSources := c.StringSlice("plaintext-input")
	for _, src := range encryptedSources {
		err = uri.ReadMapFromURI(c, src, uri.ReadDataFromEncryptedStream(pwd), parsedData)
		if err != nil {
			return err
		}
	}
	for _, src := range plaintextSources {
		err = uri.ReadMapFromURI(c, src, uri.ReadDataFromPlaintextStream(), parsedData)
		if err != nil {
			return err
		}
	}
	tmplString := c.String("template-data")
	if tmplString == "" {
		tmplData, err := uri.ReadDataFromURI(c, c.String("template"), uri.ReadDataFromPlaintextStream())
		if err != nil {
			return err
		}
		tmplString = string(tmplData)
	}
	tmpl, err := template.New("").Parse(tmplString)
	if err != nil {
		return err
	}
	tmpl.Funcs(funcs.TemplateFuncs)
	stream, err := uri.GetOutputStreamFromURI(c, c.Args().Get(0))
	if err != nil {
		return err
	}
	if err = tmpl.Execute(stream, parsedData); err != nil {
		return err
	}
	return nil
}
