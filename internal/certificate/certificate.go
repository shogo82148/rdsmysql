package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Certificate is the certificates for connecting RDS MySQL with SSL/TLS.
// It contains the intermediate and root certificates for RDS MySQL ( https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem ),
// and the root certifcate for RDS Proxy( https://www.amazontrust.com/repository/AmazonRootCA1.pem ).
const Certificate = rdsProxyCertificate + rdsCertificates

func init() {
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM([]byte(Certificate)); !ok {
		panic(errors.New("failed to append certs"))
	}
	err := mysql.RegisterTLSConfig("rdsmysql", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		panic(err)
	}
}
