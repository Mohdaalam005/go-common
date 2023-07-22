package database

import (
	"database/sql"
	"github.com/pkg/errors"

	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

type DbConnection struct {
	Conn *sql.DB
}

type DbConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

func InitDatabase(config DbConfig) (DbConnection, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disalbe",
		config.Host, config.Port, config.User, config.Pass, config.DbName)

	connection := DbConnection{}
	if db != nil {
		connection.Conn = db
	} else {
		pdb, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			return DbConnection{}, errors.Wrap(err, "while connecting to database")
		}
		connection.Conn = pdb
	}
	return connection, nil
}

func PingDB(conn *DbConnection) error {
	err := conn.Conn.Ping()
	if err != nil {
		log.Errorf("Error pinging the database %v", err)
	}
	return err
}
