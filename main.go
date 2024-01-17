package main

import (
	"database/sql"
	"fmt"

	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

type Book struct {
	Id      int
	Name    string
	Count   int
	User_id int
}

type User struct {
	Id   int
	Name string
	Age  int
}

type ReturningRow struct {
	Id        int
	UserName  string
	UserAge   int
	BookName  string
	BookCount int
}

func main() {
	connection := "host=rosie.db.elephantsql.com user=zvtgqvoa dbname=zvtgqvoa password=TKd9_asbgfbjZvM8i8B9ccGf46Y5nqIs sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name VARCHAR, age INT)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS books(id SERIAL PRIMARY KEY, name VARCHAR, count INT, user_id INT, FOREIGN KEY (user_id) REFERENCES users(id))")
	if err != nil {
		panic(err)
	}

	createUser(db)

}

func createUser(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	var userName string
	var userAge int
	fmt.Print("Enter username: ")
	fmt.Scan(&userName)
	fmt.Print("Enter user age: ")
	fmt.Scan(&userAge)
	var id int
	err = tx.QueryRow("INSERT INTO users(name, age) VALUES($1, $2) RETURNING id", userName, userAge).Scan(&id)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	var bookName string
	var bookCount int
	fmt.Print("Enter book name to rent: ")
	fmt.Scan(&bookName)
	fmt.Print("How many books want to rent: ")
	fmt.Scan(&bookCount)
	_, err = tx.Exec("INSERT INTO books(name, count, user_id) VALUES ($1, $2, $3)", bookName, bookCount, id)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		panic(err)
	}
	fmt.Print("Succesfully created user with book!!!")
}

func updateDetails(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	var menu int
	var id int
	fmt.Print("Which data do you want to update:\n1 - User`s name\n2 - User`s age\n3 - Book`s name\n4 - Book`s count")
	fmt.Scan(&menu)
	switch menu {
	case 1:
		fmt.Print("Enter user id to update: ")
		fmt.Scan(&id)
		var newUserName string
		fmt.Print("Enter new user name: ")
		fmt.Scan(&newUserName)
		_, err = tx.Exec("UPDATE users SET name = $1 WHERE id = $2", newUserName, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

	case 2:
		fmt.Print("Enter user id to update: ")
		fmt.Scan(&id)
		var newUserAge int
		fmt.Print("Enter user`s new age: ")
		fmt.Scan(&newUserAge)
		_, err = tx.Exec("UPDATE users SET age = $1 WHERE id = $2", newUserAge, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

	case 3:
		fmt.Print("Enter book id to update: ")
		fmt.Scan(&id)
		var newBookName string
		fmt.Print("Enter new book`s name: ")
		fmt.Scan(&newBookName)
		_, err = tx.Exec("UPDATE books SET name = $1 WHERE id = $2", newBookName, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

	case 4:
		fmt.Print("Enter book id to update: ")
		fmt.Scan(&id)
		var newBookCount int
		fmt.Print("Enter new book count: ")
		fmt.Scan(&newBookCount)
		_, err = tx.Exec("UPDATE books SET count = $1 WHERE id = $2", newBookCount, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func deleteData(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	var deleter int
	var id int
	fmt.Print("How do you want to delete rows:\n1 - With user id\n2 - With book id")
	fmt.Scan(&deleter)
	switch deleter {
	case 1:
		fmt.Print("Enter user id to delete: ")
		fmt.Scan(&id)
		_, err = tx.Exec("DELETE from books WHERE user_id = $1", id)
		if err != nil {
			tx.Rollback()
			panic(err)

		}
		_, err = tx.Exec("DELETE from users WHERE id = $1", id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

	case 2:
		fmt.Print("Enter book id to delete: ")
		fmt.Scan(&id)
		_, err = tx.Exec("DELETE from books WHERE id = $1", id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func getOneRow(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	var id int
	var returningrow ReturningRow
	fmt.Print("Enter user id to print: ")
	fmt.Scan(&id)
	err = tx.QueryRow("SELECT u.id, u.name, u.age, b.name, b.count from users u join books b on u.id = b.user_id WHERE u.id = $1", id).Scan(&returningrow.Id, &returningrow.UserName, &returningrow.UserAge, &returningrow.BookName, &returningrow.BookCount)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	pp.Printf("User Id: %d\tUser Name: %s\tUser age: %d\tBook name: %s\tBook count: %d", returningrow.Id, returningrow.UserName, returningrow.UserAge, returningrow.BookName, returningrow.BookCount)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func getAllData(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	var returningrows []ReturningRow
	rows, err := tx.Query("SELECT u.id, u.name, u.age, b.name, b.count from users u join books b on u.id = b.user_id")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	for rows.Next() {
		var returningrow ReturningRow
		err := rows.Scan(&returningrow.Id, &returningrow.UserName, &returningrow.UserAge, &returningrow.BookName, &returningrow.BookCount)
		if err != nil {
			panic(err)
		}
		returningrows = append(returningrows, returningrow)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	pp.Print(returningrows)
}
