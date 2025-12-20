package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"online_store/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
GetFilteredProducts выполняет сложный запрос к базе данных для получения товаров
с применением множественных фильтров и сортировки. Метод динамически строит SQL-запрос
на основе переданных параметров.
*/
func GetFilteredProducts(db *pgxpool.Pool, ctx context.Context, params *models.FilterParams) ([]models.ProductForFilter, error) {
	query := `
        SELECT 
        products.product_id, 
        products.name, 
        products.price, 
        products.stock_qty, 
        products_categories.category_id,
        products_producers.producer_id,
        categories.name as category_name,
        producers.name as producer_name,
        products.image_url -- Добавьте эту строку
    FROM products
    JOIN products_categories ON products.product_id = products_categories.product_id
    JOIN products_producers ON products.product_id = products_producers.product_id
    JOIN categories ON products_categories.category_id = categories.category_id
    JOIN producers ON products_producers.producer_id = producers.producer_id
    WHERE 1=1 
    ` // "WHERE 1=1" - точка начала для удобного добавления условий через AND

	args := []interface{}{}
	paramCount := 0

	// Фильтр по категории
	if len(params.Categories) > 0 {
		paramCount++
		placeholders := make([]string, len(params.Categories))
		for i, catID := range params.Categories {
			args = append(args, catID)
			placeholders[i] = fmt.Sprintf("$%d", paramCount+i)
		}
		query += fmt.Sprintf(" AND products_categories.category_id IN (%s)", strings.Join(placeholders, ","))
		paramCount += len(params.Categories) - 1
	}

	// Фильтр по производителю
	if len(params.Producers) > 0 {
		paramCount++
		placeholders := make([]string, len(params.Producers))
		for i, prodID := range params.Producers {
			args = append(args, prodID)
			placeholders[i] = fmt.Sprintf("$%d", paramCount+i)
		}
		query += fmt.Sprintf(" AND products_producers.producer_id IN (%s)", strings.Join(placeholders, ","))
		paramCount += len(params.Producers) - 1
	}

	// Фильтр по максимальной цене
	if params.MaxPrice > 0 {
		paramCount++
		query += fmt.Sprintf(" AND products.price <= $%d", paramCount)
		args = append(args, params.MaxPrice)
	}

	// Фильтр по минимальной цене
	if params.MinPrice > 0 {
		paramCount++
		query += fmt.Sprintf(" AND products.price >= $%d", paramCount)
		args = append(args, params.MinPrice)
	}

	if params.SearchQuery != "" {
		paramCount++
		query += fmt.Sprintf(" AND products.name ILIKE $%d", paramCount)
		args = append(args, "%"+params.SearchQuery+"%")
	}

	switch params.SortByPrice {
	case "ASC":
		query += " ORDER BY products.price ASC"
	case "DESC":
		query += " ORDER BY products.price DESC"
	default:
		// Сортировка по умолчанию (опционально)
		query += " ORDER BY products.product_id ASC"
	}

	// Выполняем запрос
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Println("Ошибка получения продуктов:", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.ProductForFilter
	for rows.Next() {
		var product models.ProductForFilter
		// Сканируем ВСЕ 6 полей!
		err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Price,
			&product.StockQuantity,
			&product.CategoryID,
			&product.ProducerID,
			&product.CategoryName,
			&product.ProducerName,
			&product.Image, // Добавьте это
		)
		if err != nil {
			log.Println("Ошибка сканирования продукта:", err)
			continue // Не прерываем весь запрос из-за одной строки
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetCategories получает все категории из БД
func GetCategories(db *pgxpool.Pool, ctx context.Context) ([]models.Category, error) {
	query := `SELECT category_id, name FROM categories ORDER BY name`
	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Println("Ошибка получения категорий:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Println("Ошибка сканирования категории:", err)
			continue
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetProducers(db *pgxpool.Pool, ctx context.Context) ([]models.Producer, error) {
	query := `SELECT producer_id, name, country FROM producers ORDER BY name`
	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Println("Ошибка получения производителей:", err)
		return nil, err
	}
	defer rows.Close()

	var producers []models.Producer
	for rows.Next() {
		var producer models.Producer
		err := rows.Scan(&producer.ID, &producer.Name, &producer.Country)
		if err != nil {
			log.Println("Ошибка сканирования производителя:", err)
			continue
		}
		producers = append(producers, producer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return producers, nil
}
