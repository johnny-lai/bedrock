package service

import (
	"github.com/gin-gonic/gin"
	"github.com/johnny-lai/bedrock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var _ = Describe("Service", func() {
	var (
		app *bedrock.Application
		svc Service
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)

		file := os.Getenv("TEST_CONFIG_YML")
		if file == "" {
			log.Fatal("Configuration file not specified. Please set TEST_CONFIG_YML variable")
		}

		app = new(bedrock.Application)
		app.Engine = gin.New()

		err := app.ReadConfigFile(file)
		if err != nil {
			log.Fatal(err)
		}

		svc = Service{}

		if err := svc.Configure(app); err != nil {
			log.Fatal(err)
		}

		if err := svc.Build(app); err != nil {
			log.Fatal(err)
		}
	})

	Describe("#Health", func() {
		It("should not raise an error", func() {
			request, _ := http.NewRequest("GET", "/health", nil)
			response := httptest.NewRecorder()
			app.Engine.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service")
}
