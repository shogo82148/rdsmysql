package rdsmysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/v2/internal/certificate"
	"golang.org/x/time/rate"
)

// check Connector implements driver.Connector.
var _ driver.Connector = (*Connector)(nil)

// Connector is an implementation of driver.Connector
type Connector struct {
	// AWSConfig is AWS Config.
	AWSConfig *aws.Config

	// MySQLConfig is a configure for connecting to MySQL servers.
	MySQLConfig *mysql.Config

	// MaxConnsPerSecond is a limit for creating new connections.
	// Zero means no limit.
	MaxConnsPerSecond int

	mu      sync.Mutex
	limiter *rate.Limiter
	// config is same as Config, but TLS configured
	config *mysql.Config
}

// Connect returns a connection to the database.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	// rate limit
	if l := c.getlimiter(); l != nil {
		if err := l.Wait(ctx); err != nil {
			return nil, err
		}
	}

	connector, err := c.newConnector(ctx)
	if err != nil {
		return nil, err
	}
	return connector.Connect(ctx)
}

func (c *Connector) newConnector(ctx context.Context) (driver.Connector, error) {
	config, err := c.newConfig()
	if err != nil {
		return nil, err
	}

	// refresh token
	cred := c.AWSConfig.Credentials
	region := c.AWSConfig.Region
	if region == "" {
		return nil, errors.New("rdsmysql: region is missing")
	}
	token, err := auth.BuildAuthToken(ctx, config.Addr, region, config.User, cred)
	if err != nil {
		return nil, fmt.Errorf("rdsmysql: fail to build auth token: %w", err)
	}
	config.Passwd = token

	// create new connector
	connector, err := mysql.NewConnector(config)
	if err != nil {
		return nil, fmt.Errorf("rdsmysql: fail to created new connector: %w", err)
	}
	return connector, nil
}

func (c *Connector) newConfig() (*mysql.Config, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.config == nil {
		clone := *c.MySQLConfig // shallow copy, but ok. we rewrite only shallow fields.

		// override configure for Amazon RDS
		// see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.AWSCLI.html
		clone.AllowCleartextPasswords = true
		clone.TLS = certificate.Config

		c.config = &clone
	}

	clone := *c.config // shallow copy, but ok. we rewrite only shallow fields.
	return &clone, nil
}

func (c *Connector) getlimiter() *rate.Limiter {
	if c.MaxConnsPerSecond == 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	limiter := c.limiter
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(c.MaxConnsPerSecond), 1)
		c.limiter = limiter
	}
	return limiter
}

// Driver returns the underlying Driver of the Connector.
func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}
