package rdsmysql

import (
	"database/sql/driver"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-sql-driver/mysql"
	pkgerrors "github.com/pkg/errors"
)

// Driver is mysql driver using IAM DB Auth.
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
type Driver struct {
	AWSSession *session.Session

	// If role is set, do sts:AssumeRole before getting the token.
	Role string
}

// Open opens new connection.
func (d *Driver) Open(name string) (driver.Conn, error) {
	config, err := mysql.ParseDSN(name)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "fail to parse dns")
	}

	// Get Credentials
	var cred *credentials.Credentials
	if d.Role == "" {
		cred = d.AWSSession.Config.Credentials
	} else {
		stssvc := sts.New(d.AWSSession)
		caller, err := stssvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err != nil {
			return nil, pkgerrors.Wrap(err, "fail to get call identity")
		}
		res, err := stssvc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn: aws.String(
				fmt.Sprintf(
					"arn:aws:iam::%s:role/"+d.Role,
					*caller.Account,
				),
			),
			RoleSessionName: aws.String("connect"),
		})
		if err != nil {
			return nil, pkgerrors.Wrap(err, "fail to assume role")
		}
		cred = credentials.NewStaticCredentials(
			*res.Credentials.AccessKeyId,
			*res.Credentials.SecretAccessKey,
			*res.Credentials.SessionToken,
		)
	}

	token, err := rdsutils.BuildAuthToken(config.Addr, *d.AWSSession.Config.Region, config.User, cred)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "fail to build auth token")
	}

	// override configure
	config.AllowCleartextPasswords = true
	config.Passwd = token
	config.TLSConfig = "rdsmysql"

	return (&mysql.MySQLDriver{}).Open(config.FormatDSN())
}
