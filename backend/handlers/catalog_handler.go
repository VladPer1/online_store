package handlers

import (
	"log"
	"net/http"
	"online_store/database"
	"online_store/models"
	"online_store/utils"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func serveCatalog(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			log.Printf("Неподдерживаемый метод %s для /catalog", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()

		// Получаем категории и производителей из БД
		categoriesFromDB, err := database.GetCategories(db, ctx)
		if err != nil {
			log.Printf("Ошибка получения категорий: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		producersFromDB, err := database.GetProducers(db, ctx)
		if err != nil {
			log.Printf("Ошибка получения производителей: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Парсим параметры фильтрации
		searchQuery := r.URL.Query().Get("search")
		minPriceStr := r.URL.Query().Get("price_min")
		maxPriceStr := r.URL.Query().Get("price_max")
		categoryParams := r.URL.Query()["category"]
		producerParams := r.URL.Query()["producer"]
		sortBy := r.URL.Query().Get("sortBy")

		// Конвертируем категории в int
		var categoryIDs []int
		for _, catStr := range categoryParams {
			if catID, err := strconv.Atoi(catStr); err == nil {
				categoryIDs = append(categoryIDs, catID)
			}
		}

		// Конвертируем производителей в int
		var producerIDs []int
		for _, prodStr := range producerParams {
			if prodID, err := strconv.Atoi(prodStr); err == nil {
				producerIDs = append(producerIDs, prodID)
			}
		}

		// Подготавливаем параметры для фильтрации
		filterParams := &models.FilterParams{
			SearchQuery: searchQuery,
			Categories:  categoryIDs, // Используем правильное поле
			Producers:   producerIDs, // Используем правильное поле
		}

		// Конвертируем цены
		if minPriceStr != "" {
			if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
				filterParams.MinPrice = minPrice
			}
		}
		if maxPriceStr != "" {
			if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
				filterParams.MaxPrice = maxPrice
			}
		}

		// Сортировка
		if sortBy == "ASC" || sortBy == "DESC" {
			filterParams.SortByPrice = sortBy
		}
		// Получаем отфильтрованные продукты
		products, err := database.GetFilteredProducts(db, r.Context(), filterParams)
		if err != nil {
			log.Printf("Ошибка получения продуктов: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Подготавливаем данные для шаблона
		data := models.CatalogData{
			Title:              "Каталог товаров - GainWave",
			Products:           products,
			Categories:         categoriesFromDB,
			Producers:          producersFromDB,
			SearchQuery:        searchQuery,
			SelectedCategories: models.MakeStringSet(categoryParams),
			SelectedBrands:     models.MakeStringSet(producerParams),
			PriceMin:           minPriceStr,
			PriceMax:           maxPriceStr,
			SortBy:             sortBy,
		}

		// Логируем для отладки
		log.Printf("Категории: %d, Производители: %d, Товары: %d",
			len(categoriesFromDB), len(producersFromDB), len(products))
		log.Printf("Параметры фильтрации: %+v", filterParams)
		// Загружаем и выполняем шаблон

		err = utils.RenderTemplate(w, "catalog.html", data)
		if err != nil {
			log.Printf("Ошибка загрузки шаблона catalog: %v", err)
			http.Error(w, "Catalog template not found", http.StatusInternalServerError)
			return
		}

	}
}
