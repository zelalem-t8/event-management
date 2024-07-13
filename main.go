package main

import (
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"go-graphql-backend/db"
	"go-graphql-backend/graph"
	"go-graphql-backend/graph/generated"
	"go-graphql-backend/graph/model"

	"golang.org/x/crypto/bcrypt"
)

// JWT secret key (replace with your own secure key in production)
var jwtKey = []byte("your_secret_key")

func main() {
	// Database connection
	dsn := "user=postgres password=root dbname=user sslmode=disable"
	err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := gin.Default()

	// CORS middleware configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Replace with your Vue frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	r.POST("/query", gin.WrapH(srv))
	r.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	// Setup routes for login and refresh tokens
	r.POST("/login", handleLogin)
	r.POST("/refresh", handleRefresh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// handleLogin handles user authentication and issues JWT token
func handleLogin(c *gin.Context) {
	var input model.LoginInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var user db.User
	query := `SELECT * FROM users WHERE username = $1`
	if err := db.DB.Get(&user, query, input.Username); err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

// handleRefresh handles refreshing JWT token
func handleRefresh(c *gin.Context) {
	// Implement refresh token logic here if needed
	c.JSON(501, gin.H{"error": "Not implemented"})
}

// generateToken generates a JWT token for the user
func generateToken(user db.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"fullname": user.FullName,
		"email":    user.Email,
		"age":      user.Age,
		"exp":      jwt.TimeFunc().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours)
		"iat":      jwt.TimeFunc().Unix(),                     // Issued At time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
