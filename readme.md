# SqAuto (Squirrel Autogen)
Automatic query builder for golang [squirrel](https://github.com/Masterminds/squirrel)
## Why???
With squirrel and golang struct scanning i often have to write a lot of boilerplate query
```go
query,args,err := sq.Querybuilder.
                     Insert("tableName").
                     Columns("id",
                           "field1",
                           "field2",
                           "field3"
                           ....).
                     Values(
                            69,
                            "foo",
                            "bar",
                            42,
                     ).Tosql()
```
I was tired of this, especially when i have a struct with 10+ fields.
With `sqauto` you can now just call one function.
```go
obj := objStruct{ Id: 69, Field1: "foo", Field2: "bar", Field3: "42"}
query, args , err:= sqauto.Insert(
                        sq.Querybuilder,
                        sqauto.Table{ Name:"tableName", Object: obj}
                        ).Tosql()
)
```
These will generate the same query. 
Sqauto will read your struct for you and generate the appropriate query.

## TODO:
0. Better readme/documentaion
1. One to Many join
2. Many to Many join


