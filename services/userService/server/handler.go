package server

import (
	"database/sql"
	"fmt"
	"userService/config"
	constants "userService/constants"
	db "userService/db"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

// Handler is a helper called by Server to handle various functions.
// It implements the bulk of the business logic.
type Handler struct {
	config    *config.Config
	dbManager *db.DbManager
	logger    *zap.Logger
}

// Error is a custom struct for errors returned by the service.
// errorCode identifies the type of error that occured.
// errorMsg gives a brief description of the error.
type Error struct {
	errorCode int32
	errorMsg  string
}

// Called by the server to create a new user.
// First encrypts the user's given password, and inserts the new row into the database.
// If successful, returns the userId.
// Else, returns an error.
func (h *Handler) CreateNewUser(username string, password string) (int64, Error) {
	exists, _, err := h.checkUserExists(username)
	if err != nil {
		// error occured when querying database
		return 0, Error{
			errorCode: constants.ERROR_DATABASE_QUERY,
			errorMsg:  constants.ERROR_DATABASE_QUERY_MSG,
		}
	}

	// user already exists, return error
	if exists {
		return 0, Error{
			errorCode: constants.ERROR_USER_ALREADY_EXISTS,
		}
	}

	// encrypt the password
	hash, err := h.getPasswordHash(password)
	if err != nil {
		h.logger.Error(
			constants.ERROR_PASSWORD_ENCRYPTION_MSG,
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_PASSWORD_ENCRYPTION,
		}
	}

	// insert user into database
	query := fmt.Sprintf("INSERT INTO users(username, password) VALUES ('%s', '%s')", username, hash)
	id, err := h.dbManager.InsertRow(query)
	if err != nil {
		// error occured when inserting user into database
		h.logger.Error(
			constants.ERROR_DATABASE_INSERT_MSG,
			zap.String("username", username),
			zap.String("query", query),
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_DATABASE_INSERT,
		}
	}

	h.logger.Info(
		constants.INFO_USER_ADD_MSG,
		zap.String("username", username),
		zap.Int64("id", id),
	)

	return id, Error{}
}

// Called by the server when it receives a login request.
// Queries the database for the user,
// and returns an error if the user does not exist
// or if something went wrong when querying the database.
// Verifies the user's given password gainst the hash in the database,
// and returns the userId if passwords match.
// Else, returns an error.
func (h *Handler) VerifyLogin(username string, password string) (int64, Error) {
	// retrieve the user
	exists, user, err := h.checkUserExists(username)
	if err != nil {
		// error occured when querying database
		return 0, Error{
			errorCode: constants.ERROR_DATABASE_QUERY,
			errorMsg:  constants.ERROR_DATABASE_QUERY_MSG,
		}
	}

	// user does not exist, return error
	if !exists {
		return 0, Error{
			errorCode: constants.ERROR_USER_DOES_NOT_EXIST,
		}
	}

	// check that the given password and stored password match
	err = h.checkPasswordMatch(user.Password, password)
	if err != nil {
		// user password error
		h.logger.Info(
			constants.ERROR_USER_PASSWORD_MSG,
			zap.String("username", username),
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_USER_PASSWORD,
		}
	}

	// log successful login
	h.logger.Info(
		constants.INFO_USER_LOGIN_MSG,
		zap.String("username", username),
		zap.Int64("userId", user.UserId),
	)

	// return userId
	return user.UserId, Error{}
}

// Called by the server during user login/signup.
// Checks if the user exists.
// Returns true if the user exists, and false otherwise.
func (h *Handler) checkUserExists(username string) (bool, db.User, error) {
	user, query, err := h.retrieveUserByUsername(username)

	if err != nil {
		if err == sql.ErrNoRows {
			// query returned no results
			h.logger.Info(
				constants.INFO_USER_DOES_NOT_EXIST,
				zap.String("username", username),
			)
			return false, db.User{}, nil
		}
		// other unexpected error occurred
		h.logger.Error(
			constants.ERROR_DATABASE_QUERY_MSG,
			zap.String("query", query),
			zap.Error(err),
		)
		return false, db.User{}, err
	}
	h.logger.Info(
		constants.INFO_USER_EXISTS,
		zap.String("username", username),
		zap.Int64("userId", user.UserId),
	)
	return true, user, nil
}

// Retrieves users from the database based on username.
// Returns the user if succesful, else returns an error.
// Also returns the query used for logging purposes.
func (h *Handler) retrieveUserByUsername(username string) (db.User, string, error) {
	var user db.User

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s'", username)
	res := h.dbManager.QueryOne(query)
	err := res.Scan(&user.UserId, &user.Username, &user.Password)

	return user, query, err
}

// Uses the bcrypt package to generate a hash from the plaintext password.
func (h *Handler) getPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

// Uses the bcrypt package to check that the given plaintext password matches the storedHash.
func (h *Handler) checkPasswordMatch(storedHash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(storedHash, []byte(password))
}
