[![GitHub Actions status](https://github.com/shogo82148/rdsmysql/workflows/Go/badge.svg)](https://github.com/shogo82148/rdsmysql)
[![GoDoc](https://godoc.org/github.com/shogo82148/rdsmysql?status.svg)](https://godoc.org/github.com/shogo82148/rdsmysql)

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
