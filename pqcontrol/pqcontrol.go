package pqcontrol

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	Success = iota
	ConnectFailed
	QueryFailed
	ScanFailed
	AuthFailed
	InsertFailed
	UpdateFailed
	DeleteFailed
	UnknownError
	UserAlreadyExists
	CryptFailed
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

func AuthAccount(username string, password string) (string, int) {
	db, err := initDB()
	if err != nil {
		log.Println("Init database failed:", err)
		return "", ConnectFailed
	}
	var hashedPassword string
	err = db.QueryRow("SELECT user_password FROM ttkkai_user WHERE user_name = $1", username).Scan(&hashedPassword)
	if err != nil {
		log.Println("Query failed:", err)
		return "", AuthFailed
	}
	if !checkPasswordHash(password, hashedPassword) {
		return "", AuthFailed
	}
	return username, Success
}

// create a new userid with snoyflke
func createUserID() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatal(err)
	}
	id := node.Generate()
	return id.Int64()
}

func CreateAccount(Username, password string) (int64, int) {
	db, err := initDB()
	if err != nil {
		log.Println("Init database failed:", err)
		return 0, ConnectFailed
	}

	var existingUsername string
	err = db.QueryRow("SELECT user_name FROM ttkkai_user WHERE user_name = $1", Username).Scan(&existingUsername)
	if err == nil {
		log.Println("Username already exists:", Username)
		return 0, UserAlreadyExists
	}

	userid := createUserID()
	hashedPassword, err := cryptPassword(password)
	if err != nil {
		log.Println("Crypt password failed:", err)
		return 0, CryptFailed
	}
	err = db.QueryRow("INSERT INTO ttkkai_user (user_id, user_name, user_password) VALUES ($1, $2, $3) ON CONFLICT (user_name) DO NOTHING RETURNING user_id", userid, Username, hashedPassword).Scan(&userid)
	if err != nil {
		log.Println("Insert user failed:", err)
		return 0, InsertFailed
	}
	return userid, Success
}

func cryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
