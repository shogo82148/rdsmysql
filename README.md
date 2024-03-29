[![GitHub Actions status](https://github.com/shogo82148/rdsmysql/workflows/Test/badge.svg)](https://github.com/shogo82148/rdsmysql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shogo82148/rdsmysql)](https://pkg.go.dev/github.com/shogo82148/rdsmysql)

# rdsmysql

The rdsmysql package is a SQL driver that allows IAM Database Authentication for Amazon RDS and Amazon Aurora.
It also supports connecting to the RDS proxy using IAM authentication.

- [IAM Database Authentication for MySQL and PostgreSQL - Amazon Relational Database Service](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html)
- [IAM Database Authentication - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html)
- [Managing connections with Amazon RDS Proxy - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy.html#rds-proxy-connecting-iam)

rdsmysql v1 works with [AWS SDK for Go v1](https://github.com/aws/aws-sdk-go):

```go
import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/shogo82148/rdsmysql"
)

func main() {
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
```

If you use [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2), use [rdsmysql v2](https://pkg.go.dev/github.com/shogo82148/rdsmysql/v2).

## Related Posts

- [How do I connect to my Amazon RDS MySQL DB instance or Aurora MySQL DB cluster using Amazon RDS Proxy?](https://aws.amazon.com/premiumsupport/knowledge-center/rds-aurora-mysql-connect-proxy/)
- [How do I allow users to authenticate to an Amazon RDS MySQL DB instance using their IAM credentials?](https://aws.amazon.com/premiumsupport/knowledge-center/users-connect-rds-iam/)
