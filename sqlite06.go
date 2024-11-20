package sqlite06

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Filename = ""
)

type UserData struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", Filename)
	if err != nil {
		fmt.Printf("error while openning %s: %s\n", Filename, err)
		return nil, err
	}
	return db, nil
}

func exists(username string) (int, error) {
	db, err := openConnection()
	if err != nil {
		fmt.Printf("error on open connection: %s\n", err)
		return -1, err
	}
	defer db.Close()
	userID := -1
	query := `SELECT ID FROM users WHERE Username='?'`
	rows, err := db.Query(query, username)
	if err != nil {
		fmt.Printf("error while query: %s\n", err)
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("error while scan: %s\n", err)
			return -1, err
		}
		userID = id
	}
	fmt.Printf("found user with id: %d\n", userID)
	return userID, nil
}

func AddUser(user UserData) (int, error) {
	db, err := openConnection()
	if err != nil {
		fmt.Printf("error on add user: %s\n", err)
		return -1, err
	}
	defer db.Close()
	userID, err := exists(user.Username)
	if err != nil {
		fmt.Printf("error on add user: %s\n", err)
		return -1, err
	}
	if userID != -1 {
		fmt.Printf("user with username %s already exists\n", user.Username)
		return -1, fmt.Errorf("user with username %s already exists", user.Username)
	}

	insertQuery := `INSERT INTO users VALUES (NULL,?)`
	_, err = db.Exec(insertQuery, user.Username)
	if err != nil {
		fmt.Printf("error on add user: %s\n", err)
		return -1, err
	}

	userID, err = exists(user.Username)
	if err != nil {
		fmt.Printf("error on add user: %s\n", err)
		return -1, nil
	}

	insertQuery = `INSERT INTO user_data VALUES(?,?,?,?)`
	_, err = db.Exec(insertQuery, userID, user.Name, user.Surname, user.Description)
	if err != nil {
		fmt.Printf("error on add user: %s\n", err)
		return -1, nil
	}
	return userID, nil
}

func DeleteUser(userID int) error {
	db, err := openConnection()
	if err != nil {
		fmt.Printf("error in delete user: %s\n", err)
		return err
	}
	defer db.Close()

	query := `DELETE FROM users WHERE ID = ?`
	_, err = db.Exec(query, userID)
	if err != nil {
		fmt.Printf("error in delete user: %s\n", err)
		return err
	}
	query = `DELETE FROM user_data WHERE UserID = ?`
	_, err = db.Exec(query, userID)
	if err != nil {
		fmt.Printf("error in delete user: %s\n", err)
		return err
	}
	return nil
}

func ListUsers() ([]UserData, error) {
	data := []UserData{}
	db, err := openConnection()
	if err != nil {
		fmt.Printf("error on list users: %s\n", err)
		return nil, err
	}
	defer db.Close()
	query := `SELECT ID, Username, Name, Surname, Description FROM users, user_data WHERE users.ID = user_data.UserID`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("error on list users: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			fmt.Printf("error on list users: %s\n", err)
			return nil, err
		}
		data = append(data, UserData{
			ID:          id,
			Username:    username,
			Name:        name,
			Surname:     surname,
			Description: description,
		})
	}
	return data, nil
}

func UpdateUser(user UserData) error {
	db, err := openConnection()
	if err != nil {
		fmt.Printf("error on update user: %s\n", err)
		return err
	}
	defer db.Close()

	userId, err := exists(user.Username)
	if err != nil {
		fmt.Printf("error on update user: %s\n", err)
		return nil
	}
	if userId == -1 {
		return fmt.Errorf("user %s doesn't exist", user.Username)
	}

	query := `UPDATE user_data SET Name = ?, Surname = ?, Description = ? WHERE UserID = ?`
	_, err = db.Exec(query, user.Name, user.Surname, user.Description, userId)
	if err != nil {
		fmt.Printf("error on update user: %s\n", err)
		return err
	}
	return nil
}
