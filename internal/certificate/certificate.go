//go:generate go run ../cmd/update_certificate/main.go

package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Certificate is the certificates for connecting RDS MySQL with SSL/TLS.
// It contains the intermediate and root certificates for [Amazon RDS MySQL] and [Amazon Aurora MySQL].
//
// [Amazon RDS MySQL]: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html
// [Amazon Aurora MySQL]: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.SSL.html
const Certificate = rdsCertificates

// Config is the tls.Config for connecting RDS MySQL with SSL/TLS.
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
