package rdsmysql

import (
	"database/sql/driver"
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-sql-driver/mysql"
	pkgerrors "github.com/pkg/errors"
)

// Driver is mysql driver using IAM DB Auth.
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
type Driver struct {
	Session *session.Session
}

// Open opens new connection.
func (d *Driver) Open(name string) (driver.Conn, error) {
	config, err := mysql.ParseDSN(name)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "fail to parse dns")
	}

	cred := d.Session.Config.Credentials
	if d.Session.Config.Region == nil {
		return nil, errors.New("rdsmysql: region is missing")
	}

	token, err := rdsutils.BuildAuthToken(config.Addr, *d.Session.Config.Region, config.User, cred)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "fail to build auth token")
	}

	// override configure
	config.AllowCleartextPasswords = true
	config.Passwd = token
	config.TLSConfig = "rdsmysql"

	return (&mysql.MySQLDriver{}).Open(config.FormatDSN())
}
