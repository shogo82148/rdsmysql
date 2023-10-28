// Package rdsmysql is a MySQL SQL driver that allows [IAM Database Authentication for Amazon RDS]
// and [IAM Database Authentication for Amazon Aurora].
// It also supports connecting with [the RDS proxy using IAM authentication].
//
// [IAM Database Authentication for Amazon RDS]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
// [IAM Database Authentication for Amazon Aurora]: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html
// [the RDS proxy using IAM authentication]: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy.html#rds-proxy-connecting-iam
package rdsmysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
)

// Driver is a MySQL driver using IAM DB Auth.
type Driver struct {
	// Session is an AWS session
	Session *session.Session
}

var _ driver.Driver = (*Driver)(nil)
var _ driver.DriverContext = (*Driver)(nil)

// Open opens a new connection.
func (d *Driver) Open(name string) (driver.Conn, error) {
	c, err := d.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return c.Connect(context.Background())
}

// OpenConnector opens a new connector.
func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	if d.Session.Config.Region == nil {
		return nil, errors.New("rdsmysql: region is missing")
	}

	config, err := mysql.ParseDSN(name)
	if err != nil {
		return nil, fmt.Errorf("rdsmysql: fail to parse dns: %w", err)
	}

	return &Connector{
		Session: d.Session,
		Config:  config,
	}, nil
}
