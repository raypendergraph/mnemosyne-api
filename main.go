package main

import (
	"fmt"
	"mnemosyne-api/adapters/inbound/http"
	repo "mnemosyne-api/adapters/outbound/neo4j"
	sys "mnemosyne-api/system"
)

func main() {
	config := sys.Configuration{
		HTTP: sys.HTTPConfiguration{Port: 8080},
	}

	errorCatalog := sys.NewCatalog(2)
	logger := configureLogging()

	var neo4jAdapter repo.Neo4JAdapter
	{
		var err sys.Error
		if neo4jAdapter, err = repo.NewRepoBase(config.Neo4J, errorCatalog); err != nil {
			logger.LogError(err)
			panic(err)
		}
	}

	port := fmt.Sprintf(":%d", config.HTTP.Port)
	if err := http.NewAdapter(config.HTTP, errorCatalog, neo4jAdapter, logger).Run(port); err != nil {
		panic(err)
	}
}

func configureLogging() sys.LoggerImpl {
	return sys.NewLogger()
}
