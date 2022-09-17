package healthcheck_test

import (
	"fmt"

	"github.com/InVisionApp/go-health"
	"github.com/go-seidon/local/internal/healthcheck"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go Client", func() {
	Context("AddChecks function", Label("unit"), func() {
		var (
			c healthcheck.HealthClient
		)

		BeforeEach(func() {
			h := health.New()
			c, _ = healthcheck.NewGohealthClient(h)
		})

		When("configs are invalid", func() {
			It("should return error", func() {
				err := c.AddChecks([]*healthcheck.HealthConfig{})

				Expect(err).To(Equal(fmt.Errorf("configs are invalid")))
			})
		})

		When("success add config", func() {
			It("should return result", func() {
				err := c.AddChecks([]*healthcheck.HealthConfig{
					{
						Name: "check-mock",
					},
				})

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Start function", Label("unit"), func() {
		var (
			c healthcheck.HealthClient
		)

		BeforeEach(func() {
			h := health.New()
			c, _ = healthcheck.NewGohealthClient(h)
		})

		When("success start", func() {
			It("should return result", func() {
				err := c.Start()

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Stop function", Label("unit"), func() {
		var (
			c healthcheck.HealthClient
		)

		BeforeEach(func() {
			h := health.New()
			c, _ = healthcheck.NewGohealthClient(h)
		})

		When("healthcheck is not running yet", func() {
			It("should return error", func() {
				err := c.Stop()

				eErr := fmt.Errorf("Healthcheck is not running - nothing to stop")
				Expect(err).To(Equal(eErr))
			})
		})
	})

	Context("State function", Label("unit"), func() {
		var (
			c healthcheck.HealthClient
		)

		BeforeEach(func() {
			h := health.New()
			c, _ = healthcheck.NewGohealthClient(h)
		})

		When("state is invalid", func() {
			It("should return error", func() {
				states, success, err := c.State()

				Expect(states).ToNot(BeNil())
				Expect(success).To(BeFalse())
				Expect(err).To(BeNil())
			})
		})
	})
})
