package grpclog_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrpcLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc Log Package")
}
