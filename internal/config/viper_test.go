package config_test

import (
	"os"
	"time"

	"github.com/go-seidon/hippo/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Viper Package", func() {

	Context("NewViperConfig function", Label("unit"), func() {
		When("parameter is not specified", func() {
			It("should return result", func() {
				res, err := config.NewViperConfig()

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("file name is specified", func() {
			It("should return result", func() {
				res, err := config.NewViperConfig(config.WithFileName(".env"))

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Get function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				os.Setenv("TEST_KEY", "TEST_VALUE")
				res, err := cfg.Get("TEST_KEY")

				Expect(res).To(Equal("TEST_VALUE"))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.Get("TEST_UNAVAILABLE_KEY")

				Expect(res).To(BeNil())
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("GetBool function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				os.Setenv("TEST_KEY", "true")
				res, err := cfg.GetBool("TEST_KEY")

				Expect(res).To(Equal(true))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.GetBool("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(false))
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("GetFloat64 function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				os.Setenv("TEST_KEY", "1.25")
				res, err := cfg.GetFloat64("TEST_KEY")

				Expect(res).To(Equal(1.25))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.GetFloat64("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(float64(0)))
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("GetInt function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				os.Setenv("TEST_KEY", "1")
				res, err := cfg.GetInt("TEST_KEY")

				Expect(res).To(Equal(1))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.GetInt("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(0))
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("GetString function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				os.Setenv("TEST_KEY", "TEST_VALUE")
				res, err := cfg.GetString("TEST_KEY")

				Expect(res).To(Equal("TEST_VALUE"))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.GetString("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(""))
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("GetTime function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				time, _ := time.Parse("2006-01-02", "2022-07-01")
				os.Setenv("TEST_KEY", "2022-07-01")
				res, err := cfg.GetTime("TEST_KEY")

				Expect(res).To(Equal(time))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.GetTime("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(time.Time{}))
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("GetDuration function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return result", func() {
				duration, _ := time.ParseDuration("60s")
				os.Setenv("TEST_KEY", "60s")
				res, err := cfg.GetDuration("TEST_KEY")

				Expect(res).To(Equal(duration))
				Expect(err).To(BeNil())
			})
		})

		When("config is not available", func() {
			It("should return result", func() {
				res, err := cfg.GetDuration("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(time.Duration(0)))
				Expect(err).To(Equal(config.ErrorConfigNotFound))
			})
		})
	})

	Context("Set function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("success set config", func() {
			It("should return result", func() {
				err := cfg.Set("TEST_KEY", "TEST_VALUE")
				val, _ := cfg.GetString("TEST_KEY")

				Expect(err).To(BeNil())
				Expect(val).To(Equal("TEST_VALUE"))
			})
		})
	})

	Context("SetDefault function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_DEFAULT_KEY", "")
		})

		When("success set config", func() {
			It("should return result", func() {
				err := cfg.SetDefault("TEST_DEFAULT_KEY", "TEST_DEFAULT_VALUE")
				val, _ := cfg.GetString("TEST_DEFAULT_KEY")

				Expect(err).To(BeNil())
				Expect(val).To(Equal("TEST_DEFAULT_VALUE"))
			})
		})
	})

	Context("IsSet function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		AfterEach(func() {
			os.Setenv("TEST_KEY", "")
		})

		When("config is available", func() {
			It("should return true", func() {
				cfg.Set("TEST_KEY", "TEST_VALUE")

				res, err := cfg.IsSet("TEST_KEY")

				Expect(res).To(Equal(true))
				Expect(err).To(BeNil())
			})
		})

		When("config is unavailable", func() {
			It("should return false", func() {
				res, err := cfg.IsSet("TEST_UNAVAILABLE_KEY")

				Expect(res).To(Equal(false))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("LoadConfig function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		When("failed load config", func() {
			It("should return error", func() {
				err := cfg.LoadConfig()

				Expect(err.Error()).To(Equal(`Config File "config" Not Found in "[]"`))
			})
		})
	})

	Context("ParseConfig function", Label("integration"), func() {
		var (
			cfg config.Config
		)

		BeforeEach(func() {
			cfg, _ = config.NewViperConfig()
		})

		When("success parse config", func() {
			It("should return result", func() {
				res := struct{}{}
				cfg.LoadConfig()
				err := cfg.ParseConfig(&res)

				Expect(err).To(BeNil())
			})
		})
	})

})
