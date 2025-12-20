package database

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"online_store/models"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckAdminStatus(db *pgxpool.Pool, ctx context.Context, userID int) (bool, error) {

	var status string
	query := `SELECT status FROM users WHERE id = $1`

	err := db.QueryRow(ctx, query, userID).Scan(&status)

	if err != nil {
		log.Println("Ошибка получения статуса Администратора", err)
		return false, err
	}

	return status == "admin", nil

}

func GetAllProductsForAdmin(db *pgxpool.Pool, ctx context.Context) ([]models.ProductForAdmin, error) {

	query :=
		`SELECT products.product_id, products.name, products.description, price, stock_qty, 
		image_url, producers.producer_id, categories.category_id, producers.name, categories.name
	 FROM products
	 LEFT JOIN products_producers 
	 ON products_producers.product_id =products.product_id 
	 LEFT JOIN producers 
	 ON products_producers.producer_id =producers.producer_id
	 LEFT JOIN products_categories 
	 ON products_categories.product_id =products.product_id
	 LEFT JOIN categories 
	 ON products_categories.category_id =categories.category_id
	 ORDER BY products.product_id`
	rows, err := db.Query(ctx, query)

	if err != nil {
		log.Println("Ошибка получения списка продуктов: ", err)
		return nil, err
	}

	defer rows.Close()
	var Products []models.ProductForAdmin

	for rows.Next() {
		var Product models.ProductForAdmin

		err := rows.Scan(
			&Product.ProductID, &Product.Name, &Product.Description,
			&Product.Price, &Product.StockQuantity, &Product.ImageURL,
			&Product.ProducerID, &Product.CategoryID,
			&Product.ProducerName, &Product.CategoryName)

		if err != nil {
			log.Println("Ошибка получения списка продуктов: ", err)
			return nil, err
		}
		Products = append(Products, Product)
	}
	return Products, nil
}

func GetProductForAdmin(db *pgxpool.Pool, ctx context.Context, ProductID int) (models.ProductForAdmin, error) {

	var Product models.ProductForAdmin
	query :=
		`SELECT products.product_id, products.name, products.description, price, stock_qty, 
		image_url, producers.producer_id, categories.category_id, producers.name, categories.name
	 FROM products
	 LEFT JOIN products_producers 
	 ON products_producers.product_id =products.product_id 
	 LEFT JOIN producers 
	 ON products_producers.producer_id =producers.producer_id
	 LEFT JOIN products_categories 
	 ON products_categories.product_id =products.product_id
	 LEFT JOIN categories 
	 ON products_categories.category_id =categories.category_id
	 WHERE products.product_id = $1`
	err := db.QueryRow(ctx, query, ProductID).Scan(&Product.ProductID, &Product.Name, &Product.Description,
		&Product.Price, &Product.StockQuantity, &Product.ImageURL, &Product.ProducerID,
		&Product.CategoryID, &Product.ProducerName, &Product.CategoryName)

	if err != nil {
		log.Println("Ошибка получения продукта из бд: ", err)
		return Product, err
	}

	return Product, nil
}

func UpdateProductForAdmin(db *pgxpool.Pool, ctx context.Context, product models.ProductForAdmin) error {

	tx, err := db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		log.Println("Ошибка начала транзакции:", err)
		return err
	}
	defer tx.Rollback(ctx)
	queryUpdateProducts := `UPDATE products 
							SET name = $1, description = $2, price = $3, 
							stock_qty = $4, image_url = $5
							WHERE product_id = $6`
	// 1. Обновляем основную информацию о товаре
	_, err = tx.Exec(ctx, queryUpdateProducts,
		product.Name, product.Description, product.Price,
		product.StockQuantity, product.ImageURL, product.ProductID)

	if err != nil {
		log.Println("Ошибка обновления products:", err)
		return err
	}
	queryUpdatePC := `
        UPDATE products_categories 
        SET category_id = $1 
        WHERE product_id = $2`
	// 2. Обновляем категорию товара
	_, err = tx.Exec(ctx, queryUpdatePC,
		product.CategoryID, product.ProductID)

	if err != nil {
		log.Println("Ошибка обновления products_categories:", err)
		return err
	}
	queryUpdatePP := `UPDATE products_producers 
					SET producer_id = $1 
					WHERE product_id = $2`
	// 3. Обновляем производителя товара
	_, err = tx.Exec(ctx, queryUpdatePP,
		product.ProducerID, product.ProductID)

	if err != nil {
		log.Println("Ошибка обновления products_producers:", err)
		return err
	}

	// Если все успешно - коммитим транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Ошибка коммита транзакции:", err)
		return err
	}

	log.Printf("Товар %d успешно обновлен", product.ProductID)
	return nil
}

func DeleteProductForAdmin(db *pgxpool.Pool, ctx context.Context, ProductID int) error {
	//Реализовано каскадное удаление
	query := `DELETE FROM Products WHERE product_id=$1`
	_, err := db.Exec(ctx, query, ProductID)
	return err
}

func CreateProductForAdmin(db *pgxpool.Pool, ctx context.Context, Product models.ProductForAdmin) error {

	queryInsertProducts := `INSERT INTO products(name, description, price, stock_qty, image_url) 
	VALUES ($1,$2,$3,$4,$5) RETURNING product_id`

	tx, err := db.Begin(ctx)

	if err != nil {
		log.Println("Ошибка начала транзакции:", err)
		return err
	}
	defer tx.Rollback(ctx)

	var productID int
	err = tx.QueryRow(ctx, queryInsertProducts,
		Product.Name,
		Product.Description,
		Product.Price,
		Product.StockQuantity,
		Product.ImageURL,
	).Scan(&productID)

	if err != nil {
		log.Println("Ошибка вставка в таблицу products:", err)
		return err
	}

	queryInsertProductsProducers := `INSERT INTO products_producers VALUES ($1,$2)`

	_, err = tx.Exec(ctx, queryInsertProductsProducers, productID, Product.ProducerID)

	if err != nil {
		log.Println("Ошибка вставка в таблицу products_producers:", err)
		return err
	}

	queryInsertProductsCategories := `INSERT INTO products_categories VALUES ($1,$2)`

	_, err = tx.Exec(ctx, queryInsertProductsCategories, productID, Product.CategoryID)

	if err != nil {
		log.Println("Ошибка вставка в таблицу products_categories :", err)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Println("Ошибка commit транзакции :", err)
		return err
	}

	return nil
}

func MakeAdmin(db *pgxpool.Pool, ctx context.Context, userID int) error {
	query := `UPDATE users SET status = 'admin' WHERE id = $1`
	_, err := db.Exec(ctx, query, userID)
	return err
}

func UnbanUser(db *pgxpool.Pool, ctx context.Context, userID int) error {
	query := `UPDATE users SET status = 'user' WHERE id = $1`
	_, err := db.Exec(ctx, query, userID)
	return err
}

func BanUser(db *pgxpool.Pool, ctx context.Context, userID int) error {
	query := `UPDATE users SET status = 'banned' WHERE id = $1`
	_, err := db.Exec(ctx, query, userID)
	return err
}

func ExportProductsToCSV(db *pgxpool.Pool, ctx context.Context, writer *csv.Writer) error {
	// Получаем данные из products
	rows, err := db.Query(ctx, `
		SELECT product_id, name, description, price, stock_qty, image_url
		FROM products 
		ORDER BY product_id
	`)
	if err != nil {
		log.Println("Ошибка получения данных:", err)
		return err
	}
	defer rows.Close()

	// Записываем данные
	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.ProductID, &product.Name, &product.Description,
			&product.Price, &product.StockQuantity, &product.Image)
		if err != nil {
			log.Println("Ошибка сканирования:", err)
			return err
		}

		// Преобразуем Price в string (если он numeric в БД)
		priceStr := strconv.FormatFloat(float64(product.Price), 'f', 2, 64)

		record := []string{
			strconv.Itoa(product.ProductID),
			product.Name,
			product.Description,
			priceStr, // ← Исправлено: преобразуем в string
			strconv.Itoa(product.StockQuantity),
			product.Image,
		}

		if err := writer.Write(record); err != nil {
			log.Println("Ошибка записи строки CSV:", err)
			return err
		}
	}

	return nil
}

func ImportProductsCSV(db *pgxpool.Pool, ctx context.Context, file io.Reader) (int, error) {
	// Читаем CSV
	reader := csv.NewReader(file)
	reader.Comma = ','

	// Пропускаем заголовок
	headers, err := reader.Read()
	if err != nil {
		log.Println("Ошибка чтения заголовков CSV:", err)
		return 0, err
	}
	log.Printf("Заголовки CSV: %v", headers)

	// Начинаем транзакцию
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("Ошибка начала транзакции:", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Импортируем данные
	imported := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Ошибка чтения строки CSV:", err)
			continue
		}

		if len(record) < 6 {
			log.Println("Пропущена неполная строка:", record)
			continue
		}

		// Парсим данные из CSV
		productID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Ошибка парсинга product_id:", record[0])
			continue
		}

		price, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			log.Println("Ошибка парсинга price:", record[3])
			continue
		}

		stockQty, err := strconv.Atoi(record[4])
		if err != nil {
			log.Println("Ошибка парсинга stock_qty:", record[4])
			continue
		}

		// Вставляем или обновляем данные
		_, err = tx.Exec(ctx, `
			INSERT INTO products (product_id, name, description, price, stock_qty, image_url)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (product_id) DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				price = EXCLUDED.price,
				stock_qty = EXCLUDED.stock_qty,
				image_url = EXCLUDED.image_url
		`, productID, record[1], record[2], float32(price), stockQty, record[5])

		if err != nil {
			log.Println("Ошибка импорта строки:", err)
			continue
		}
		imported++
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		log.Println("Ошибка коммита транзакции:", err)
		return 0, err
	}

	log.Printf("Импорт завершен: %d записей", imported)
	return imported, nil
}
