package system

type Configuration struct {
	HTTP  HTTPConfiguration  `yaml:"http"`
	Neo4J Neo4JConfiguration `yaml:"neo4j"`
}

type HTTPConfiguration struct {
	Port uint `yaml:"port"`
}

type Neo4JConfiguration struct {
	URL string `yaml:"url"`
}
