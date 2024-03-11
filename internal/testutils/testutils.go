package testutils

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
)

const TestUser = "rdsmysql"

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

	if err := initializeUser(context.Background(), db); err != nil {
		t.Fatal(err)
	}
}

func initializeUser(ctx context.Context, db *sql.DB) error {
	var cnt int
	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM `mysql`.`user` WHERE `user` = ?", TestUser)
	if err := row.Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return nil // already initialized
	}

	_, err := db.ExecContext(ctx, "CREATE USER '"+TestUser+"' IDENTIFIED WITH AWSAuthenticationPlugin AS 'RDS'")
	if err != nil {
		return err
	}

	return nil
}