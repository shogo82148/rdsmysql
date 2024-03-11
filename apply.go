package rdsmysql

import (
	"context"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-sql-driver/mysql"
)

// Apply applies the IAM DB Auth to mysqlConfig.
//
// It overrides the following fields of mysqlConfig:
//   - AllowCleartextPasswords: true
//   - TLS: the certificate of Amazon RDS
//   - Passwd: the auth token
//   - BeforeConnect: refresh the auth token
func Apply(config *mysql.Config, session *session.Session) error {
	// override configure for Amazon RDS
	// see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.AWSCLI.html
	config.AllowCleartextPasswords = true
	config.TLS = TLSConfig.Clone()

	// refresh token
	cred := session.Config.Credentials
	region := session.Config.Region
	if region == nil {
		return fmt.Errorf("rdsmysql: region is missing")
	}
	addr := ensureHavePort(config.Addr)
	beforeConnect := func(ctx context.Context, config *mysql.Config) error {
		token, err := rdsutils.BuildAuthToken(addr, *region, config.User, cred)
		if err != nil {
			return fmt.Errorf("rdsmysql: fail to build auth token: %w", err)
		}
		config.Passwd = token
		return nil
	}

	if err := config.Apply(mysql.BeforeConnect(beforeConnect)); err != nil {
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
