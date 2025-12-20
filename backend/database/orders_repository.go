package database

import (
	"context"
	"fmt"
	"log"

	"online_store/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOrders(db *pgxpool.Pool, ctx context.Context, UserID int) ([]models.OrdersProfile, error) {
	var Orders []models.OrdersProfile

	query := `SELECT orders.order_id, name, quantity, total_amount
			  FROM order_items
			  JOIN orders
			  ON orders.order_id = order_items.order_id
			  JOIN products
			  ON products.product_id = order_items.product_id
			  WHERE orders.user_id = $1`

	rows, err := db.Query(ctx, query, UserID)
	if err != nil {
		log.Println("Ошибка выполнения запроса в БД", err)
		return Orders, err
	}
	defer rows.Close()
	var Order models.OrdersProfile

	for rows.Next() {
		err := rows.Scan(&Order.OrderID, &Order.ProductName, &Order.Quantity, &Order.TotalAmount)
		if err != nil {
			log.Println("Ошибка сканирования", err)
		}
		Orders = append(Orders, Order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Orders, err
}

func CreateOrder(db *pgxpool.Pool, ctx context.Context, userID int) (int, error) {

	cart_id, err := GetCartIDUser(db, ctx, userID)
	if err != nil {
		log.Println("ошибка : ", err)
		return 0, err
	}

	cart, err := GetUserCart(db, ctx, cart_id)

	if err != nil {
		log.Println("ошибка: ", err)
		return 0, err
	}

	var SumCartProducts float32
	var TotalItems = 0

	for i := range cart {
		cart[i].TotalPrice = cart[i].Price * float32(cart[i].Quantity)
		SumCartProducts += cart[i].TotalPrice
		TotalItems += cart[i].Quantity
	}

	if len(cart) == 0 {
		log.Println("ошибка: заказ пуст")
		return 0, fmt.Errorf("ошибка: заказ пуст")
	}

	queryUpdateQtyProducts := `UPDATE products SET stock_qty = stock_qty - $1 WHERE product_id = $2`

	tx, err := db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		log.Println("ошибка создания транзакции: ", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	for i, cart_item := range cart {
		_, err = tx.Exec(ctx, queryUpdateQtyProducts, cart_item.Quantity, cart_item.ProductID)

		if err != nil {
			fmt.Printf("ошибка обновления количества товара на итерации %d: %s", i, err)
			return 0, err
		}
	}

	queryDeleteCartItems := `DELETE FROM cart_items WHERE cart_id = $1`

	_, err = tx.Exec(ctx, queryDeleteCartItems, cart_id)

	if err != nil {
		log.Println("ошибка очищения корзины: ", err)
		return 0, err
	}

	queryCreateOrder := `INSERT INTO orders(user_id, total_amount)
	VALUES ($1, $2) RETURNING order_id`

	var orderID int
	err = tx.QueryRow(ctx, queryCreateOrder, userID, SumCartProducts).Scan(&orderID)
	if err != nil {
		log.Println("ошибка заполнения таблицы orders: ", err)
		return 0, err
	}

	queryCreateOrderItems := `INSERT INTO order_items(order_id, product_id, quantity, unit_price)
	VALUES ($1, $2, $3, $4) `

	for _, CartItems := range cart {
		_, err = tx.Exec(ctx, queryCreateOrderItems, orderID, CartItems.ProductID, CartItems.Quantity, CartItems.Price)
		if err != nil {
			log.Println("ошибка заполнения таблицы order_items: ", err)
			return 0, err
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Println("Ошибка коммита транзакции:", err)
		return 0, err
	}

	log.Printf("Заказ успешно создан с %d товарами", len(cart))

	return orderID, nil
}
