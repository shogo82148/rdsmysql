package rdsmysql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-sql-driver/mysql"
	_ "github.com/shogo82148/rdsmysql/internal/certificate" // install certificate.
	"golang.org/x/time/rate"
)

// check Connector implements driver.Connector.
var _ driver.Connector = &Connector{}

// Connector is an implementation of driver.Connector
type Connector struct {
	Session           *session.Session
	Config            *mysql.Config
	MaxConnsPerSecond int

	mu      sync.Mutex
	limiter *rate.Limiter
	config  *mysql.Config
}

// Connect returns a connection to the database.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	// rate limit
	if l := c.getlimiter(); l != nil {
		if err := l.Wait(ctx); err != nil {
			return nil, err
		}
	}

	cred := c.Session.Config.Credentials
	region := c.Session.Config.Region
	if region == nil {
		return nil, errors.New("rdsmysql: region is missing")
	}

	c.mu.Lock()
	if c.config == nil {
		copy := *c.Config // shallow copy, but ok. we rewrite only shallow fields.

		// format and parse dns.
		// because TLS config is loaded by ParseDNS.
		copy.TLSConfig = "rdsmysql"
		config, err := mysql.ParseDSN(copy.FormatDSN())
		if err != nil {
			c.mu.Unlock()
			return nil, fmt.Errorf("fail to parse dsn: %w", err)
		}
		c.config = config
	}
	config := c.config

	token, err := rdsutils.BuildAuthToken(config.Addr, *region, config.User, cred)
	if err != nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("fail to build auth token: %w", err)
	}

	// override configure
	config.AllowCleartextPasswords = true
	config.Passwd = token
	config.TLSConfig = "rdsmysql"

	connector, err := mysql.NewConnector(config)
	if err != nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("fail to created new connector: %w", err)
	}
	c.mu.Unlock()

	return connector.Connect(ctx)
}

func (c *Connector) getlimiter() *rate.Limiter {
	if c.MaxConnsPerSecond == 0 {
		return nil
	}
	c.mu.Lock()
	limiter := c.limiter
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(c.MaxConnsPerSecond), 1)
		c.limiter = limiter
	}
	c.mu.Unlock()
	return limiter
}

// Driver returns the underlying Driver of the Connector.
func (c *Connector) Driver() driver.Driver {
	return &Driver{
		Session: c.Session,
	}
}
