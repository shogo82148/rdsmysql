package rdsmysql

import (
	"crypto/x509"
	"testing"
)

func TestCertificate(t *testing.T) {
	t.Parallel()

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(Certificate)); !ok {
		t.Error("failed to append certs")
	}
}
