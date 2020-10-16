package whiterabbit

// Config ...
type Config interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetEncrypted() bool
}

// DefaultConfig for testing
type DefaultConfig struct {
	host     string
	port     int
	user     string
	password string
}

// GetHost ...
func (df DefaultConfig) GetHost() string {
	return "localhost"
}

// GetPort ...
func (df DefaultConfig) GetPort() int {
	return 7687
}

// GetUser ...
func (df DefaultConfig) GetUser() string {
	return "neo4j"
}

// GetPassword ...
func (df DefaultConfig) GetPassword() string {
	return "root"
}

// GetEncrypted ...
func (df DefaultConfig) GetEncrypted() bool {
	return false
}
