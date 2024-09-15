package entities

import (
	"errors"
	"regexp"
	"time"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// UserRole represents predefined roles a user can have
type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
	UserRoleGuest    UserRole = "guest"
	// Additional roles can be added here
)

// User represents a user in the system
type User struct {
	UserID    int64      // Unique identifier for the user
	Username  string     // Username of the user
	Email     string     // Email address
	Password  string     // Hashed password (never store plain text passwords)
	Status    UserStatus // Status of the user account
	Roles     []UserRole // Roles assigned to the user
	CreatedAt time.Time  // Timestamp when the user was created
	Age       int        // Age of the user
}

// NewUser creates a new User with necessary validations
func NewUser(userID int64, username, email, password string, roles []UserRole, age int) (*User, error) {
	user := &User{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Password:  password, // In practice, hash the password before storing
		Status:    UserStatusActive,
		Roles:     roles,
		CreatedAt: time.Now(),
		Age:       age,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate ensures all business invariants for the User are met
func (u *User) Validate() error {
	if u.UserID <= 0 {
		return errors.New("user ID must be a positive integer")
	}
	if u.Username == "" {
		return errors.New("username must not be empty")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasRequiredPasswordComplexity(u.Password) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	if !isValidUserStatus(u.Status) {
		return errors.New("invalid user status")
	}
	if len(u.Roles) == 0 {
		return errors.New("user must have at least one role")
	}
	for _, role := range u.Roles {
		if !isValidUserRole(role) {
			return errors.New("invalid user role: " + string(role))
		}
	}
	if u.Age < 18 {
		return errors.New("user must be at least 18 years old")
	}
	return nil
}

// UpdateEmail updates the user's email after validation
func (u *User) UpdateEmail(newEmail string) error {
	if !isValidEmail(newEmail) {
		return errors.New("invalid email format")
	}
	u.Email = newEmail
	return nil
}

// ChangePassword updates the user's password after validation
func (u *User) ChangePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasRequiredPasswordComplexity(newPassword) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	u.Password = newPassword // In practice, hash the password before storing
	return nil
}

// UpdateStatus changes the user's status following business rules
func (u *User) UpdateStatus(newStatus UserStatus) error {
	if !isValidUserStatus(newStatus) {
		return errors.New("invalid user status")
	}
	// Business rule: Suspended users can only be reactivated
	if u.Status == UserStatusSuspended && newStatus != UserStatusActive {
		return errors.New("suspended users can only be reactivated")
	}
	u.Status = newStatus
	return nil
}

// AddRole adds a new role to the user
func (u *User) AddRole(role UserRole) error {
	if !isValidUserRole(role) {
		return errors.New("invalid user role")
	}
	for _, r := range u.Roles {
		if r == role {
			return errors.New("user already has the role: " + string(role))
		}
	}
	u.Roles = append(u.Roles, role)
	return nil
}

// RemoveRole removes a role from the user
func (u *User) RemoveRole(role UserRole) error {
	if !isValidUserRole(role) {
		return errors.New("invalid user role")
	}
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			return nil
		}
	}
	return errors.New("user does not have the role: " + string(role))
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// Helper function to validate password complexity
func hasRequiredPasswordComplexity(password string) bool {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper && hasLower && hasNumber
}

// Helper function to validate user status
func isValidUserStatus(status UserStatus) bool {
	return status == UserStatusActive || status == UserStatusInactive || status == UserStatusSuspended
}

// Helper function to validate user roles
func isValidUserRole(role UserRole) bool {
	return role == UserRoleAdmin || role == UserRoleCustomer || role == UserRoleGuest
}
