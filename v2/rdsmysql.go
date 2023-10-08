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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-sql-driver/mysql"
)

// Driver is a mysql driver using IAM DB Auth.
type Driver struct {
	// AWSConfig is AWS Config.
	AWSConfig *aws.Config
}

var _ driver.Driver = (*Driver)(nil)
var _ driver.DriverContext = (*Driver)(nil)

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
	if d.AWSConfig.Region == "" {
		return nil, errors.New("rdsmysql: region is missing")
	}

	config, err := mysql.ParseDSN(name)
	if err != nil {
		return nil, fmt.Errorf("fail to parse dns: %w", err)
	}

	return &Connector{
		AWSConfig:   d.AWSConfig,
		MySQLConfig: config,
	}, nil
}
