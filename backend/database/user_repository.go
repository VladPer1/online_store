package database

import (
	"context"
	"fmt"
	"log"
	"online_store/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(db *pgxpool.Pool, ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (email, password, name, status) VALUES ($1, $2, $3, $4 ) RETURNING id`
	err := db.QueryRow(ctx, query, user.Email, user.Password, user.Name, user.Status).Scan(&user.ID)
	if err != nil {
		log.Println("Ошибка создания пользователя:", err)
	}

	return err
}

func GetUser(db *pgxpool.Pool, ctx context.Context, id int) (*models.User, error) {

	var user models.User

	query := `SELECT id, email, password, name, status FROM users WHERE id=$1`

	err := db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status)
	if err != nil {
		// Если пользователь не найден
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		// Другие ошибки базы данных
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil

}

func GetUserByEmail(db *pgxpool.Pool, ctx context.Context, email string) (*models.User, error) {

	var user models.User

	query := `SELECT id, email, password, name, status FROM users WHERE email=$1`
	// Определяем запрос на основе идентификатора

	// Выполняем запрос и проверяем ошибки
	err := db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status)
	if err != nil {
		// Если пользователь не найден
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		// Другие ошибки базы данных
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil

}

func UpdateUser(db *pgxpool.Pool, ctx context.Context, user *models.User) error {
	query := `UPDATE users SET email=$1, password=$2, name=$3, status =$4 WHERE id=$5`
	_, err := db.Exec(ctx, query, user.Email, user.Password, user.Name, user.Status, user.ID)
	return err
}

func DeleteUser(db *pgxpool.Pool, ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := db.Exec(ctx, query, id)
	return err
}

func CreatePlaceholderUser(db *pgxpool.Pool, ctx context.Context, user *models.PlaceholderUser) error {

	query := `INSERT INTO pending_registrations (email, username, password, verification_code) 
              VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(ctx, query, user.Email, user.Name, user.Password, user.VerificationCode)

	if err != nil {
		log.Println("Ошибка создания пользователя:", err)
	}

	return err
}

func GetPlaceholderUser(db *pgxpool.Pool, ctx context.Context, email string) (*models.PlaceholderUser, error) {
	var user models.PlaceholderUser

	query := `SELECT email, username, password, verification_code
			 FROM pending_registrations WHERE email=$1`

	// Выполняем запрос и проверяем ошибки
	err := db.QueryRow(ctx, query, email).Scan(&user.Email, &user.Name, &user.Password, &user.VerificationCode)

	if err != nil {
		// Если пользователь не найден
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		// Другие ошибки базы данных
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

func DeletePlaceholderUser(db *pgxpool.Pool, ctx context.Context, email string) error {
	query := `DELETE FROM pending_registrations WHERE email=$1`
	_, err := db.Exec(ctx, query, email)
	return err
}

func CreateUserInFPD(db *pgxpool.Pool, ctx context.Context, user *models.PlaceholderUser) error {

	query := `INSERT INTO forgot_password_date (email, verification_code) 
              VALUES ($1, $2)`
	_, err := db.Exec(ctx, query, user.Email, user.VerificationCode)

	if err != nil {
		log.Println("Ошибка создания пользователя:", err)
	}

	return err
}

func GetUserFromFPD(db *pgxpool.Pool, ctx context.Context, email string) (*models.PlaceholderUser, error) {
	var user models.PlaceholderUser

	query := `SELECT email, verification_code
			 FROM forgot_password_date WHERE email=$1`

	// Выполняем запрос и проверяем ошибки
	err := db.QueryRow(ctx, query, email).Scan(&user.Email, &user.VerificationCode)

	if err != nil {
		// Если пользователь не найден
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		// Другие ошибки базы данных
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

func DeleteUserFromFPD(db *pgxpool.Pool, ctx context.Context, email string) error {
	query := `DELETE FROM forgot_password_date WHERE email=$1`
	_, err := db.Exec(ctx, query, email)
	return err
}
