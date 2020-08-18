package driver

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
)
// Read only environments variable
const (
	CassandraUrl      = "CASSANDRA_URL"
	CassandraKeyspace = "CASSANDRA_KEYSPACE"
)

// InitCluster initializes cluster specified in appEnv.env file
// Creates Email table and returns cassandra session
func InitCluster() *gocql.Session {
	cassandra := CheckEnvVar(CassandraUrl)
	keyspace := CheckEnvVar(CassandraKeyspace)
	cluster := CreateCluster(cassandra, keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
	CreateEmailTable(keyspace, session)
	return session
}

// CheckEnvVar checks whether environment variable key exists
func CheckEnvVar(key string) string {
	cassandra, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal("Environment variable doesn't exists :" + key)
	}
	return cassandra
}

// CreateCluster creates cluster with provided host and keyspace
// Returns cluster config
func CreateCluster(host string, keyspace string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(host)
	CreateKeyspace(keyspace, cluster)
	cluster.Keyspace = keyspace
	return cluster
}

// CreateKeyspace creates keyspace if it doesnt exists in specified cluster
func CreateKeyspace(keyspace string, cluster *gocql.ClusterConfig) {
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
	defer session.Close()

	createKeyspaceQuery := "CREATE KEYSPACE IF NOT EXISTS " + keyspace + " WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};"

	if err := session.Query(createKeyspaceQuery).Exec(); err != nil {
		log.Fatal("Error creting keyspace: ", err)
	}

	log.Println("Keyspace created: ", keyspace)
}

// CreateEmailTable creates table and index on email address if does not exist in specified keyspace
func CreateEmailTable(keyspace string, s *gocql.Session) {
	tableName := keyspace + ".Email"
	createTableQuery := "CREATE TABLE IF NOT EXISTS " + tableName + "\n" +
		"(\n    email        text," +
		"\n    title        text," +
		"\n    content      text," +
		"\n    magic_number int," +
		"\n    PRIMARY KEY (email, magic_number, content)\n);"
	createIndex := "CREATE INDEX  IF NOT EXISTS mNumber_index ON " + tableName + "(magic_number);"
	if err := s.Query(createTableQuery).Exec(); err != nil {
		log.Println("Table was not created: ", err)
	}
	s.Query(createIndex).Exec()

	log.Println(fmt.Sprintf("Table %s was created", tableName))
}
