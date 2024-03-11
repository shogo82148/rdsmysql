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

// Apply applies the IAM DB Auth to mysqlConfig.
//
// It overrides the following fields of mysqlConfig:
//   - AllowCleartextPasswords: true
//   - TLS: the certificate of Amazon RDS
//   - Passwd: the auth token
//   - BeforeConnect: refresh the auth token
func Apply(mysqlConfig *mysql.Config, awsConfig aws.Config) error {
	// override configure for Amazon RDS
	// see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.AWSCLI.html
	mysqlConfig.AllowCleartextPasswords = true
	mysqlConfig.TLS = TLSConfig

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
