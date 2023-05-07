package certificate

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed" // for go:embed directive
	"errors"

	"github.com/go-sql-driver/mysql"
)

// rdsProxyCertificate is the root certificate for RDS Proxy.
// It comes from https://www.amazontrust.com/repository/AmazonRootCA1.pem.
// See https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy.html#rds-proxy-connecting-iam for more details.
//
//go:embed rds_proxy.pem
var rdsProxyCertificate string

// rdsCertificates the intermediate and root certificates for RDS MySQL.
// It comes from https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem.
// See https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html#UsingWithRDS.SSL.CertificatesAllRegions for more details.
//
//go:embed rds.pem
var rdsCertificates string

// Certificate is the certificates for connecting RDS MySQL with SSL/TLS.
// It contains the intermediate and root certificates for RDS MySQL ( https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem ),
// and the root certificates for RDS Proxy( https://www.amazontrust.com/repository/AmazonRootCA1.pem ).
var Certificate []byte

// Config is the tls.Config for connecting RDS MySQL with SSL/TLS.
var Config *tls.Config

func init() {
	Certificate = []byte(rdsProxyCertificate + rdsCertificates)
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(Certificate); !ok {
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
