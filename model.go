package whiterabbit

// Model go struct to add neo4j specific fields
type Model struct {
	ID     int64             `json:"id"`                  // neo4j node or relationship ID
	Labels map[string]string `json:"attribute,omitempty"` // any label not defined mapping struct
}

func (m Model) SetId(id int64) {
	m.ID = id
}
