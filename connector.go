package rdsmysql

import (
	"context"
	"database/sql/driver"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/time/rate"
)

// check Connector implements driver.Connector.
var _ driver.Connector = (*Connector)(nil)

// Connector is a MySQL connector using IAM DB Auth.
// It implements [database/sql/driver.Connector].
type Connector struct {
	// Session is AWS Session.
	Session *session.Session

	// Config is a configure for connecting to MySQL servers.
	Config *mysql.Config

	// MaxConnsPerSecond is a limit for creating new connections.
	// Zero means no limit.
	MaxConnsPerSecond int

	// once guards config and limiter
	once sync.Once

	limiter *rate.Limiter

	// connector is the IAM DB Auth configured MySQL connector.
	connector driver.Connector

	// err is an error occurred during initialization
	err error
}

func (c *Connector) init() {
	// shallow copy, but ok. we rewrite only shallow fields.
	config := new(mysql.Config)
	*config = *c.Config

	// override configure for Amazon RDS
	if err := Apply(config, c.Session); err != nil {
		c.err = err
		return
	}
	connector, err := mysql.NewConnector(config)
	if err != nil {
		c.err = err
		return
	}
	c.connector = connector

	// create limiter
	if c.MaxConnsPerSecond > 0 {
		c.limiter = rate.NewLimiter(rate.Limit(c.MaxConnsPerSecond), 1)
	}
}

// Connect returns a connection to the database.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	// initialize
	c.once.Do(c.init)
	if c.err != nil {
		return nil, c.err
	}

	// rate limit
	if l := c.limiter; l != nil {
		if err := l.Wait(ctx); err != nil {
			return nil, err
		}
	}

	return c.connector.Connect(ctx)
}

// Driver returns the underlying [database/sql/driver.Driver] of the [Connector].
func (c *Connector) Driver() driver.Driver {
	return &Driver{
		Session: c.Session,
	}
}
