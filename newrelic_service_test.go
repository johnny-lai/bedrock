package bedrock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gin-gonic/gin"
	"testing"
)

var _ = Describe("NewRelicService", func() {
	var (
		app  *Application
		svc  *NewRelicService
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		app = new(Application)
		svc = new(NewRelicService)
	})

	Describe("#Configure", func() {
		It("should read the newrelic dictionary", func() {
			app.ConfigBytes = []byte("newrelic:\n  licensekey: abc\n  appname: whatever\n  verbose: true")

			// Verify initial state
			Expect(svc.Config.LicenseKey).To(Equal(""))
			Expect(svc.Config.AppName).To(Equal(""))
			Expect(svc.Config.Verbose).To(Equal(false))

			// Load config
			Expect(svc.Configure(app)).To(BeNil())

			// Verify config was loaded properly
			Expect(svc.Config.LicenseKey).To(Equal("abc"))
			Expect(svc.Config.AppName).To(Equal("whatever"))
			Expect(svc.Config.Verbose).To(Equal(true))
		})
	})

	Describe("#Build", func() {
		It("should not panic", func() {
			app.Engine = gin.New()

			// Build
			Expect(svc.Build(app)).To(BeNil())
		})
	})
})

func TestNewRelicService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewRelicService")
}