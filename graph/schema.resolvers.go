package graph

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go-graphql-backend/db"
	"go-graphql-backend/graph/generated"
	"go-graphql-backend/graph/model"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Define the JWT secret key
var jwtSecret = []byte("your_secret_key")

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Signup(ctx context.Context, input model.SignupInput) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := db.User{
		FullName: input.FullName,
		Username: input.Username,
		Email:    input.Email,
		Age:      input.Age,
		Password: string(hashedPassword),
	}

	query := `INSERT INTO users (full_name, username, email, age, password) 
              VALUES (:full_name, :username, :email, :age, :password) RETURNING id`

	// Using sqlx.Named to prepare the query with named parameters
	stmt, err := db.DB.NamedQuery(query, user)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Executing the query to insert the user and retrieve the ID
	if stmt.Next() {
		err := stmt.Scan(&user.ID)
		if err != nil {
			return nil, err
		}
	}

	// Convert user.ID to string
	// Log user.ID to verify type and value
	//fmt.Printf("User created with ID: %s\n", userID)

	// Return the GraphQL model.User object
	userId, _ := strconv.Atoi(user.ID)
	return &model.User{
		ID:       userId,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	var user db.User
	query := `SELECT * FROM users WHERE username = $1`
	if err := db.DB.Get(&user, query, input.Username); err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, err
	}
	userId, _ := strconv.Atoi(user.ID)
	return &model.AuthPayload{
		Token: token,
		User: &model.User{
			ID:       userId,
			FullName: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Age:      user.Age,
		},
	}, nil
}

func generateToken(user db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
