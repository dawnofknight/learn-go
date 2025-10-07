package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type App struct {
	DB *sql.DB
}

func main() {
	// envs or defaults
	dsn := GetDSN()
	// data := &User{}

	// connect
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxIdleTime(2 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	// ping with timeout
	if err := pingWithTimeout(db, 5*time.Second); err != nil {
		log.Fatalf("DB not reachable: %v", err)
	}

	app := &App{DB: db}

	r := SetupRouter(app)

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func pingWithTimeout(db *sql.DB, d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return db.PingContext(ctx)
}

// Handlers

func (a *App) createUser(c *gin.Context) {
	var in User
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	res, err := a.DB.ExecContext(ctx,
		`INSERT INTO users (name, email) VALUES (?, ?)`,
		in.Name, in.Email,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	u, err := a.getUserByID(ctx, uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "created but fetch failed"})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (a *App) listUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	rows, err := a.DB.QueryContext(ctx, `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}
	c.JSON(http.StatusOK, users)
}

func (a *App) getUser(c *gin.Context) {
	id, err := paramID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	u, err := a.getUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (a *App) updateUser(c *gin.Context) {
	id, err := paramID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var in User
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	_, err = a.DB.ExecContext(ctx,
		`UPDATE users SET name = ?, email = ? WHERE id = ?`,
		in.Name, in.Email, id,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := a.getUserByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "updated but fetch failed"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (a *App) deleteUser(c *gin.Context) {
	id, err := paramID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	res, err := a.DB.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	aff, _ := res.RowsAffected()
	c.JSON(http.StatusOK, gin.H{"deleted": aff})
}

// helpers

func (a *App) getUserByID(ctx context.Context, id uint64) (User, error) {
	var u User
	err := a.DB.QueryRowContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func paramID(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
