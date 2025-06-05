# RAG-API

RAG-API is a Go-based application that provides a RESTful API for managing and querying vector embeddings using PostgreSQL with the `pgvector` extension. It supports embedding generation, vector storage, and similarity search operations.

## Features

- **Vector Embedding Storage**: Store high-dimensional vector embeddings in a PostgreSQL database.
- **Similarity Search**: Query the most relevant documents based on cosine similarity.
- **RESTful API**: Expose endpoints for embedding data and searching for relevant documents.
- **Backend Integration**: Supports embedding and generation backends for processing data.

## Requirements

- **Go**: Version 1.20 or higher.
- **PostgreSQL**: Version 14 or higher with the `pgvector` extension installed.
- **Dependencies**: Managed via `go.mod`.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/rag-api.git
   cd rag-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the PostgreSQL database:
   - Install the `pgvector` extension in your PostgreSQL instance.
   - Create the required tables:
     ```sql
     CREATE TABLE ollama_embeddings (
         doc_id TEXT PRIMARY KEY,
         embedding VECTOR(1024),
         metadata JSONB
     );
     ```

4. Update the database connection string in `internal/app/server.go`:
   ```go
   var databaseURL = "postgres://user:password@localhost:5432/dbname?sslmode=disable"
   ```

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```

2. API Endpoints:
   - **Ping**: Test the server status.
     ```bash
     GET /ping
     Response: {"message": "pong"}
     ```

   - **Insert Data**: Insert a document and its embedding into the database.
     ```bash
     POST /data
     Body: {"data": "Your content here"}
     Response: {"message": "Data inserted successfully"}
     ```

   - **Search**: Search for relevant documents based on a query.
     ```bash
     GET /search/:query
     Response: {"query": "Your query", "response": "Generated response"}
     ```

## Project Structure

- `internal/app/server.go`: Main application logic and API endpoints.
- `pkg/db/pgvector.go`: PostgreSQL vector database implementation.
- `pkg/db/vectordb.go`: Interface and utility functions for vector database operations.
- `pkg/ollama`: Backend integration for embedding and generation.

## Configuration

- **Embedding Backend**:
  - Host: `http://localhost:11434`
  - Embedding Model: `mxbai-embed-large`
  - Generation Model: `llama3`

- **Environment Variables**:
  - Update the database URL and backend hosts in the code or use environment variables for better security.

## Dependencies

- [Gin](https://github.com/gin-gonic/gin): HTTP web framework.
- [pgx](https://github.com/jackc/pgx): PostgreSQL driver for Go.
- [pgvector-go](https://github.com/pgvector/pgvector-go): Go client for `pgvector`.

## License

This project is licensed under the Apache License 2.0. See the `LICENSE` file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.