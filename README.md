[![GitHub Actions status](https://github.com/shogo82148/rdsmysql/workflows/Test/badge.svg)](https://github.com/shogo82148/rdsmysql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shogo82148/rdsmysql)](https://pkg.go.dev/github.com/shogo82148/rdsmysql)

# rdsmysql

The rdsmysql package is a SQL driver that allows IAM Database Authentication for Amazon RDS and Amazon Aurora.

- [IAM Database Authentication for MySQL and PostgreSQL - Amazon Relational Database Service](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html)
- [IAM Database Authentication - Amazon Aurora](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html)

``` go
c := aws.NewConfig().WithRegion("ap-northeast-1")
s := session.Must(session.NewSession(c))
d := &Driver{
    Session: s,
}
sql.Register("rdsmysql", d)

db, err := sql.Open("rdsmysql", "user:@tcp(db-foobar.ap-northeast-1.rds.amazonaws.com:3306)/")
if err != nil {
    t.Fatal(err)
}
defer db.Close()
```
