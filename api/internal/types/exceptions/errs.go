package exceptions

const (
	// Common errors

	ParameterBindingError = "Parameter binding error."

	// Mock errors

	MockError = "Mock error."

	// User related errors

	PhoneBoundError         = "The phone number has already been bound."
	UserNotExistsError      = "User dose not found."
	InvalidPhoneError       = "Invalid phone number."
	IncorrectPasswordError  = "Password is incorrect"
	AvatarSizeExceededError = "Avatar size exceeded 10MB."

	// Product related errors

	PicSizeExceededError     = "Pic size exceeded 10MB."
	ProductNotFoundError     = "Product not found."
	ProductSoldError         = "Product sold out."
	UserNotProductOwnerError = "User is not the owner of the product."
	ProductOffShelvesError   = "Product is off shelves."

	// Assert related errors

	ResourceNotFoundError = "Resource not found."

	// Order related errors

	ProductNotAvailableError = "Product not available."
	OrderNotFoundError       = "Order not found."
	UserNotOrderOwnerError   = "User is not the owner of the order."
	OrderHasNotShipped       = "Product has not been shipped."
	UserNotOrderSellerError  = "User is not the seller of the order."
	OrderNotPaidError        = "Order has not been paid."
	OrderNotToBePaidError    = "Order's status is not to be paid."
	OrderNotRelatedError     = "You are not related with this order."
	OrderCanNotCanceled      = "Order can not be canceled."
)
