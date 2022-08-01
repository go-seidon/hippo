package logging_test

import (
	"bytes"
	"fmt"

	"github.com/go-seidon/local/internal/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Logrus Formater Package", func() {

	Context("Format function", Label("unit"), func() {
		var (
			formatter *logging.GoFormatter
			entry     *logrus.Entry
		)

		BeforeEach(func() {
			formatter = &logging.GoFormatter{}
			entry = &logrus.Entry{}
		})

		When("success format message", func() {
			It("should return result", func() {
				res, err := formatter.Format(entry)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("data contain error", func() {
			It("should return result", func() {
				entry.Data = logrus.Fields{
					"error": fmt.Errorf("some error"),
				}
				res, err := formatter.Format(entry)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("data contain field", func() {
			It("should return result", func() {
				entry.Data = logrus.Fields{
					"key": "value",
				}
				res, err := formatter.Format(entry)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("buffer is specified", func() {
			It("should return result", func() {
				entry.Buffer = &bytes.Buffer{}
				res, err := formatter.Format(entry)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("pretty print is used", func() {
			It("should return result", func() {
				formatter.PrettyPrint = true
				res, err := formatter.Format(entry)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

})
