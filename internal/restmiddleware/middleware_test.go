package restmiddleware_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRestMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest Middleware Package")
}
