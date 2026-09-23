package rdsmysql

import (
	"database/sql"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/v2/internal/testutils"
)

func TestApply(t *testing.T) {
	testutils.Setup(t)

	ctx := t.Context()

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
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := db.PingContext(ctx); err != nil {
		t.Error(err)
	}
}
