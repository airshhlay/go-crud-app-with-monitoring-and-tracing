package errors

// Error is a custom struct for errors returned by the service.
// errorCode identifies the type of error that occured.
// errorMsg gives a brief description of the error.
type Error struct {
	ErrorCode int32
	ErrorMsg  string
	Err       error
}

// Returns Message if Err is nil
func (err Error) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.ErrorMsg
}
