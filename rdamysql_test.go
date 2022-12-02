package rdsmysql_test

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/shogo82148/rdsmysql"
)

func ExampleOpen() {
	// register authentication information
	c := aws.NewConfig().WithRegion("ap-northeast-1")
	s := session.Must(session.NewSession(c))
	d := &rdsmysql.Driver{
		Session: s,
	}
	sql.Register("rdsmysql", d)

	db, err := sql.Open("rdsmysql", "user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// do something with db
}
