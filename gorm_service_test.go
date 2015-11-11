package bedrock

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("GormService", func() {
	var (
		app *Application
		svc *GormService
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		app = new(Application)
		svc = new(GormService)
	})

	Describe("#Configure", func() {
		It("should read the airbrake dictionary", func() {
			app.ConfigBytes = []byte("db:\n  user: abc\n  password: whatever\n  host: localhost\n  database: a_db")

			// Verify initial state
			Expect(svc.Config.User).To(Equal(""))
			Expect(svc.Config.Password).To(Equal(""))
			Expect(svc.Config.Host).To(Equal(""))
			Expect(svc.Config.Database).To(Equal(""))

			// Load config
			Expect(svc.Configure(app)).To(BeNil())

			// Verify config was loaded properly
			Expect(svc.Config.User).To(Equal("abc"))
			Expect(svc.Config.Password).To(Equal("whatever"))
			Expect(svc.Config.Host).To(Equal("localhost"))
			Expect(svc.Config.Database).To(Equal("a_db"))
		})
	})

	Describe("#Build", func() {
		It("should not panic", func() {
			// Build
			Expect(svc.Build(app)).To(BeNil())
		})
	})
})

func TestGormService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GormService")
}
