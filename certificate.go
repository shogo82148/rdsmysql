//go:generate go run internal/cmd/update_certificate/main.go

package rdsmysql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// Certificates is the certificates for connecting RDS MySQL with SSL/TLS.
// It contains the intermediate and root certificates for [Amazon RDS MySQL] and [Amazon Aurora MySQL].
//
// [Amazon RDS MySQL]: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html
// [Amazon Aurora MySQL]: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.SSL.html
const Certificates = rdsCertificates + rdsGovCloudCertificates

// TLSConfig is the tls.TLSConfig for connecting RDS MySQL with SSL/TLS.
//
// It trusts only the Amazon RDS root certificates in [Certificates]. [Amazon RDS Proxy]
// presents certificates issued by AWS Certificate Manager (ACM) instead, which this
// configuration does not trust; use [NewTLSConfig] with useSystemCertPool set to true to
// trust the system's CA certificates instead.
//
// [Amazon RDS Proxy]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
var TLSConfig *tls.Config

func init() {
	tlsConfig, err := NewTLSConfig(false)
	if err != nil {
		panic(err)
	}
	TLSConfig = tlsConfig
	if err := mysql.RegisterTLSConfig("rdsmysql", TLSConfig); err != nil {
		panic(err)
	}
}

// NewTLSConfig returns a new [tls.Config] for connecting to Amazon RDS MySQL.
//
// If useSystemCertPool is false, it trusts only the Amazon RDS root certificates in
// [Certificates]. If useSystemCertPool is true, it trusts the system's CA certificates
// instead of [Certificates]; enable this when connecting through [Amazon RDS Proxy], which
// presents certificates issued by AWS Certificate Manager (ACM) rather than the Amazon RDS
// root CA.
//
// [Amazon RDS Proxy]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
func NewTLSConfig(useSystemCertPool bool) (*tls.Config, error) {
	if useSystemCertPool {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("rdsmysql: failed to load the system cert pool: %w", err)
		}
		return &tls.Config{RootCAs: pool}, nil
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM([]byte(Certificates)); !ok {
		return nil, errors.New("rdsmysql: failed to append certs")
	}
	return &tls.Config{RootCAs: rootCertPool}, nil
}
