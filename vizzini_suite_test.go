package vizzini_test

import (
	"flag"
	"fmt"
	"log"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"

	"github.com/cloudfoundry-incubator/receptor"
)

var client receptor.Client
var domain string
var stack string

var receptorAddress, receptorUsername, receptorPassword string

func init() {
	flag.StringVar(&receptorAddress, "receptor-address", "receptor.10.244.0.34.xip.io", "http address for the receptor (required)")
	flag.StringVar(&receptorUsername, "receptor-username", "", "receptor username")
	flag.StringVar(&receptorUsername, "receptor-password", "", "receptor password")
	flag.Parse()

	if receptorAddress == "" {
		log.Fatal("i need a receptor-address to talk to Diego...")
	}
}

func TestReceptorSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReceptorSuite Suite")
}

func NewGuid() string {
	u, err := uuid.NewV4()
	Ω(err).ShouldNot(HaveOccurred())
	return u.String()
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(10 * time.Second)
	domain = fmt.Sprintf("vizzini-%d", GinkgoParallelNode())
	stack = "lucid64"

	client = receptor.NewClient(receptorAddress, receptorUsername, receptorPassword)
})

var _ = AfterSuite(func() {
	Ω(client.GetAllTasksByDomain(domain)).Should(BeEmpty())
	Ω(client.GetAllDesiredLRPsByDomain(domain)).Should(BeEmpty())
})
