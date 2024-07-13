// db/models.go
package db

type User struct {
	ID       string `db:"id"`
	FullName string `db:"full_name"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Age      int    `db:"age"`
	Password string `db:"password"`
}
