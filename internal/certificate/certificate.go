package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Certificate is the certificates for connecting RDS MySQL with SSL/TLS.
// It contains the intermediate and root certificates for RDS MySQL ( https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem ),
// and the root certificates for RDS Proxy( https://www.amazontrust.com/repository/AmazonRootCA1.pem ).
const Certificate = rdsProxyCertificate + rdsCertificates

var Config *tls.Config

func init() {
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM([]byte(Certificate)); !ok {
		panic(errors.New("failed to append certs"))
	}
	Config = &tls.Config{
		RootCAs: rootCertPool,
	}
	err := mysql.RegisterTLSConfig("rdsmysql", Config)
	if err != nil {
		panic(err)
	}
}
