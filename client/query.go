package RSSClient

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
)

// Read query from STDIN
// reader := bufio.NewReader(os.Stdin)
// fmt.Print("Enter text search: ")
// queryText, _ := reader.ReadString('\n')

//TODO: Optimize query
func preprocessQuery(queryText string) string {
	space := regexp.MustCompile(`\s+`)
	queryText = space.ReplaceAllString(queryText, " ")
	queryText = strings.TrimSpace(queryText)
	queryText = strings.Replace(queryText, " ", " | ", -1)
	// fmt.Println(queryText)
	return queryText
}

func query(c *redisearch.Client, queryText string) []string {

	query := redisearch.NewQuery(queryText)

	// Searching for text
	docs, total, err := c.Search(query)

	if err != nil {
		log.Fatal(err)
	}

	//TODO: split to a function to parse the results
	results := make([]string, 0, len(docs)+1)

	results = append(results, "Total results: "+strconv.Itoa(total))
	for _, doc := range docs {
		results = append(results, "Start at:"+doc.Properties["start"].(string)+"\tSub:"+doc.Properties["sub"].(string))
		// fmt.Println(results[i])
	}

	return results
}

func Search(c *redisearch.Client, phrase string) []string {
	queryText := preprocessQuery(phrase)
	return query(c, queryText)
}
