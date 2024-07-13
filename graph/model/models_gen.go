// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthPayload struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Mutation struct {
}

type Query struct {
}

type SignupInput struct {
	FullName string `json:"fullName"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type User struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}