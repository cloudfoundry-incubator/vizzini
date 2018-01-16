package vizzini_test

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"code.cloudfoundry.org/bbs/models"
	. "code.cloudfoundry.org/vizzini/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TLS Proxy", func() {
	var lrp *models.DesiredLRP

	BeforeEach(func() {
		if !enableContainerProxyTests {
			Skip("container proxy tests are disabled")
		}

		lrp = DesiredLRPWithGuid(guid)
		Expect(bbsClient.DesireLRP(logger, lrp)).To(Succeed())
		Eventually(ActualGetter(logger, guid, 0)).Should(BeActualLRPWithState(guid, 0, models.ActualLRPStateRunning))
	})

	It("proxies traffic to the application process inside the container", func() {
		directURL := "https://" + TLSDirectAddressFor(guid, 0, 8080)

		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// the instance cred cert uses the internal container ip which causes
					// the request to fail since we can only talk to the host ip adddress
					InsecureSkipVerify: true,
				},
			},
		}
		resp, err := client.Get(directURL)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	Describe("has a valid certificate", func() {
		var (
			certs []*x509.Certificate
			lrp   models.ActualLRP
		)

		BeforeEach(func() {
			conn, err := tls.Dial("tcp", TLSDirectAddressFor(guid, 0, 8080), &tls.Config{
				// the instance cred cert uses the internal container ip which causes the
				// request to fail since we can only talk to the host ip adddress. ignore
				// the cert verification and do some manual assertion on the cert
				// contents
				InsecureSkipVerify: true,
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Handshake()
			Expect(err).NotTo(HaveOccurred())

			connState := conn.ConnectionState()
			Expect(connState.HandshakeComplete).To(BeTrue())
			certs = connState.PeerCertificates
			Expect(certs).To(HaveLen(2)) // the instance identity cert + CA
			lrp, err = ActualGetter(logger, guid, 0)()
			Expect(err).NotTo(HaveOccurred())
		})

		It("has a common name that matches the instance guid", func() {
			Expect(certs[0].Subject.CommonName).To(Equal(lrp.InstanceGuid))
		})
	})
})