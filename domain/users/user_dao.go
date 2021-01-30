package users

import (
	"database/sql"
	"github.com/Abacode7/bookstore_users-api/utils/errors"
	"github.com/Abacode7/bookstore_users-api/utils/logger"
)

const (
	insertQuery                 = `INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);`
	getQuery                    = `SELECT * FROM users WHERE id=?;`
	findQuery                   = `SELECT * FROM users WHERE status=?;`
	updateQuery                 = `UPDATE users SET first_name=?, last_name=?, email=?, status=?, password=? WHERE id=?;`
	deleteQuery                 = `DELETE FROM users WHERE id=?;`
	findByEmailQuery = `SELECT * FROM users WHERE email = ? AND status = ?;`
)

type IUserDao interface {
	Save(User) (*User, *errors.RestErr)
	Get(int64) (*User, *errors.RestErr)
	FindByStatus(string) ([]User, *errors.RestErr)
	Update(User) (*User, *errors.RestErr)
	Delete(int64) *errors.RestErr
	FindByEmail(string) (*User, *errors.RestErr)
}

type userDao struct {
	client *sql.DB
}

/// NewUserDao is a constructor for userDao
func NewUserDao(db *sql.DB) IUserDao {
	return &userDao{db}
}

/// Save stores the user in the database
func (ud *userDao) Save(user User) (*User, *errors.RestErr) {
	stmt, err := ud.client.Prepare(insertQuery)
	if err != nil {
		logger.Error("error preparing insert query", err)
		restErr := errors.NewInternalServerError("Prepare err: database error")
		return nil, restErr
	}
	defer stmt.Close()

	result, execErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if execErr != nil {
		logger.Error("error executing prepare query", execErr)
		restErr := errors.NewInternalServerError("Exec error: database error")
		return nil, restErr
	}
	userId, queryErr := result.LastInsertId()
	if queryErr != nil {
		logger.Error("error retrieveing last insert id", queryErr)
		restErr := errors.NewInternalServerError("Result err: database error")
		return nil, restErr
	}
	user.Id = userId

	return &user, nil
}

/// Gets a user with id userID
func (ud *userDao) Get(userID int64) (*User, *errors.RestErr) {
	stmt, err := ud.client.Prepare(getQuery)
	if err != nil {
		logger.Error("error preparing get query", err)
		sqlErr := errors.NewInternalServerError("database error")
		return nil, sqlErr
	}
	defer stmt.Close()

	rows, queryErr := stmt.Query(userID)
	if queryErr != nil {
		logger.Error("error getting user with userId", queryErr)
		err := errors.NewInternalServerError("database error")
		return nil, err
	}

	count := 0
	var user User
	if rows.Next() {
		rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password)
		count++
	}
	if count < 1 {
		err := errors.NewBadRequestError("user with id doesnt exist")
		return nil, err
	}
	return &user, nil
}

/// Update modifies the values of a user with specified id
func (ud *userDao) Update(user User) (*User, *errors.RestErr) {
	stmt, prepErr := ud.client.Prepare(updateQuery)
	if prepErr != nil {
		logger.Error("error preparing update query", prepErr)
		err := errors.NewInternalServerError("database error")
		return nil, err
	}
	defer stmt.Close()

	result, stmtErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)
	if stmtErr != nil {
		logger.Error("error when trying to update user", stmtErr)
		err := errors.NewInternalServerError("database error")
		return nil, err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		logger.Error("error when trying to retrieve rows affected", err)
		rowErr := errors.NewInternalServerError("database error")
		return nil, rowErr
	}
	if rowsAff < 1 {
		err := errors.NewBadRequestError("user with id doesn't exist")
		return nil, err
	}
	return &user, nil
}

/// Delete removes user with userId from the database
func (ud *userDao) Delete(userId int64) *errors.RestErr {
	stmt, prepErr := ud.client.Prepare(deleteQuery)
	if prepErr != nil {
		logger.Error("error preparing delete query", prepErr)
		err := errors.NewInternalServerError("database error")
		return err
	}
	defer stmt.Close()

	result, stmtErr := stmt.Exec(userId)
	if stmtErr != nil {
		logger.Error("error executing delete query", stmtErr)
		err := errors.NewInternalServerError("database error")
		return err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		logger.Error("error retrieving rows affected", err)
		er := errors.NewInternalServerError("database error")
		return er
	}
	if rowsAff < 1 {
		err := errors.NewBadRequestError("user with id doesnt exist")
		return err
	}
	return nil
}

/// FindByStatus gets all users with given status
func (ud *userDao) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, prepErr := ud.client.Prepare(findQuery)
	if prepErr != nil {
		logger.Error("error preparing findByStatus query", prepErr)
		err := errors.NewInternalServerError("database error")
		return nil, err
	}
	defer stmt.Close()

	rows, stmtErr := stmt.Query(status)
	if stmtErr != nil {
		logger.Error("error executing findByStatus query", stmtErr)
		err := errors.NewInternalServerError("database error")
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password)
		users = append(users, user)
	}
	return users, nil
}

/// FindByEmailAndPassword gets the user with given email and password
func (ud *userDao) FindByEmail(email string) (*User, *errors.RestErr) {
	stmt, prepErr := ud.client.Prepare(findByEmailQuery)
	if prepErr != nil {
		logger.Error("error executing findByEmailAndPassword query", prepErr)
		err := errors.NewInternalServerError("database error")
		return nil, err
	}
	rows := stmt.QueryRow(email, StatusActive)
	var user User
	err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("user not found", err)
			return nil, errors.NewNotFoundError("invalid data: user not found")
		}
		logger.Error("error scanning retrieved data", err)
		return nil, errors.NewInternalServerError("database error")
	}
	return &user, nil
}
