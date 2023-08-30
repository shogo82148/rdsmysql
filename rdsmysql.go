// Package rdsmysql is a SQL driver that allows IAM Database Authentication for Amazon RDS and Amazon Aurora.
// It also supports connecting to the RDS proxy using IAM authentication.
//
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html
// https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy.html#rds-proxy-connecting-iam
package rdsmysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
)

// Driver is a mysql driver using IAM DB Auth.
//
//	// configure AWS session
//	awsConfig := aws.NewConfig().WithRegion("ap-northeast-1")
//	awsSession := session.Must(session.NewSession(awsConfig))
//
//	// configure the driver
//	driver := &rdsmysql.Driver{
//	  Session: awsSession,
//	}
//	sql.Register("rdsmysql", driver)
//
//	// additional code for using the `rdsmysql` driver
type Driver struct {
	// Session is an AWS session
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
	if d.Session.Config.Region == nil {
		return nil, errors.New("rdsmysql: region is missing")
	}

	config, err := mysql.ParseDSN(name)
	if err != nil {
		return nil, fmt.Errorf("fail to parse dns: %w", err)
	}

	return &Connector{
		Session: d.Session,
		Config:  config,
	}, nil
}
