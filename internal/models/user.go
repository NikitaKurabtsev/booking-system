package models

type User struct {
	ID           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username,omitempty"`
	Email        string `db:"email" json:"email,omitempty"`
	PasswordHash string `db:"password_hash"`
}
