package integration

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/laurentbh/whiterabbit"
)

// LoadFixure ...
func LoadFixure(files []string) {
	neo, err := whiterabbit.Open(Cfg{})
	if err != nil {
		panic(err)
	}
	defer neo.Close()

	con, _ := neo.GetConnection()
	defer con.Close()

	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		cypher := string(content)

		cmds := strings.Split(cypher, ";")

		for _, c := range cmds {
			if len(c) != 0 {
				res, err := con.GetSession().Run(c,
					map[string]interface{}{})
				if err != nil {
					panic(err)
				}
				if res.Err() != nil {
					panic(res.Err())

				}
			}
		}
	}
}
