//go:generate go run internal/cmd/update_certificate/main.go

package rdsmysql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Certificates is the certificates for connecting RDS MySQL with SSL/TLS.
// It contains the intermediate and root certificates for [Amazon RDS MySQL] and [Amazon Aurora MySQL].
//
// [Amazon RDS MySQL]: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html
// [Amazon Aurora MySQL]: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.SSL.html
const Certificates = rdsCertificates

// TLSConfig is the tls.TLSConfig for connecting RDS MySQL with SSL/TLS.
var TLSConfig *tls.Config

func init() {
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM([]byte(Certificates)); !ok {
		panic(errors.New("failed to append certs"))
	}
	TLSConfig = &tls.Config{
		RootCAs: rootCertPool,
	}
	err := mysql.RegisterTLSConfig("rdsmysql", TLSConfig)
	if err != nil {
		panic(err)
	}
}
