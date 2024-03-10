package rdsmysql

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestCertificate(t *testing.T) {
	t.Parallel()

	pemCerts := []byte(Certificates)
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			t.Errorf("unexpected block type: %s", block.Type)
			continue
		}
		if _, err := x509.ParseCertificate(block.Bytes); err != nil {
			t.Errorf("failed to parse certificate: %v", err)
		}
	}
}
