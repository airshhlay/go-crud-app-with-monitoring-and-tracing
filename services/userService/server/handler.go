package server

import (
	"database/sql"
	"fmt"
	"userService/config"
	constants "userService/constants"
	db "userService/db"
	customErr "userService/errors"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

const (
	createNewUserOp     = "AddUser"
	getUserByUsernameOp = "GetUserByUsername"
	getUserByIDOp       = "GetUserById"
	queryTypeInsert     = "INSERT"
	queryTypeSelect     = "SELECT"
	trueStr             = "true"
	falseStr            = "false"
)

// Handler is a helper called by Server to handle various functions.
// It implements the bulk of the business logic.
type Handler struct {
	config    *config.Config
	dbManager *db.DatabaseManager
	logger    *zap.Logger
}

// CreateNewUser is called by the server to create a new user during Signup.
// First encrypts the user's given password, and inserts the new row into the database.
// If successful, returns the userID.
// Else, returns an error.
func (h *Handler) CreateNewUser(username string, password string) (int64, error) {
	exists, _, err := h.checkUserExists(username)
	if err != nil {
		// error occured when querying database
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorDatabaseQuery,
			ErrorMsg:  constants.ErrorDatabaseQueryMsg,
		}
	}

	// user already exists, return error
	if exists {
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorUserAlreadyExists,
		}
	}

	// encrypt the password
	hash, err := h.getPasswordHash(password)
	if err != nil {
		h.logger.Error(
			constants.ErrorPasswordEncryptionMsg,
			zap.Error(err),
		)
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorPasswordEncryption,
		}
	}

	// insert user into database
	id, err := h.insertNewUser(username, hash)

	h.logger.Info(
		constants.InfoUserAdd,
		zap.String("username", username),
		zap.Int64("id", id),
	)

	return id, nil
}

// VerifyLogin is called by the server when it receives a login request.
// Queries the database for the user,
// and returns an error if the user does not exist
// or if something went wrong when querying the database.
// Verifies the user's given password gainst the hash in the database,
// and returns the userID if passwords match.
// Else, returns an error.
func (h *Handler) VerifyLogin(username string, password string) (int64, error) {
	// retrieve the user
	exists, user, err := h.checkUserExists(username)
	if err != nil {
		// error occured when querying database
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorDatabaseQuery,
			ErrorMsg:  constants.ErrorDatabaseQueryMsg,
		}
	}

	// user does not exist, return error
	if !exists {
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorUserDoesNotExist,
		}
	}

	// check that the given password and stored password match
	err = h.checkPasswordMatch(user.Password, password)
	if err != nil {
		// user password error
		h.logger.Info(
			constants.ErrorUserPasswordMsg,
			zap.String("username", username),
			zap.Error(err),
		)
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorUserPassword,
		}
	}

	// log successful login
	h.logger.Info(
		constants.InfoUserLogin,
		zap.String("username", username),
		zap.Int64("userID", user.UserID),
	)

	// return userID
	return user.UserID, nil
}

// checkUserExists is a helper function.
// Checks if the user exists.
// Returns true if the user exists, and false otherwise.
func (h *Handler) checkUserExists(username string) (bool, db.User, error) {
	user, query, err := h.retrieveUserByUsername(username)

	if err != nil {
		if err == sql.ErrNoRows {
			// query returned no results
			h.logger.Info(
				constants.InfoUserDoesNotExist,
				zap.String("username", username),
				zap.String("query", query),
			)
			return false, db.User{}, nil
		}
		// other unexpected error occurred
		return false, db.User{}, err
	}
	h.logger.Info(
		constants.InfoUserExists,
		zap.String("username", username),
		zap.Int64("userID", user.UserID),
	)
	return true, user, nil
}

// retrieveUserByUsername is a helper function that retrieves users from the database based on username.
// Returns the user if succesful, else returns an error.
// Also returns the query used for logging purposes.
func (h *Handler) retrieveUserByUsername(username string) (db.User, string, error) {
	var user db.User

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s'", username)
	err := h.dbManager.QueryOne(query, getUserByUsernameOp, &user.UserID, &user.Username, &user.Password)
	// err := res.Scan(&user.UserID, &user.Username, &user.Password)

	return user, query, err
}

// insertNewUser is a helper function to insert a new user into the database. It returs the last inserted ID, as well as an error if any.
func (h *Handler) insertNewUser(username string, hash []byte) (int64, error) {
	query := fmt.Sprintf("INSERT INTO users(username, password) VALUES ('%s', '%s')", username, hash)

	id, err := h.dbManager.InsertRow(query, createNewUserOp)
	if err != nil {
		// error occured when inserting user into database
		return 0, &customErr.Error{
			ErrorCode: constants.ErrorDatabaseInsert,
		}
	}

	return id, nil
}

// getPasswordHash uses the bcrypt package to generate a hash from the plaintext password.
func (h *Handler) getPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

// checkPasswordMatch uses the bcrypt package to check that the given plaintext password matches the storedHash.
func (h *Handler) checkPasswordMatch(storedHash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(storedHash, []byte(password))
}
