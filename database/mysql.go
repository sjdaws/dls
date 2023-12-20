package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Database struct {
    connection *sql.DB
}

type query struct {
    rows *sql.Rows
}

type statement struct {
    prepared *sql.Stmt
}

// Connect connects to the database
func Connect(dbname string, host string, password string, username string) (*Database, error) {
    if dbname == "" {
        dbname = "dls"
    }

    connection, err := createConnection(dbname, host, password, username)
    if err != nil {
        return nil, err
    }

    db := &Database{
        connection: connection,
    }

    err = db.connection.Ping()
    if err != nil {
        connection, err = createConnection(dbname, host, password, username)
        if err != nil {
            return nil, err
        }

        db = &Database{
            connection: connection,
        }
        err = db.create(dbname)

        if err != nil {
            return nil, fmt.Errorf("unable to create database %s: %v", dbname, err)
        }
    }

    db.connection.SetConnMaxLifetime(time.Minute * 3)
    db.connection.SetMaxIdleConns(10)
    db.connection.SetMaxOpenConns(10)

    err = db.migrate(dbname)
    if err != nil {
        return nil, fmt.Errorf("unable to complete database migration: %v", err)
    }

    return db, nil
}

func (d *Database) Close() {
    err := d.connection.Close()
    if err != nil {
        log.Printf("ERROR: unable to close database connection: %v", err)
    }
}

func (q *query) close() {
    err := q.rows.Close()
    if err != nil {
        log.Printf("ERROR: unable to close query: %v", err)
    }
}

func (s *statement) close() {
    err := s.prepared.Close()
    if err != nil {
        log.Printf("ERROR: unable to close statement: %v", err)
    }
}

func buildDsn(dbname string, host string, password string, username string) string {
    if password != "" {
        password = fmt.Sprintf(":%s", password)
    }

    return fmt.Sprintf("%s%s@tcp(%s)/%s?multiStatements=true&parseTime=true", username, password, host, dbname)
}

func createConnection(dbname string, host string, password string, username string) (*sql.DB, error) {
    connection, err := sql.Open("mysql", buildDsn(dbname, host, password, username))
    if err != nil {
        return nil, fmt.Errorf("unable to create database connection: %v", err)
    }

    return connection, nil
}
