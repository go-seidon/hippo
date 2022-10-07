package config_test

import (
	"testing"

	"github.com/go-seidon/hippo/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Package")
}

var _ = Describe("Logging Package", func() {

	Context("WithFileName function", Label("unit"), func() {
		When("parameter is specified", func() {
			It("should return result", func() {
				opt := config.WithFileName(".env")
				var res config.ConfigOption
				opt(&res)

				Expect(res.FileName).To(Equal(".env"))
			})
		})
	})

})
