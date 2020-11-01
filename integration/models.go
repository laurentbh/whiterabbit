package integration

import "github.com/laurentbh/whiterabbit"

// Category ...
type Category struct {
	whiterabbit.Model
	Name string
}

// Ingredient ...
type Ingredient struct {
	whiterabbit.Model
	Name string
}

// User structure for tests
type User struct {
	whiterabbit.Model
	Name string
	Age  int
}
