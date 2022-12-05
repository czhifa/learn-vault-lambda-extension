package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/vault/api"

	_ "github.com/go-sql-driver/mysql"
)

const functionName = "demo-function"

// Payload captures the basic payload we're sending for demonstration
// Ex: {"payload": "hello"}
type Payload struct {
	Message string `json:"payload"`
}

// String prints the payload recieved
func (m Payload) String() string {
	return m.Message
}

// HandleRequest reads credentials from /tmp and uses them to query the database
// for users. The database is determined by the DATABASE_URL environment
// variable, and the username and password are retrieved from the secret.
func HandleRequest(ctx context.Context, payload Payload) error {
	logger := log.New(os.Stderr, fmt.Sprintf("[%s] ", functionName), 0)
	logger.Println("Received:", payload.String())
	secretRaw, err := ioutil.ReadFile("/tmp/vault_secret.json")
	logger.Println("Reading file /tmp/vault_secret.json")
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// read token
	// tokenRaw, err := ioutil.ReadFile("/tmp/vault/token")
	// if err != nil {
	// 	return fmt.Errorf("error reading file: %w", err)
	// }

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("no DATABASE_URL, exiting")
	}

	// First decode the JSON into a map[string]interface{}
	var secret api.Secret
	b := bytes.NewBuffer(secretRaw)
	dec := json.NewDecoder(b)
	// While decoding JSON values, interpret the integer values as `json.Number`s
	// instead of `float64`.
	dec.UseNumber()

	if err := dec.Decode(&secret); err != nil {
		return err
	}

	// read users from database
	logger.Println("username: ")
	logger.Println("    ", secret.Data["username"])
	logger.Println("password: ")
	logger.Println("    ", secret.Data["password"])
	logger.Println("dbURL: ")
	logger.Println("    ", dbURL)

	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/lambdadb", secret.Data["username"], secret.Data["password"], dbURL)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}

	var titles []string
	rows, err := db.QueryContext(ctx, "SELECT title FROM article")
	if err != nil {
		return err
	}
	// defer rows.Close()
	for rows.Next() {
		var title string
		if err = rows.Scan(&title); err != nil {
			return err
		}
		titles = append(titles, title)
	}
	logger.Println("titles: ")
	for i := range titles {
		logger.Println("    ", titles[i])
	}

	return nil
}

// func callQuery() error {
// 	logger := log.New(os.Stderr, fmt.Sprintf("[%s] ", functionName), 0)

// 	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql", "root", "rootpassword", "127.0.0.1")
// 	db, err := sql.Open("mysql", connStr)
// 	if err != nil {
// 		return err
// 	}

// 	var users []string
// 	//rows, err := db.QueryContext(ctx, "SELECT user FROM mysql.user")
// 	rows, err := db.Query("select user from mysql.user")
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var user string
// 		if err = rows.Scan(&user); err != nil {
// 			return err
// 		}
// 		users = append(users, user)
// 	}
// 	logger.Println("users: ")
// 	for i := range users {
// 		logger.Println("    ", users[i])
// 	}

// 	return nil
// }

func main() {
	//callQuery()
	lambda.Start(HandleRequest)
}
