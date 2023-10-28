package rdsmysql_test

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql"
)

func ExampleConnector() {
	// configure AWS session
	awsConfig := aws.NewConfig().WithRegion("ap-northeast-1")
	awsSession := session.Must(session.NewSession(awsConfig))

	// configure the connector
	cfg, err := mysql.ParseDSN("user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
	if err != nil {
		panic(err)
	}
	connector := &rdsmysql.Connector{
		Session: awsSession,
		Config:  cfg,
	}

	// open the database
	db := sql.OpenDB(connector)
	defer db.Close()

	// ... do something using db ...
}

func ExampleDriver() {
	// register authentication information
	awsConfig := aws.NewConfig().WithRegion("ap-northeast-1")
	awsSession := session.Must(session.NewSession(awsConfig))
	driver := &rdsmysql.Driver{
		Session: awsSession,
	}
	sql.Register("rdsmysql", driver)

	db, err := sql.Open("rdsmysql", "user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ... do something using db ...
}
