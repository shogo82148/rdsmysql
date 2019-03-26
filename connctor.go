package rdsmysql

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-sql-driver/mysql"
	_ "github.com/shogo82148/rdsmysql/internal/certificate" // install certificate.
	"golang.org/x/xerrors"
)

// check Connector implements driver.Connctor.
var _ driver.Connector = &Connector{}

// Connector is an implementation of driver.Connector
type Connector struct {
	Session *session.Session
	Config  *mysql.Config
}

// Connect returns a connection to the database.
func (c *Connector) Connect(context.Context) (driver.Conn, error) {
	config := *c.Config // shollow copy

	cred := c.Session.Config.Credentials
	region := c.Session.Config.Region
	if region == nil {
		return nil, errors.New("rdsmysql: region is missing")
	}
	token, err := rdsutils.BuildAuthToken(config.Addr, *region, config.User, cred)
	if err != nil {
		return nil, xerrors.Errorf("fail to build auth token: %w", err)
	}

	// override configure
	config.AllowCleartextPasswords = true
	config.Passwd = token
	config.TLSConfig = "rdsmysql"

	return (&mysql.MySQLDriver{}).Open(config.FormatDSN())
}

// Driver returns the underlying Driver of the Connector.
func (c *Connector) Driver() driver.Driver {
	return &Driver{
		Session: c.Session,
	}
}
