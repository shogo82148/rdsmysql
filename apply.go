package rdsmysql

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
)

// ApplyOption configures the behavior of [Apply].
type ApplyOption func(*applyOptions)

type applyOptions struct {
	useSystemCertPool bool
}

// WithSystemCertPool makes [Apply] trust the system's CA certificate pool instead of the
// Amazon RDS root certificates. Enable this when connecting through [Amazon RDS Proxy], which
// presents certificates issued by AWS Certificate Manager (ACM) rather than the Amazon RDS
// root CA.
//
// [Amazon RDS Proxy]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html
func WithSystemCertPool() ApplyOption {
	return func(o *applyOptions) {
		o.useSystemCertPool = true
	}
}

// Apply applies the IAM DB Auth to mysqlConfig.
//
// It overrides the following fields of mysqlConfig:
//   - AllowCleartextPasswords: true
//   - TLS: the certificate of Amazon RDS (see [WithSystemCertPool] to trust the system's CA pool instead)
//   - Passwd: the auth token
//   - BeforeConnect: refresh the auth token
func Apply(mysqlConfig *mysql.Config, awsConfig aws.Config, opts ...ApplyOption) error {
	var options applyOptions
	for _, opt := range opts {
		opt(&options)
	}

	tlsConfig := TLSConfig
	if options.useSystemCertPool {
		config, err := NewTLSConfig(true)
		if err != nil {
			return fmt.Errorf("rdsmysql: fail to build TLS config: %w", err)
		}
		tlsConfig = config
	}

	// override configure for Amazon RDS
	// see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.AWSCLI.html
	mysqlConfig.AllowCleartextPasswords = true
	mysqlConfig.TLS = tlsConfig.Clone()

	// refresh token
	cred := awsConfig.Credentials
	region := awsConfig.Region
	if region == "" {
		return errors.New("rdsmysql: region is missing")
	}
	addr := ensureHavePort(mysqlConfig.Addr)
	beforeConnect := func(ctx context.Context, config *mysql.Config) error {
		token, err := auth.BuildAuthToken(ctx, addr, region, config.User, cred)
		if err != nil {
			return fmt.Errorf("rdsmysql: fail to build auth token: %w", err)
		}
		config.Passwd = token
		return nil
	}

	if err := mysqlConfig.Apply(mysql.BeforeConnect(beforeConnect)); err != nil {
		return fmt.Errorf("rdsmysql: fail to apply beforeConnect: %w", err)
	}
	return nil
}

// ensureHavePort ensures that addr has a port.
func ensureHavePort(addr string) string {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return net.JoinHostPort(addr, "3306")
	}
	return addr
}
