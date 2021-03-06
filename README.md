# whiterabbit

Small library to map go struct with  [neo4j](https://neo4j.com/) entities.


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
con.FindAllNodes(User{})
con.FindNodesClause(User{}, map[string]interface{}{"Name": "user", "Age": 10, }, whiterabbit.StartsWith)
```

If the struct contains `whiterabbit.Model`, all node's attributes that are not matching an exported field of the struct, will be copy in `whiterabbit.Model` [see](#white-rabbit-model)

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
ret, _ := ConvertNode(node, User{}))
u = ret.(User)
```
- if present `whiterabbit.Model` will be initialized, [see](#white-rabbit-model)
- `ConvertNode` has a variadic third argument representing all struct candidate for the node

### white rabbit model
adding `whiterabbit.Model` to struct allows:
- get `neo4j` node IDs
- retrieve nodes' attributes not declared in struct
```go
type Model struct {
	ID     int64             // neo4j node or relationship ID
	Labels map[string]string // any label not defined mapping struct
}
```
## relations
```go
relations, err := con.MatchRelation("is_a", Ingredient{}, Category{})
```
## transactions
using `Connection.InTransaction(f)`
where `f` is 
```go
func(con *Connection) ([]neo4j.Result, error)
```
transaction is rollbacked if `f` returns an error




