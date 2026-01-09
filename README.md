# Recipe Search Service
This service ingests recipe JSON datasets and provides context-aware search.

# Flow diagram

```
            +--------------------+
            | Recipe JSON Dataset |
            +----------+---------+
                       |
                       v
               +-------+--------+
               | Ingestor       |
               | - normalize    |
               | - embed        |
               | - upsert       |
               +-------+--------+
                       |
                       v
               +-------+--------+
               | Qdrant         |
               | vectors+payload|
               +-------+--------+
                       ^
                       |
            +----------+---------+
            | Search Service     |
            | - parse context    |
            | - embed query      |
            | - filter+search    |
            | - rerank/explain   |
            +----------+---------+
                       ^
                       |
               +-------+--------+
               | HTTP /search   |
               +----------------+
```

## Requirements

- Go
- [Qdrant](https://qdrant.tech/documentation/overview):

```sh
docker pull qdrant/qdrant

docker run -p 6333:6333 qdrant/qdrant
```

## Run the server

```sh
go run ./cmd/server
```
