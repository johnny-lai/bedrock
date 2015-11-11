package bedrock

import (
	"encoding/base64"
	"io/ioutil"
	"os"
)

type TemplateContext struct {
}

// Returns an environment variable
func (c *TemplateContext) Env(key string) string {
	return os.Getenv(key)
}

// Returns the contents of the specified file. Useful for pulling secrets
// that are mounted on the filesystem into the config file.
func (c *TemplateContext) Cat(file string) string {
	bytes, _ := ioutil.ReadFile(file)
	// Ignore error
	return string(bytes)
}

// Returns a Base64 representation of the specified string.
func (c *TemplateContext) ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
