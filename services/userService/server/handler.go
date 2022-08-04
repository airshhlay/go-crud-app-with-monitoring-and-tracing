package server

import (
	"database/sql"
	"fmt"
	"userService/config"
	"userService/constants"
	"userService/db"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	config    *config.Config
	dbManager *db.DbManager
	logger    *zap.Logger
}

type Error struct {
	errorCode int32
	errorMsg  string
}

// checks if the user exists in the database
func (h *Handler) CheckUserExists(username string) (bool, Error) {
	_, err := h.retrieveUserByUsername(username)

	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Info(
				"User does not exist, able to create account",
				zap.String("username", username),
			)
			return false, Error{}
		}
		return false, Error{
			errorCode: constants.ERROR_DATABASE,
			errorMsg:  constants.ERROR_DATABASE_MSG,
		}
	}
	return true, Error{}
}

func (h *Handler) retrieveUserByUsername(username string) (db.User, error) {
	var user db.User

	query := fmt.Sprintf("SELECT * FROM users WHERE username=%s", username)
	res := h.dbManager.QueryOne(query)
	err := res.Scan(&user.UserId, &user.Username, &user.Password)

	if err != nil {
		h.logger.Error(
			"Error occured when querying database for user",
			zap.String("query", query),
			zap.Error(err),
		)
	}

	return user, err
}

func (h *Handler) CreateNewUser(username string, password string) (int64, Error) {
	// encrypt the password
	hash, err := h.getPasswordHash(password)
	if err != nil {
		h.logger.Error(
			"Error occured when encrypting password",
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_PASSWORD_ENCRYPTION,
			errorMsg:  constants.ERROR_PASSWORD_ENCRYPTION_MSG,
		}
	}
	query := fmt.Sprintf("INSERT INTO users(username, password) VALUES (%s, %b)", username, hash)
	id, err := h.dbManager.InsertRow(query)
	if err != nil {
		// error occured when inserting user into database
		h.logger.Error(
			"Error occured when inserting user into database",
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_DATABASE,
			errorMsg:  constants.ERROR_DATABASE_MSG,
		}
	}

	h.logger.Info(
		"User successfully added to database",
		zap.String("username", username),
		zap.Int64("id", id),
	)

	return id, Error{}
}

func (h *Handler) VerifyLogin(username string, password string) (int64, Error) {
	// retrieve the user
	user, err := h.retrieveUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			// user does not exist
			h.logger.Info(
				"User does not exist",
				zap.String("username", username),
			)
			return 0, Error{
				errorCode: constants.ERROR_USER_DOES_NOT_EXIST,
				errorMsg:  constants.ERROR_USER_DOES_NOT_EXIST_MSG,
			}
		}
		// unexpected error occured when querying database
		return 0, Error{
			errorCode: constants.ERROR_DATABASE,
			errorMsg:  constants.ERROR_DATABASE_MSG,
		}
	}

	// generate the hash from the given password
	hash, err := h.getPasswordHash(password)
	if err != nil {
		// encryption error
		h.logger.Error(
			"Encryption error",
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_PASSWORD_ENCRYPTION,
			errorMsg:  constants.ERROR_PASSWORD_ENCRYPTION_MSG,
		}
	}

	// check that the given password and stored password match
	err = h.checkPasswordMatch(user.Password, hash)
	if err != nil {
		// user password error
		h.logger.Info(
			"User attempted login with password mismatch",
			zap.String("username", username),
			zap.Error(err),
		)
		return 0, Error{
			errorCode: constants.ERROR_USER_PASSWORD,
			errorMsg:  constants.ERROR_USER_PASSWORD_MSG,
		}
	}

	return user.UserId, Error{}
}

func (h *Handler) getPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (h *Handler) checkPasswordMatch(password []byte, hash []byte) error {
	return bcrypt.CompareHashAndPassword(password, hash)
}
