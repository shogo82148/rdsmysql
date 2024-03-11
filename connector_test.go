package rdsmysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/internal/testutils"
)

func TestConnector(t *testing.T) {
	testutils.Setup(t)

	awsConfig := aws.NewConfig().WithRegion(testutils.Region)
	awsSession := session.Must(session.NewSession(awsConfig))

	config := mysql.NewConfig()
	config.User = testutils.User
	config.Addr = testutils.Host
	connector := &Connector{
		Session:           awsSession,
		Config:            config,
		MaxConnsPerSecond: 10,
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		t.Error(err)
	}
}
