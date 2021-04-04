package users

import (
	"database/sql"
	"github.com/Abacode7/bookstore_utils-go/v2/logger"
	"github.com/Abacode7/bookstore_utils-go/v2/rest_error"
)

const (
	insertUserQuery   = `INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);`
	getUserQuery      = `SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE id=?;`
	findByStatusQuery = `SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;`
	updateUserQuery   = `UPDATE users SET first_name=?, last_name=?, email=?, status=?, password=? WHERE id=?;`
	deleteUserQuery   = `DELETE FROM users WHERE id=?;`
	findByEmailQuery  = `SELECT * FROM users WHERE email = ? AND status = ?;`
)

type IUserDao interface {
	Save(User) (*User, rest_error.RestErr)
	Get(int64) (*User, rest_error.RestErr)
	FindByStatus(string) (Users, rest_error.RestErr)
	Update(User) (*User, rest_error.RestErr)
	Delete(int64) rest_error.RestErr
	FindByEmail(string) (*User, rest_error.RestErr)
}

type userDao struct {
	client *sql.DB
}

/// NewUserDao is a constructor for userDao
func NewUserDao(db *sql.DB) IUserDao {
	return &userDao{db}
}

/// Save stores the user in the database
func (ud *userDao) Save(user User) (*User, rest_error.RestErr) {
	stmt, err := ud.client.Prepare(insertUserQuery)
	if err != nil {
		logger.Error("error preparing insert query", err)
		restErr := rest_error.NewInternalServerError("Prepare err: database error")
		return nil, restErr
	}
	defer stmt.Close()

	result, execErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if execErr != nil {
		logger.Error("error executing prepare query", execErr)
		restErr := rest_error.NewInternalServerError("Exec error: database error")
		return nil, restErr
	}
	userId, queryErr := result.LastInsertId()
	if queryErr != nil {
		logger.Error("error retrieving last insert id", queryErr)
		restErr := rest_error.NewInternalServerError("Result err: database error")
		return nil, restErr
	}
	user.Id = userId

	return &user, nil
}

/// Gets a user with id userID
func (ud *userDao) Get(userID int64) (*User, rest_error.RestErr) {
	stmt, err := ud.client.Prepare(getUserQuery)
	if err != nil {
		logger.Error("error preparing get query", err)
		sqlErr := rest_error.NewInternalServerError("database error")
		return nil, sqlErr
	}
	defer stmt.Close()

	var user User
	row := stmt.QueryRow(userID)
	rowErr := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password)
	if rowErr != nil {
		if rowErr == sql.ErrNoRows {
			return nil, rest_error.NewNotFoundError("invalid user id: user not found")
		}
		logger.Error("error scanning user data", rowErr)
		return nil, rest_error.NewInternalServerError("database error")
	}
	return &user, nil
}

/// Update modifies the values of a user with specified id
func (ud *userDao) Update(user User) (*User, rest_error.RestErr) {
	stmt, prepErr := ud.client.Prepare(updateUserQuery)
	if prepErr != nil {
		logger.Error("error preparing update query", prepErr)
		return nil, rest_error.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, stmtErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)
	if stmtErr != nil {
		logger.Error("error when trying to update user", stmtErr)
		return nil, rest_error.NewInternalServerError("database error")
	}
	return &user, nil
}

/// Delete removes user with userId from the database
func (ud *userDao) Delete(userId int64) rest_error.RestErr {
	stmt, prepErr := ud.client.Prepare(deleteUserQuery)
	if prepErr != nil {
		logger.Error("error preparing delete query", prepErr)
		return rest_error.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result, stmtErr := stmt.Exec(userId)
	if stmtErr != nil {
		logger.Error("error executing delete query", stmtErr)
		return rest_error.NewInternalServerError("database error")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		logger.Error("error retrieving rows affected", err)
		return rest_error.NewInternalServerError("database error")
	}
	if rowsAff < 1 {
		return rest_error.NewBadRequestError("user with id doesnt exist")
	}
	return nil
}

/// FindByStatus gets all users with given status
func (ud *userDao) FindByStatus(status string) (Users, rest_error.RestErr) {
	stmt, prepErr := ud.client.Prepare(findByStatusQuery)
	if prepErr != nil {
		logger.Error("error preparing findByStatus query", prepErr)
		return nil, rest_error.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, stmtErr := stmt.Query(status)
	if stmtErr != nil {
		logger.Error("error executing findByStatus query", stmtErr)
		return nil, rest_error.NewInternalServerError("database error")
	}
	defer rows.Close()

	users := make(Users, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
		if err != nil {
			logger.Error("error scanning retrieved data", stmtErr)
			return nil, rest_error.NewInternalServerError("database error")
		}
		users = append(users, user)
	}
	return users, nil
}

/// FindByEmailAndPassword gets the user with given email and password
func (ud *userDao) FindByEmail(email string) (*User, rest_error.RestErr) {
	stmt, prepErr := ud.client.Prepare(findByEmailQuery)
	if prepErr != nil {
		logger.Error("error executing findByEmailAndPassword query", prepErr)
		err := rest_error.NewInternalServerError("database error")
		return nil, err
	}
	rows := stmt.QueryRow(email, StatusActive)
	var user User
	err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("user not found", err)
			return nil, rest_error.NewNotFoundError("invalid data: user not found")
		}
		logger.Error("error scanning retrieved data", err)
		return nil, rest_error.NewInternalServerError("database error")
	}
	return &user, nil
}
