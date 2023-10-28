[![GitHub Actions status](https://github.com/shogo82148/rdsmysql/workflows/Test/badge.svg)](https://github.com/shogo82148/rdsmysql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shogo82148/rdsmysql/v2)](https://pkg.go.dev/github.com/shogo82148/rdsmysql/v2)

# rdsmysql

The rdsmysql package is a SQL driver that allows IAM Database Authentication for Amazon RDS and Amazon Aurora.
It also supports connecting to the RDS proxy using IAM authentication.

- [IAM Database Authentication for MySQL and PostgreSQL - Amazon Relational Database Service](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html)
- [IAM Database Authentication - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html)
- [Managing connections with Amazon RDS Proxy - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy.html#rds-proxy-connecting-iam)

rdsmysql v2 works with [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2):

```go
import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/rdsmysql/v2"
)

func main() {
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
```

If you use [AWS SDK for Go v1](https://github.com/aws/aws-sdk-go), use [rdsmysql v1](https://pkg.go.dev/github.com/shogo82148/rdsmysql).

## Related Posts

- [How do I connect to my Amazon RDS MySQL DB instance or Aurora MySQL DB cluster using Amazon RDS Proxy?](https://aws.amazon.com/premiumsupport/knowledge-center/rds-aurora-mysql-connect-proxy/)
- [How do I allow users to authenticate to an Amazon RDS MySQL DB instance using their IAM credentials?](https://aws.amazon.com/premiumsupport/knowledge-center/users-connect-rds-iam/)
