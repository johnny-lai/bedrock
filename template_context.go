package bedrock

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
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

// Returns a lower-case representation of the specified string.
func (c *TemplateContext) ToLower(s string) string {
	return strings.ToLower(s)
}

// Returns a upper-case representation of the specified string.
func (c *TemplateContext) ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Returns a upper-case representation of the specified string.
func (c *TemplateContext) ToUnderscore(s string) string {
	return strings.Replace(s, "-", "_", -1)
}
