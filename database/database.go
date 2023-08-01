package database

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	defaultMaxIdleConnections = 5
	defaultMaxOpenConnections = 50
)

var db *sql.DB

type DbConnection struct {
	Conn *sql.DB
}

type DbConfig struct {
	Host            string
	Port            int
	User            string
	Pass            string
	DbName          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxLifetimeMins int
}

func InitDatabase(config DbConfig) (DbConnection, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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
		if config.MaxIdleConn == 0 {
			log.Debugf("InitDatabase:: setting SetMaxIdleConns to %d", defaultMaxIdleConnections)
			connection.Conn.SetMaxIdleConns(defaultMaxIdleConnections)
		} else {
			log.Debugf("InitDatabase:: setting SetMaxIdleConns to %d", config.MaxIdleConn)
			connection.Conn.SetMaxIdleConns(config.MaxIdleConn)
		}

		if config.MaxOpenConn == 0 {
			log.Debugf("InitDatabase:: setting SetMaxOpenConns to %d", defaultMaxOpenConnections)
			connection.Conn.SetMaxOpenConns(defaultMaxOpenConnections)
		} else {
			log.Debugf("InitDatabase:: setting SetMaxOpenConns to %d", config.MaxOpenConn)
			connection.Conn.SetMaxOpenConns(config.MaxOpenConn)
		}

		if config.MaxLifetimeMins != 0 {
			log.Debugf("InitDatabase:: setting ConnMaxLifetime to %d mins", config.MaxLifetimeMins)
			connection.Conn.SetConnMaxLifetime(time.Minute * time.Duration(config.MaxLifetimeMins))
		}
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
