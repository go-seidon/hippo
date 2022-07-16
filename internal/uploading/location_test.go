package uploading_test

import (
	"time"

	"github.com/go-seidon/local/internal/mock"
	"github.com/go-seidon/local/internal/uploading"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Daily Rotate Service", func() {
	Context("NewDailyRotate function", Label("unit"), func() {
		var (
			p     uploading.NewDailyRotateParam
			clock *mock.MockClock
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			clock = mock.NewMockClock(ctrl)
			p = uploading.NewDailyRotateParam{}
		})

		When("success create rotator", func() {
			It("should return result", func() {
				res := uploading.NewDailyRotate(p)

				Expect(res).ToNot(BeNil())
			})
		})

		When("clock is specified", func() {
			It("should return result", func() {
				p.Clock = clock
				res := uploading.NewDailyRotate(p)

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("GetLocation function", Label("unit"), func() {
		var (
			s     uploading.UploadLocation
			clock *mock.MockClock
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			clock = mock.NewMockClock(ctrl)

			s = uploading.NewDailyRotate(uploading.NewDailyRotateParam{
				Clock: clock,
			})
		})

		When("function is called", func() {
			It("should return result", func() {
				currentTimestamp, _ := time.Parse("2006-01-02", "2022-02-28")
				clock.
					EXPECT().
					Now().
					Return(currentTimestamp).
					Times(1)

				res := s.GetLocation()

				Expect(res).To(Equal("2022/02/28"))
			})
		})

		When("less than 2 digit value", func() {
			It("should return result", func() {
				currentTimestamp, _ := time.Parse("2006-1-2", "1990-2-1")
				clock.
					EXPECT().
					Now().
					Return(currentTimestamp).
					Times(1)

				res := s.GetLocation()

				Expect(res).To(Equal("1990/02/01"))
			})
		})
	})
})
