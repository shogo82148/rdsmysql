[![GitHub Actions status](https://github.com/shogo82148/rdsmysql/workflows/Test/badge.svg)](https://github.com/shogo82148/rdsmysql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shogo82148/rdsmysql)](https://pkg.go.dev/github.com/shogo82148/rdsmysql)

# rdsmysql

The rdsmysql package is a SQL driver that allows IAM Database Authentication for Amazon RDS and Amazon Aurora.
It also supports connecting to the RDS proxy using IAM authentication.

- [IAM Database Authentication for MySQL and PostgreSQL - Amazon Relational Database Service](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html)
- [IAM Database Authentication - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html)
- [Managing connections with Amazon RDS Proxy - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy.html#rds-proxy-connecting-iam)

``` go
// configure AWS session
awsConfig := aws.NewConfig().WithRegion("ap-northeast-1")
awsSession := session.Must(session.NewSession(awsConfig))

// configure the driver
driver := &rdsmysql.Driver{
    Session: awsSession,
}
sql.Register("rdsmysql", driver)

db, err := sql.Open("rdsmysql", "user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
if err != nil {
    t.Fatal(err)
}
defer db.Close()
```

## Related Posts

- [How do I connect to my Amazon RDS MySQL DB instance or Aurora MySQL DB cluster using Amazon RDS Proxy?](https://aws.amazon.com/premiumsupport/knowledge-center/rds-aurora-mysql-connect-proxy/)
- [How do I allow users to authenticate to an Amazon RDS MySQL DB instance using their IAM credentials?](https://aws.amazon.com/premiumsupport/knowledge-center/users-connect-rds-iam/)
