package testutils

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func Setup(t *testing.T) {
	t.Helper()

	user := os.Getenv("RDSMYSQL_USER")
	password := os.Getenv("RDSMYSQL_PASSWORD")
	host := os.Getenv("RDSMYSQL_HOST")
	if host == "" {
		t.Skip("RDSMYSQL_HOST is not set; skip integrated test")
		return
	}
	config := mysql.NewConfig()
	config.User = user
	config.Passwd = password
	config.Net = "tcp"
	config.Addr = host

	conn, err := mysql.NewConnector(config)
	if err != nil {
		t.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		t.Fatal(err)
	}
}
