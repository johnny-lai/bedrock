package bedrock

import (
	"io/ioutil"
	"os"
)

type TemplateContext struct {
}

func (c *TemplateContext) Env(key string) string {
	return os.Getenv(key)
}

func (c *TemplateContext) Cat(file string) string {
	bytes, _ := ioutil.ReadFile(file)
	// Ignore error
	return string(bytes)
}

func (c *TemplateContext) Base64(s string) string {
	return s
}
