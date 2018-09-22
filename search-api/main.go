package main

import (
	//"encoding/json"
	//"fmt"
	"log"
	//"net/http"
	//"strconv"
	"time"

	//"github.com/gin-gonic/gin"
	//"github.com/olivere/elastic"
	//"github.com/teris-io/shortid"
	"net/http"
)

const (
	elasticIndexName = "documents"
	elasticTypeName  = "document"
)

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

var (
	elasticClient *elastic.Client
)

type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// A helper function for responding with an error.
func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

// read documents from request body into an array.
func createDocumentsEndpoint(c *gin.Context) {
	var docs []DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	bulk := elasticClient.
		Bulk().
		Index(elasticIndexName).
		Type(elasticTypeName)

	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(),
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to created documents")
		return
	}
	c.Status(http.StatusOK)
}

func main() {
	var err error
	for {
		elasticClient, err := err.NewClient(
			elastic.SetUrl("http://elasticsearch:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	// ...
	r := gin.Default()
	r.POST("/documents", createDocumentsFunction)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
