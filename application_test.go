package bedrock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("App", func() {
	var (
		app *Application
	)

	BeforeEach(func() {
		app = new(Application)
	})
})

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App")
}
