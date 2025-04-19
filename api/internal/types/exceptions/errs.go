package exceptions

const (
	// Common errors

	ParameterBindingError = "Parameter binding error."

	// Mock errors

	MockError = "Mock error."

	// User related errors

	PhoneBoundError        = "The phone number has already been bound."
	UserNotExistsError     = "User dose not found."
	InvalidPhoneError      = "Invalid phone number."
	IncorrectPasswordError = "Password is incorrect"
)
