package grpchandler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrpcHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc Handler Package")
}
