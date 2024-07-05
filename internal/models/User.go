package models

// User represents a user in the system
// @Description User represents a user in the system
type User struct {
	// ID is the unique identifier of the user
	ID uint `json:"id"`

	// FirstName is the first name of the user
	// @required
	// @example John
	FirstName string `json:"first_name"`

	// LastName is the last name of the user
	// @required
	// @example Doe
	LastName string `json:"last_name"`

	// Email is the email address of the user
	// @required
	// @example john.doe@example.com
	Email string `json:"email"`

	// PassportNumber is the passport number of the user
	// @required
	// @example ABC12345
	PassportNumber string `json:"passport_number"`

	// Task is the task associated with the user
	Task Task `json:"task"`
}
