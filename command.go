package whiterabbit

import (
	"strings"

	"github.com/laurentbh/whiterabbit/internal"
)

// CreateNode ...
func (db *DB) CreateNode(value interface{}) error {

	session, _ := db.GetSession()

	defer session.Close()

	mapping, _ := internal.GetMapping(value)
	cyp := createNodeCypher(mapping)

	result, err := session.Run(
		cyp,
		mapping.Values,
	)
	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
func createNodeCypher(mapping internal.Mapping) (ret string) {

	var builder strings.Builder
	// REATE (n:Ingredient{name: "🥩 beef"});
	builder.WriteString("CREATE (n:")
	builder.WriteString(mapping.Label)

	if len(mapping.Attributes) == 0 {
		builder.WriteString(")")
		return builder.String()
	}
	builder.WriteString("{")

	if len(mapping.Attributes) > 0 {
		sep := false
		for k := range mapping.Attributes {
			if sep {
				builder.WriteString(", ")
			}
			builder.WriteString(k)
			builder.WriteString(": $")
			builder.WriteString(k)

			sep = true
		}
	}
	builder.WriteString("})")
	ret = builder.String()
	return

}
