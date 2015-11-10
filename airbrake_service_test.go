package bedrock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gin-gonic/gin"
	"testing"
)

var _ = Describe("AirbrakeService", func() {
	var (
		app  *Application
		svc  *AirbrakeService
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		app = new(Application)
		svc = new(AirbrakeService)
	})

	Describe("#Configure", func() {
		It("should read the airbrake dictionary", func() {
			app.ConfigBytes = []byte("airbrake:\n  projectid: 222\n  projectkey: something")

			// Verify initial state
			Expect(svc.Config.ProjectID).To(Equal(int64(0)))
			Expect(svc.Config.ProjectKey).To(Equal(""))
			Expect(svc.Notifier).To(BeNil())

			// Load config
			Expect(svc.Configure(app)).To(BeNil())

			// Verify config was loaded properly
			Expect(svc.Config.ProjectID).To(Equal(int64(222)))
			Expect(svc.Config.ProjectKey).To(Equal("something"))

			// Should not be initialized in configure
			Expect(svc.Notifier).To(BeNil())
		})
	})

	Describe("#Build", func() {
		It("should set OnException", func() {
			app.Engine = gin.Default()

			// Verify initial state
			Expect(svc.Notifier).To(BeNil())
			Expect(app.OnException).To(BeNil())

			// Build
			Expect(svc.Build(app)).To(BeNil())

			// Verify state change
			Expect(svc.Notifier).NotTo(BeNil())
			Expect(app.OnException).NotTo(BeNil())
		})
	})
})

func TestAirbrakeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AirbrakeService")
}