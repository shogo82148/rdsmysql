// Package rdsmysql is a SQL driver that allows IAM Database Authentication for Amazon RDS and Amazon Aurora.
//
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html
package rdsmysql

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
)

// Driver is a mysql driver using IAM DB Auth.
//
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html
type Driver struct {
	Session *session.Session
}

var _ driver.Driver = &Driver{}
var _ driver.DriverContext = &Driver{}

// Open opens new connection.
func (d *Driver) Open(name string) (driver.Conn, error) {
	c, err := d.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return c.Connect(context.Background())
}

// OpenConnector opens new connection.
func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	config, err := mysql.ParseDSN(name)
	if err != nil {
		return nil, fmt.Errorf("fail to parse dns: %w", err)
	}

	return &Connector{
		Session: d.Session,
		Config:  config,
	}, nil
}
