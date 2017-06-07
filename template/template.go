package template

import (
	"path"
	"strings"
	"text/template"
	"time"
)

var TemplateFuncs template.FuncMap = make(template.FuncMap)

func init() {
	TemplateFuncs["base"] = path.Base
	TemplateFuncs["split"] = strings.Split
	TemplateFuncs["json"] = unmarshalJsonObject
	TemplateFuncs["jsonArray"] = unmarshalJsonArray
	TemplateFuncs["dir"] = path.Dir
	TemplateFuncs["map"] = createMap
	TemplateFuncs["getenv"] = getenv
	TemplateFuncs["join"] = strings.Join
	TemplateFuncs["datetime"] = time.Now
	TemplateFuncs["toUpper"] = strings.ToUpper
	TemplateFuncs["toLower"] = strings.ToLower
	TemplateFuncs["contains"] = strings.Contains
	TemplateFuncs["replace"] = strings.Replace
	//TemplateFuncs["lookupIP"] = LookupIP
	//TemplateFuncs["lookupSRV"] = LookupSRV
	TemplateFuncs["fileExists"] = fileExist
}
