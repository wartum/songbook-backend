package db_access

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

var connection_url string = "postgres://postgres:Songs@localhost:5432/songs"
var db *sql.DB

type User struct {
	Id       string
	Name     string
	Password string
	Salt     string
}

func CheckHealth() {
	connection := create_connection()
	defer connection.Close(context.Background())

	log.Println("Connection succesfull")
}

func create_connection() *pgx.Conn {
	connection, err := pgx.Connect(context.Background(), connection_url)
	if err != nil {
		log.Fatalf("Could not connect to database. %v\n", err)
	}
	return connection
}

func GetUser(name string) (User, error) {
	connection := create_connection()
	defer connection.Close(context.Background())

	query := fmt.Sprintf("SELECT id, name, password, salt FROM users WHERE name='%v'", name)
	row, err := connection.Query( context.Background(), query)
	if err != nil {
		return User{}, fmt.Errorf("SQL query execution failed. %v\n", err)
	}

	var user User
	row.Next()
	row.Scan(&user.Id, &user.Name, &user.Password, &user.Salt)

	if user.Name == "" {
		return User{}, fmt.Errorf("User doesn't exist")
	}

	return user, nil
}
