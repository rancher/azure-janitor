package action_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAzureJanitorAction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Azure Janitor Suite")
}
