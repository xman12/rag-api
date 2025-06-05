package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rag-api/pkg/db"
	backend "rag-api/pkg/ollama"
	"time"
)

var (
	host        = "http://localhost:11434"
	emdHost     = "http://localhost:11434"
	embModel    = "mxbai-embed-large"
	genModel    = "llama3"
	databaseURL = "postgres://user:password@localhost:5432/dbname?sslmode=disable"
)

type Query struct {
	Value string `uri:"query" binding:"required"`
}

type Data struct {
	Value string `form:"data" json:"data" binding:"required"`
}

func Run() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/search/:query", handleSearch)
	r.POST("/data", handleData)

	r.Run()
}

// handleData handles the POST request to insert data into the vector database.
func handleData(c *gin.Context) {
	var data Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data parameter",
		})
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	embeddingBackend := backend.NewOllamaBackend(emdHost, embModel, time.Duration(20*time.Second))
	vectorDB, err := db.NewPGVector(databaseURL)
	if err != nil {
		log.Fatalf("Error initializing vector database: %v", err)
	}
	defer vectorDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Embed the data using the specified embedding backend
	dataEmbedding, err := embeddingBackend.Embed(ctx, data.Value, headers)
	if err != nil {
		log.Fatalf("Error generating data embedding: %v", err)
	}
	log.Println("Vector embeddings generated")

	// Insert the embedded data into the vector database
	err = vectorDB.InsertDocument(ctx, data.Value, dataEmbedding)
	if err != nil {
		log.Fatalf("Error inserting data into vector database: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data inserted successfully",
	})
}

// handleSearch handles the GET request to search for relevant documents based on the query.
func handleSearch(c *gin.Context) {

	var query Query
	if err := c.ShouldBindUri(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameter",
		})
		return
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	embeddingBackend := backend.NewOllamaBackend(emdHost, embModel, time.Duration(20*time.Second))
	generationBackend := backend.NewOllamaBackend(host, genModel, time.Duration(20*time.Second))

	vectorDB, err := db.NewPGVector(databaseURL)
	if err != nil {
		log.Fatalf("Error initializing vector database: %v", err)
	}
	// Make sure to close the connection when done
	defer vectorDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close the connection when done
	defer vectorDB.Close()

	// Embed the query using the specified embedding backend
	queryEmbedding, err := embeddingBackend.Embed(ctx, query.Value, headers)
	if err != nil {
		log.Fatalf("Error generating query embedding: %v", err)
	}
	log.Println("Vector embeddings generated")

	// Retrieve relevant documents for the query embedding
	retrievedDocs, err := vectorDB.QueryRelevantDocuments(ctx, queryEmbedding, "ollama")
	if err != nil {
		log.Fatalf("Error retrieving relevant documents: %v", err)
	}

	// Log the retrieved documents to see if they include the inserted content
	for _, doc := range retrievedDocs {
		log.Printf("Retrieved Document: %v", doc)
	}

	// Augment the query with retrieved context
	augmentedQuery := db.CombineQueryWithContext(query.Value, retrievedDocs)

	prompt := backend.NewPrompt().
		AddMessage("system", "Вы — помощник ИИ. Используйте предоставленный контекст, чтобы ответить на вопрос пользователя как можно точнее. Не отвечает на вопросы не загруженного контекста. Пиши только по русски. После ответа сбрасывай контекст общения. Если результат не найден, то верни текст 'Данные не найдены'").
		AddMessage("user", augmentedQuery).
		SetParameters(backend.Parameters{
			MaxTokens:   150, // Supported by LLaMa
			Temperature: 0.7, // Supported by LLaMa
			TopP:        0.9, // Supported by LLaMa
		})

	// Generate response with the specified generation backend
	response, err := generationBackend.Generate(ctx, prompt)

	c.JSON(http.StatusOK, gin.H{
		"query":    query.Value,
		"response": response,
	})
}
