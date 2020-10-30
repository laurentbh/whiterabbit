# whiterabbit

A "wanna be ORM" for [neo4j](https://neo4j.com/).

An easy way to map go struct with neo4j entities


## Connection to neo4j
first step is to open the database
```go
db, err := whiterabbit.Open(config)
```
- where config is any struct that satisfy the `Config` interface
- close the database `db.Close()`

## Getting a connection
a connection, ie a `session` in neo4j world can be obtain:
```go
con, err := db.GetConnection()
```
`sessions` are meant to be short-live in neo4j world, so they shall be closed as soon as possible `con.Close()`


## creating a node
- nodes create have the name of the struct as label
- every exported field of the struct is an attribute of the node

For example:
```go
type User struct {
    Name     string
    Age      int
    password string
}
user := User{Name: "first", Age: 10}
con.CreateNode(user)
```
will create a `Node` labelled `User` with 2 attributes, `Name` and `Age`.
`password` is ignored

## fetching nodes
```go
type User struct {
	whiterabbit.Model
    Name     string
    Age      int
    password string
}
con.FindNodes(User{})
```

If the struct contains `whiterabbit.Model`, all node's attributes that are not matching an exported field of the struct, will be copy in `whiterabbit.Model` [see](###white-rabbit-model)

## converting nodes
`neo4j.Node` can be converted to `struct`:
```go
var node neo4j.Node
type User struct {
	whiterabbit.Model
    Name     string
    Age      int
    password string
}
ret, _ := ConvertNode(node, []interface{}{User{}}))
u = ret.(User)
```
- `ConvertNode` takes an array of `struct` as second parameter representing all potential internal target of the node.
- if present `whiterabbit.Model` will be initialized, [see](###white-rabbit-model)

### white rabbit model
This struct can be added to your struct to retrieve nodes' attributes not declare in your struct.
```go
type Model struct {
	ID     int64             // neo4j node or relationship ID
	Labels map[string]string // any label not defined mapping struct
}
```

## transactions
using `Connection.InTransaction(f)`
where `f` is 
```go
func(con *Connection) ([]neo4j.Result, error)
```
transaction is rollbacked if `f` returns an error




