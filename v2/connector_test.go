package rdsmysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/v2/internal/testutils"
)

func TestConnector(t *testing.T) {
	testutils.Setup(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(testutils.Region))
	if err != nil {
		t.Fatal(err)
	}

	config := mysql.NewConfig()
	config.User = testutils.User
	config.Addr = testutils.Host
	connector := &Connector{
		AWSConfig:         awsConfig,
		MySQLConfig:       config,
		MaxConnsPerSecond: 10,
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Error(err)
	}
}
