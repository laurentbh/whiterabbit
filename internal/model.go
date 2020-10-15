package internal

// Model go struct to add neo4j specific fields
type Model struct {
	ID     int64             // neo4j node or relationship ID
	Labels map[string]string // any label not defined mapping struct
}
