[![Build Status](https://travis-ci.com/shogo82148/rdsmysql.svg?branch=master)](https://travis-ci.com/shogo82148/rdsmysql)
[![GoDoc](https://godoc.org/github.com/shogo82148/rdsmysql?status.svg)](https://godoc.org/github.com/shogo82148/rdsmysql)

# rdsmysql

The rdsmysql package is a SQL driver for Amazon RDS.
It allows [Authentication and Access Control for Amazon RDS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAM.html) using IAM.

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