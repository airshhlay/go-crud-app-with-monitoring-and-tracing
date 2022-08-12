package server

import (
	"database/sql"
	"fmt"
	"userService/config"
	constants "userService/constants"
	db "userService/db"
	customErr "userService/errors"
	metrics "userService/metrics"

	"github.com/prometheus/client_golang/prometheus"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

const (
	CREATE_NEW_USER_OP      = "AddUser"
	GET_USER_BY_USERNAME_OP = "GetUserByUsername"
	GET_USER_BY_ID_OP       = "GetUserById"
	QUERY_TYPE_INSERT       = "INSERT"
	QUERY_TYPE_SELECT       = "SELECT"
	TRUE_STR                = "true"
	FALSE_STR               = "false"
)

// Handler is a helper called by Server to handle various functions.
// It implements the bulk of the business logic.
type Handler struct {
	config    *config.Config
	dbManager *db.DbManager
	logger    *zap.Logger
}

// Called by the server to create a new user.
// First encrypts the user's given password, and inserts the new row into the database.
// If successful, returns the userId.
// Else, returns an error.
func (h *Handler) CreateNewUser(username string, password string) (int64, error) {
	exists, _, err := h.checkUserExists(username)
	if err != nil {
		// error occured when querying database
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_DATABASE_QUERY,
			ErrorMsg:  constants.ERROR_DATABASE_QUERY_MSG,
		}
	}

	// user already exists, return error
	if exists {
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_USER_ALREADY_EXISTS,
		}
	}

	// encrypt the password
	hash, err := h.getPasswordHash(password)
	if err != nil {
		h.logger.Error(
			constants.ERROR_PASSWORD_ENCRYPTION_MSG,
			zap.Error(err),
		)
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_PASSWORD_ENCRYPTION,
		}
	}

	// insert user into database
	id, err := h.insertNewUser(username, hash)

	h.logger.Info(
		constants.INFO_USER_ADD_MSG,
		zap.String("username", username),
		zap.Int64("id", id),
	)

	return id, nil
}

// Called by the server when it receives a login request.
// Queries the database for the user,
// and returns an error if the user does not exist
// or if something went wrong when querying the database.
// Verifies the user's given password gainst the hash in the database,
// and returns the userId if passwords match.
// Else, returns an error.
func (h *Handler) VerifyLogin(username string, password string) (int64, error) {
	// retrieve the user
	exists, user, err := h.checkUserExists(username)
	if err != nil {
		// error occured when querying database
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_DATABASE_QUERY,
			ErrorMsg:  constants.ERROR_DATABASE_QUERY_MSG,
		}
	}

	// user does not exist, return error
	if !exists {
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_USER_DOES_NOT_EXIST,
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
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_USER_PASSWORD,
		}
	}

	// log successful login
	h.logger.Info(
		constants.INFO_USER_LOGIN_MSG,
		zap.String("username", username),
		zap.Int64("userId", user.UserId),
	)

	// return userId
	return user.UserId, nil
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

	// time database query
	querySuccess := TRUE_STR // whether the query was successful
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(h.config.ServiceLabel, QUERY_TYPE_SELECT, GET_USER_BY_USERNAME_OP, querySuccess).Observe(v)
	}))

	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()
	res := h.dbManager.QueryOne(query)
	err := res.Scan(&user.UserId, &user.Username, &user.Password)

	return user, query, err
}

func (h *Handler) insertNewUser(username string, hash []byte) (int64, error) {
	// time database query
	querySuccess := TRUE_STR // whether the query was successful
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(h.config.ServiceLabel, QUERY_TYPE_INSERT, CREATE_NEW_USER_OP, querySuccess).Observe(v)
	}))

	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

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
		querySuccess = FALSE_STR
		return 0, &customErr.Error{
			ErrorCode: constants.ERROR_DATABASE_INSERT,
		}
	}

	return id, nil
}

// Uses the bcrypt package to generate a hash from the plaintext password.
func (h *Handler) getPasswordHash(password string) ([]byte, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.PasswordEncryptionDuration.WithLabelValues().Observe(v)
	}))
	// observe duration at the end of this function
	defer func() {
		timer.ObserveDuration()
	}()

	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

// Uses the bcrypt package to check that the given plaintext password matches the storedHash.
func (h *Handler) checkPasswordMatch(storedHash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(storedHash, []byte(password))
}
