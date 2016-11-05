package bedrock

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"text/template"
)

type ConfigBytes []byte

// Reads the specified config file. Note that bedrock.Application will process
// the config file, using text/template, with the following extra functions:
//
//     {{.Env "ENVIRONMENT_VARIABLE"}}
//     {{.Cat "File name"}}
//     {{.Base64 "a string"}}
func ReadConfigFile(file string) (ConfigBytes, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, errors.New("config path not valid")
	}

	tmpl, err := template.New(path.Base(file)).ParseFiles(file)
	if err != nil {
		return nil, err
	}

	var configBytes bytes.Buffer
	tc := TemplateContext{}
	err = tmpl.Execute(&configBytes, &tc)
	if err != nil {
		return nil, err
	}

	return ConfigBytes(configBytes.Bytes()), nil
}

func (c ConfigBytes) Unmarshal(dst interface{}) error {
	return yaml.Unmarshal(c, dst)
}
