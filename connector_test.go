package rdsmysql

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
)

func newTestConnector() *Connector {
	cred := credentials.NewStaticCredentials("AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", "")
	awsConfig := aws.NewConfig().WithRegion("ap-northeast-1").WithCredentials(cred)
	awsSession := session.Must(session.NewSession(awsConfig))
	mysqlConfig, err := mysql.ParseDSN("user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
	if err != nil {
		panic(err)
	}
	return &Connector{
		Session: awsSession,
		Config:  mysqlConfig,
	}
}

func TestNewConnector(t *testing.T) {
	_, err := newTestConnector().newConnector()
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkNewConnector(b *testing.B) {
	c := newTestConnector()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.newConnector()
		}
	})
}
