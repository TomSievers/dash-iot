package user

import "database/sql"

type User struct {
	Id          int64
	Name        string
	Password    []byte
	Permissions string
}

func AddUser(tx *sql.Tx, user User) (int64, error) {
	stmt, err := tx.Prepare("INSERT INTO user (name, password, permissions) VALUES (?, ?, ?)")

	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(user.Name, user.Password, user.Permissions)

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func GetUserByName(tx *sql.Tx, name string) (*User, error) {

	stmt, err := tx.Prepare("SELECT id, name, password, permissions FROM user WHERE name = ?")

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(name)

	user := User{}

	row.Scan(&user.Id, &user.Name, &user.Password, &user.Permissions)

	return nil, nil
}
