package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func CreateNewDocument(id string, start string, sub string) redisearch.Document {
	fmt.Printf("id: %s, start: %s, sub %s \n", id, start, sub)
	doc := redisearch.NewDocument(id, 1.0)
	doc.Set("start", start)
	doc.Set("sub", sub)
	return doc
}

func readData(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var mapSubWithStart []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mapSubWithStart = append(mapSubWithStart, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return mapSubWithStart
}

func ExampleClient() {

	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClient("localhost:6379", "myIndex")

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).AddField(redisearch.NewTextField("sub"))

	// Drop an existing index. If the index does not exist an error is returned
	c.Drop()

	// Create the index with the given schema
	if err := c.CreateIndex(sc); err != nil {
		log.Fatal(err)
	}

	// Read data from test and index
	data := readData("data/test1.txt")
	for i := 0; i < len(data); i += 2 {
		doc := CreateNewDocument("video_part"+strconv.Itoa(i), data[i], data[i+1])
		if err := c.Index([]redisearch.Document{doc}...); err != nil {
			log.Fatal(err)
		}
	}

	// // Searching with limit and sorting
	// docs, total, err := c.Search(redisearch.NewQuery("hello world").
	// 	Limit(0, 2).
	// 	SetReturnFields("title"))

	// fmt.Println(docs[0].Id, docs[0].Properties["title"], total, err)
	// // Output: doc1 Hello world 1 <nil>
}

func main() {

	ExampleClient()

}
