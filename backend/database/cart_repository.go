package database

import (
	"context"
	"fmt"
	"log"

	"online_store/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserCart(db *pgxpool.Pool, ctx context.Context, cart_id int) ([]models.CartItem, error) {
	query := `
        SELECT 
            cart_items.product_id, 
            cart_items.quantity, 
            products.name, 
            products.description, 
            products.price, 
            products.image_url,
            categories.name as category_name,
            producers.name as producer_name
        FROM cart_items
        INNER JOIN products ON products.product_id = cart_items.product_id
        INNER JOIN products_categories ON products.product_id = products_categories.product_id
        INNER JOIN categories ON categories.category_id = products_categories.category_id
        INNER JOIN products_producers ON products.product_id = products_producers.product_id
        INNER JOIN producers ON producers.producer_id = products_producers.producer_id
        WHERE cart_items.cart_id = $1`

	var spisok []models.CartItem
	rows, err := db.Query(ctx, query, cart_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var stroka models.CartItem
		err = rows.Scan(
			&stroka.ProductID,
			&stroka.Quantity,
			&stroka.ProductName,
			&stroka.Description,
			&stroka.Price,
			&stroka.ImageURL,
			&stroka.CategoryName,
			&stroka.ProducerName,
		)

		if err != nil {
			log.Println("Ошибка сканирования:", err)
			continue
		}
		spisok = append(spisok, stroka)
	}
	return spisok, nil
}

func GetCartIDUser(db *pgxpool.Pool, ctx context.Context, userID int) (int, error) {

	var cart_id int
	query := `SELECT cart_id FROM carts WHERE user_id = $1`

	err := db.QueryRow(ctx, query, userID).Scan(&cart_id)

	if err != nil {
		// Если пользователь не найден
		if err.Error() == "no rows in result set" {
			return -1, fmt.Errorf("user not found")
		}
		// Другие ошибки базы данных
		return -1, fmt.Errorf("database error: %w", err)
	}

	return cart_id, nil
}

func CreateCartForUser(db *pgxpool.Pool, ctx context.Context, userID int) (int, error) {

	query := `INSERT INTO carts(user_id) VALUES ($1) RETURNING cart_id`
	var cartID int
	err := db.QueryRow(ctx, query, userID).Scan(&cartID)
	return cartID, err
}
func AddProductToCart(db *pgxpool.Pool, ctx context.Context, cartID int, productID int, quantity int) error {

	query :=
		`INSERT INTO cart_items(cart_id, product_id, quantity)
	VALUES ($1, $2, $3) ON CONFLICT (cart_id, product_id) 
	DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`

	_, err := db.Exec(ctx, query, cartID, productID, quantity)
	return err

}

func DeleteProductFromCart(db *pgxpool.Pool, ctx context.Context, cartID int, productID int) error {
	query := `DELETE FROM cart_items WHERE product_id = $1`

	_, err := db.Exec(ctx, query, productID)
	return err
}
