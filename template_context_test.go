package bedrock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

var _ = Describe("TemplateContext", func() {
	var (
		tc TemplateContext
	)

	Describe("#Env", func() {
		It("should read the environment variables", func() {
			os.Setenv("TEMPLATE_CONTEXT_TEST", "something")

			ans := tc.Env("TEMPLATE_CONTEXT_TEST")
			Expect(ans).To(Equal("something"))
		})
	})

	Describe("#Cat", func() {
		It("should read the contents of a file", func() {
			f, err := ioutil.TempFile("", "template_context_test")
			Expect(err).To(BeNil())

			file := f.Name()
			defer os.Remove(file)

			f.WriteString("something nothing")
			f.Close()

			ans := tc.Cat(file)
			Expect(ans).To(Equal("something nothing"))
		})
	})

	Describe("#ToBase64", func() {
		It("should convert to base64", func() {
			ans := tc.ToBase64("something")
			Expect(ans).To(Equal("c29tZXRoaW5n"))
		})
	})
})

func TestTemplateContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TemplateContext")
}
