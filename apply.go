package rdsmysql

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/internal/certificate"
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
	config.TLS = certificate.Config

	// refresh token
	cred := session.Config.Credentials
	region := session.Config.Region
	if region == nil {
		return fmt.Errorf("rdsmysql: region is missing")
	}
	beforeConnect := func(ctx context.Context, config *mysql.Config) error {
		token, err := rdsutils.BuildAuthToken(config.Addr, *region, config.User, cred)
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
