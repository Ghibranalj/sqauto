# SqAuto (Squirrel Autogen)
Automatic query builder for golang [squirrel](https://github.com/Masterminds/squirrel)
## Why???
With squirrel and golang struct scanning i often have to write a lot of boilerplate query
```go
query,args,err := sq.Querybuilder.
                     Insert("table").
                     Columns("id",
                           "field1",
                           "field2",
                           "field3"
                           ....).
                     Values(
                            id,
                            val1,
                            val2,
                            val3
                     ).Tosql()
```
I was tired of this, especially when i have a struct with 10+ fields.
With `sqauto` you can now just call one function.
```go
query, args , err:= sqauto.Insert(
                        sq.Querybuilder,
                        sqauto.Table{Name="tableName", 
                        Object:structOBj{
                        id:id, val1: var1, val2: var2,
                        }}).Tosql()
)
```
These will generate the same query.


