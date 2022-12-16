package grpclog_test

import (
	"github.com/go-seidon/hippo/internal/grpclog"
	"github.com/go-seidon/provider/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message Package", func() {

	Context("MarshallJSON function", Label("unit"), func() {
		When("success marshall message", func() {
			It("should return result", func() {
				data := &testdata.SimpleMessage{}
				msg := grpclog.NewMessage(data)

				res, err := msg.MarshalJSON()

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

})
