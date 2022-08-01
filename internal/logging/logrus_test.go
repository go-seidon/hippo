package logging_test

import (
	"context"
	"fmt"

	"github.com/go-seidon/local/internal/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logrus Package", func() {

	Context("NewLogrusLog function", Label("unit"), func() {
		When("option is not specified", func() {
			It("should return result", func() {
				res := logging.NewLogrusLog()

				Expect(res).ToNot(BeNil())
			})
		})

		When("option is specified", func() {
			It("should return result", func() {
				opt := logging.WithAppContext("mock-name", "mock-version")
				res := logging.NewLogrusLog(opt)

				Expect(res).ToNot(BeNil())
			})
		})

		When("debugging is enabled", func() {
			It("should return result", func() {
				opt1 := logging.WithAppContext("mock-name", "mock-version")
				opt2 := logging.EnableDebugging()
				res := logging.NewLogrusLog(opt1, opt2)

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("Info function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Info("mock-log")

			})
		})
	})

	Context("Debug function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Debug("mock-log")

			})
		})
	})

	Context("Error function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Error("mock-log")

			})
		})
	})

	Context("Warn function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Warn("mock-log")

			})
		})
	})

	Context("Infof function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Infof("%s", "mock-log")

			})
		})
	})

	Context("Debugf function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Debugf("%s", "mock-log")

			})
		})
	})

	Context("Errorf function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Errorf("%s", "mock-log")

			})
		})
	})

	Context("Warnf function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Warnf("%s", "mock-log")

			})
		})
	})

	Context("Infoln function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Infoln("mock-log")

			})
		})
	})

	Context("Debugln function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Debugln("mock-log")

			})
		})
	})

	Context("Errorln function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Errorln("mock-log")

			})
		})
	})

	Context("Warnln function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				logger.Warnln("mock-log")

			})
		})
	})

	Context("WithFields function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				res := logger.WithFields(map[string]interface{}{})

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("WithError function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				res := logger.WithError(fmt.Errorf("some error"))

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("WithContext function", Label("unit"), func() {
		var (
			logger logging.Logger
		)

		BeforeEach(func() {
			logger = logging.NewLogrusLog()
		})

		When("success send log", func() {
			It("should return nil", func() {
				res := logger.WithContext(context.Background())

				Expect(res).ToNot(BeNil())
			})
		})
	})

})
