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

func TestApply(t *testing.T) {
	testutils.Setup(t)

	awsConfig := aws.NewConfig().WithRegion(testutils.Region)
	awsSession := session.Must(session.NewSession(awsConfig))

	config := mysql.NewConfig()
	config.User = testutils.User
	config.Addr = testutils.Host
	if err := Apply(config, awsSession); err != nil {
		t.Fatal(err)
	}

	conn, err := mysql.NewConnector(config)
	if err != nil {
		t.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		t.Error(err)
	}
}
