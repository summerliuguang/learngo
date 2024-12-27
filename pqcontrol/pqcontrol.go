package pqcontrol

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
)

const (
	Success = iota
	ConnectFailed
	QueryFailed
	ScanFailed
	AuthFailed
)

type CONNECT_DATA struct {
	user     string
	password string
	host     string
	port     int
	dbname   string
}

var db *sql.DB
var once sync.Once
var readConn CONNECT_DATA

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initReadConn() {
	readConn.user = os.Getenv("DB_READ_ROLE")
	readConn.password = os.Getenv("DB_READ_PASS")
	readConn.host = os.Getenv("DB_HOST")
	readConn.dbname = os.Getenv("DB_NAME")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Println("Convert port failed:", err)
		port = 5432
	}
	readConn.port = port
}

func initDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		initReadConn()
		connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", readConn.user, readConn.password, readConn.host, readConn.port, readConn.dbname)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("Open database failed:", err)
			return
		}
		err = db.Ping()
		if err != nil {
			log.Println("Ping database failed:", err)
			return
		}
	})
	return db, err
}

func GetUsers() ([]string, int) {
	db, err := initDB()
	if err != nil {
		log.Println("Init database failed:", err)
		return nil, ConnectFailed
	}

	var names []string
	rows, err := db.Query("SELECT user_name FROM ttkkai_user")
	if err != nil {
		log.Println("Query failed:", err)
		return nil, QueryFailed
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Println("Scan failed:", err)
			return nil, ScanFailed
		}
		names = append(names, name)
	}
	return names, Success
}

func GetUserById(userid string) (string, int) {
	db, err := initDB()
	if err != nil {
		log.Println("Init database failed:", err)
		return "", ConnectFailed
	}

	var name string
	err = db.QueryRow("SELECT user_name FROM ttkkai_user WHERE user_id = $1", userid).Scan(&name)
	if err != nil {
		log.Println("Query failed:", err)
		return "", QueryFailed
	}
	return name, Success
}

func AuthAccount(userid, password string) (string, int) {
	db, err := initDB()
	if err != nil {
		log.Println("Init database failed:", err)
		return "", ConnectFailed
	}
	var name string
	err = db.QueryRow("SELECT user_name FROM ttkkai_user WHERE user_id = $1 AND user_password = $2", userid, password).Scan(&name)
	if err != nil {
		log.Println("Query failed:", err)
		return "", AuthFailed
	}
	return name, Success
}
