package models

import "gorm.io/gorm"

type CreateUserDTO struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
}

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type AnonymousChangePasswordDTO struct {
	NewPassword string `json:"new_password"`
}

type AuthenticationTypeDTO struct {
	Otp          string
	IsSmsEnabled bool
}
type UpdateUserDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserDTO struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	PhoneNumber        string `json:"phone_number"`
	IsTwoFactorEnabled bool   `json:"two_factor_enabled"`
	IsLockoutEnabled   bool   `json:"is_lockout_enabled"`
}

type User struct {
	gorm.Model
	FirstName            string `Gorm:"size:255"`
	LastName             string `Gorm:"size:255"`
	Username             string `Gorm:"unique"`
	Password             string
	Email                string `Gorm:"unique"`
	PhoneNumber          string
	EmailConfirmed       bool `Gorm:"default:false"`
	PhoneNumberConfirmed bool `Gorm:"default:false"`
	Role                 int
	TwoFactorEnabled     bool   `Gorm:"default:false"`
	AuthenticationKey    string `Gorm:"size:255"`
	AuthenticationType   int    `Gorm:"default:0"`
	LockoutEnabled       bool   `Gorm:"default:false"`
	RefreshToken         string `Gorm:"size:255"`
	Tasks                []Task `Gorm:"foreignKey:UserID,constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
}

type Role int

const (
	UserRole Role = iota
	Admin
)

var Roles = map[Role]string{
	UserRole: "user",
	Admin:    "admin",
}

type AuthenticationType int

const (
	Sms AuthenticationType = iota
	Google
	Default
)

var AuthenticationTypes = map[AuthenticationType]string{
	Sms:     "sms",
	Google:  "google_Authentication",
	Default: "none",
}
