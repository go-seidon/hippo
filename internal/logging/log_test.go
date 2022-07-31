package logging_test

import (
	"testing"

	"github.com/go-seidon/local/internal/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Package")
}

var _ = Describe("Logging Package", func() {

	Context("WithAppContext function", Label("unit"), func() {
		When("parameter is specified", func() {
			It("should return result", func() {
				opt := logging.WithAppContext("mock-name", "mock-version")
				var res logging.LogOption
				opt(&res)

				Expect(res.AppName).To(Equal("mock-name"))
				Expect(res.AppVersion).To(Equal("mock-version"))
				Expect(res.DebuggingEnabled).To(BeFalse())
				Expect(res.PrettyPrintEnabled).To(BeFalse())
				Expect(res.StackSkip).To(BeNil())
			})
		})
	})

	Context("EnableDebugging function", Label("unit"), func() {
		When("function is called", func() {
			It("should return result", func() {
				opt := logging.EnableDebugging()
				var res logging.LogOption
				opt(&res)

				Expect(res.DebuggingEnabled).To(BeTrue())
				Expect(res.AppName).To(Equal(""))
				Expect(res.AppVersion).To(Equal(""))
				Expect(res.PrettyPrintEnabled).To(BeFalse())
				Expect(res.StackSkip).To(BeNil())
			})
		})
	})

	Context("EnablePrettyPrint function", Label("unit"), func() {
		When("function is called", func() {
			It("should return result", func() {
				opt := logging.EnablePrettyPrint()
				var res logging.LogOption
				opt(&res)

				Expect(res.PrettyPrintEnabled).To(BeTrue())
				Expect(res.DebuggingEnabled).To(BeFalse())
				Expect(res.AppName).To(Equal(""))
				Expect(res.AppVersion).To(Equal(""))
				Expect(res.StackSkip).To(BeNil())
			})
		})
	})

	Context("AddStackSkip function", Label("unit"), func() {
		When("add one stack skip", func() {
			It("should return result", func() {
				opt := logging.AddStackSkip("some-pkg")
				var res logging.LogOption
				opt(&res)

				Expect(res.StackSkip).To(Equal([]string{
					"some-pkg",
				}))
				Expect(res.PrettyPrintEnabled).To(BeFalse())
				Expect(res.DebuggingEnabled).To(BeFalse())
				Expect(res.AppName).To(Equal(""))
				Expect(res.AppVersion).To(Equal(""))
			})
		})

		When("add two stack skip", func() {
			It("should return result", func() {
				opt1 := logging.AddStackSkip("some-pkg-1")
				opt2 := logging.AddStackSkip("some-pkg-2")
				var res logging.LogOption
				opt1(&res)
				opt2(&res)

				Expect(res.StackSkip).To(Equal([]string{
					"some-pkg-1",
					"some-pkg-2",
				}))
				Expect(res.PrettyPrintEnabled).To(BeFalse())
				Expect(res.DebuggingEnabled).To(BeFalse())
				Expect(res.AppName).To(Equal(""))
				Expect(res.AppVersion).To(Equal(""))
			})
		})
	})

})
