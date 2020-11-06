package integration

// Cfg for testing
type Cfg struct {
	host     string
	port     int
	user     string
	password string
}

// GetHost ...
func (df Cfg) GetHost() string {
	return "localhost"
}

// GetPort ...
func (df Cfg) GetPort() int {
	return 7688
}

// GetUser ...
func (df Cfg) GetUser() string {
	return "neo4j"
}

// GetPassword ...
func (df Cfg) GetPassword() string {
	return "root"
}

// GetEncrypted ...
func (df Cfg) GetEncrypted() bool {
	return false
}
