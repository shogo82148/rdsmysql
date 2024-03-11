package rdsmysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/internal/testutils"
)

func TestApply(t *testing.T) {
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
	if err := Apply(config, awsConfig); err != nil {
		t.Fatal(err)
	}

	conn, err := mysql.NewConnector(config)
	if err != nil {
		t.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Error(err)
	}
}
