package driver

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
)

const (
	CassandraUrl      = "CASSANDRA_URL"
	CassandraKeyspace = "CASSANDRA_KEYSPACE"
)

func InitCluster() *gocql.Session {
	cassandra := checkEnvVar(CassandraUrl)
	keyspace := checkEnvVar(CassandraKeyspace)
	cluster := createCluster(cassandra, keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
	createEmailTable(keyspace, session)
	return session
}

func checkEnvVar(key string) string {
	cassandra, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal("Environment variable doesn't exists :" + key)
	}
	return cassandra
}

func createCluster(host string, keyspace string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(host)
	createKeyspace(keyspace, cluster)
	cluster.Keyspace = keyspace
	return cluster
}

func createKeyspace(keyspace string, cluster *gocql.ClusterConfig) {
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

func createEmailTable(keyspace string, s *gocql.Session) {
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

func deleteEmailTable(keyspace string, s *gocql.Session) {
	tableName := keyspace + ".Email"
	deleteTableQuery := "DROP TABLE IF EXISTS " + keyspace + ".Email;"

	s.Query(deleteTableQuery).Exec()

	log.Println(fmt.Sprintf("Table %s was droped\n", tableName))
}
