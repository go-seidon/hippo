package grpc_log_test

import (
	"fmt"

	grpc "github.com/go-seidon/local/internal/grpc"
	grpc_log "github.com/go-seidon/local/internal/grpc-log"
	mock_grpclog "github.com/go-seidon/local/internal/grpc-log/mock"
	grpc_test "github.com/go-seidon/local/internal/grpc-test"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stream Package", func() {

	Describe("Log Server Stream", func() {
		var (
			lss    grpc.ServerStream
			ss     *mock_grpclog.MockServerStream
			logger *mock_logging.MockLogger
			msg    *grpc_test.TestData
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			ss = mock_grpclog.NewMockServerStream(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			lss = grpc_log.NewLogServerStream(ss, logger)
			msg = &grpc_test.TestData{}
		})

		Context("SendMsg function", Label("unit"), func() {
			When("error happened", func() {
				It("should return error", func() {
					expectErr := fmt.Errorf("network error")

					ss.
						EXPECT().
						SendMsg(gomock.Eq(msg)).
						Return(expectErr).
						Times(1)

					err := lss.SendMsg(msg)

					Expect(err).To(Equal(expectErr))
				})
			})

			When("error not happened", func() {
				It("should return result", func() {
					ss.
						EXPECT().
						SendMsg(gomock.Eq(msg)).
						Return(nil).
						Times(1)

					logger.
						EXPECT().
						WithFields(gomock.Any()).
						Return(logger).
						Times(1)

					logger.
						EXPECT().
						Info(gomock.Eq("send stream")).
						Times(1)

					err := lss.SendMsg(msg)

					Expect(err).To(BeNil())
				})
			})
		})

		Context("RecvMsg function", Label("unit"), func() {
			When("error happened", func() {
				It("should return error", func() {
					expectErr := fmt.Errorf("network error")

					ss.
						EXPECT().
						RecvMsg(gomock.Eq(msg)).
						Return(expectErr).
						Times(1)

					err := lss.RecvMsg(msg)

					Expect(err).To(Equal(expectErr))
				})
			})

			When("error not happened", func() {
				It("should return result", func() {
					ss.
						EXPECT().
						RecvMsg(gomock.Eq(msg)).
						Return(nil).
						Times(1)

					logger.
						EXPECT().
						WithFields(gomock.Any()).
						Return(logger).
						Times(1)

					logger.
						EXPECT().
						Info(gomock.Eq("receive stream")).
						Times(1)

					err := lss.RecvMsg(msg)

					Expect(err).To(BeNil())
				})
			})
		})
	})

})
