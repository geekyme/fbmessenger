package fbmessenger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo/reporters"
	"testing"

	"os"
	"path/filepath"
)

func TestWebhooks(t *testing.T) {
	RegisterFailHandler(Fail)
	testReportsPath, _ := filepath.Abs("./test-reports")
	os.MkdirAll(testReportsPath, 0777)
	junitReporter := reporters.NewJUnitReporter(filepath.Join(testReportsPath, "fbmessenger-junit.xml"))
	RunSpecsWithDefaultAndCustomReporters(t, "FBMessenger Suite", []Reporter{junitReporter})
}
