package rdsmysql_test

import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/v2"
)

func ExampleConnector() {
	// configure AWS SDK
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		panic(err)
	}

	// configure the connector
	mysqlConfig, err := mysql.ParseDSN("user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
	if err != nil {
		panic(err)
	}
	connector := &rdsmysql.Connector{
		AWSConfig:   awsConfig,
		MySQLConfig: mysqlConfig,
	}

	// open the database
	db := sql.OpenDB(connector)
	defer db.Close()

	// ... do something using db ...
}

func ExampleDriver() {
	// configure AWS SDK
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		panic(err)
	}

	driver := &rdsmysql.Driver{
		AWSConfig: awsConfig,
	}
	sql.Register("rdsmysql", driver)

	db, err := sql.Open("rdsmysql", "user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ... do something using db ...
}
