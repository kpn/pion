package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	// tests to run
	_ "github.com/kpn/pion/test/e2e/bucket"
	_ "github.com/kpn/pion/test/e2e/file"
)

func RunE2ETests(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t, "pion e2e suite", []Reporter{reporters.NewJUnitReporter("it_cov.xml")})
}
