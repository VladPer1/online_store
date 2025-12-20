package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"online_store/database"
	"online_store/models"
	"online_store/utils"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

// AdminPanel - главная страница админ-панели.
func AdminPanelHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Проверяем права администратора
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil {
			log.Println("Ошибка проверки статуса пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			log.Println("Пользователь не является администратором")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Подготавливаем данные для шаблона
		data := struct {
			Title          string
			CurrentSection string
			Action         string
			Message        string
			Success        bool
			Products       []models.ProductForAdmin
			Categories     []models.Category
			Producers      []models.Producer
			Product        models.ProductForAdmin
		}{
			Title:          "Админ-панель - GainWave",
			CurrentSection: "products", // по умолчанию
		}

		// Получаем параметры из URL
		if section := r.URL.Query().Get("section"); section != "" {
			data.CurrentSection = section
		}
		data.Action = r.URL.Query().Get("action")

		// Загружаем данные в зависимости от раздела
		switch data.CurrentSection {
		case "products":
			// Загружаем категории и производители для форм
			categories, err := database.GetCategories(db, ctx)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Ошибка загрузки категорий:", err)
				return
			} else {
				data.Categories = categories
			}

			producers, err := database.GetProducers(db, ctx)
			if err != nil {
				log.Println("Ошибка загрузки производителей:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				data.Producers = producers
			}

			switch data.Action {
			case "edit":
				// Режим редактирования - загружаем данные товара
				if productIDStr := r.URL.Query().Get("id"); productIDStr != "" {
					productID, err := strconv.Atoi(productIDStr)
					if err == nil {
						product, err := database.GetProductForAdmin(db, ctx, productID)
						if err == nil {
							data.Product = product
						} else {
							log.Println("Ошибка загрузки товара:", err)
							http.Error(w, "Internal Server Error", http.StatusInternalServerError)
							return
						}
					}
				}
			case "add":
				// Режим добавления - пустой продукт
				data.Product = models.ProductForAdmin{}
			default:
				// Режим просмотра списка - загружаем товары
				products, err := database.GetAllProductsForAdmin(db, ctx)
				if err != nil {
					log.Println("Ошибка загрузки товаров:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				} else {
					data.Products = products
				}
			}

		case "ban":
			// Раздел бана пользователей - не требует дополнительных данных
			// Просто показываем формы для бана/разбана

		}

		// Обработка сообщений
		if msg := r.URL.Query().Get("message"); msg != "" {
			data.Message = msg
			data.Success = r.URL.Query().Get("success") == "true"
		}

		// Рендерим шаблон
		err = utils.RenderTemplate(w, "admin.html", data)

		if err != nil {
			log.Println("Ошибка рендеринга шаблона", err)
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
	}
}

// AdminSaveProduct - сохранение товара (добавление и редактирование).
func AdminSaveProductHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Проверяем права администратора
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil {
			log.Println("Ошибка проверки статуса пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return

		}

		if !isAdmin {
			log.Println("Пользователь не является администратором")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Парсим форму
		productID, _ := strconv.Atoi(r.FormValue("product_id"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 32)
		stock, _ := strconv.Atoi(r.FormValue("stock_quantity"))
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
		producerID, _ := strconv.Atoi(r.FormValue("producer_id"))

		product := models.ProductForAdmin{
			ProductID:     productID,
			Name:          r.FormValue("name"),
			Description:   r.FormValue("description"),
			Price:         float32(price),
			StockQuantity: stock,
			ImageURL:      r.FormValue("image_url"),
			CategoryID:    categoryID,
			ProducerID:    producerID,
		}

		var message string
		var success bool

		// Сохраняем в базу
		if productID == 0 {
			// В HTML шаблоне при редактировании productID передается как 0, а при создании товара как 123
			// Добавление нового товара
			err = database.CreateProductForAdmin(db, ctx, product)
			message = "Товар успешно добавлен"
		} else {
			// Редактирование существующего товара
			err = database.UpdateProductForAdmin(db, ctx, product)
			message = "Товар успешно обновлен"
		}

		if err != nil {
			log.Println(err)
			message = "Ошибка сохранения: " + err.Error()
			success = false
		} else {
			success = true
		}

		// Редирект с сообщением
		redirectURL := "/admin?section=products&message=" + message + "&success=" + strconv.FormatBool(success)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

// AdminDeleteProduct - удаление товара.
func AdminDeleteProductHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Проверяем права администратора
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil {
			log.Println("Ошибка проверки статуса пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			log.Println("Пользователь не является администратором")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		productID, _ := strconv.Atoi(r.FormValue("product_id"))
		err = database.DeleteProductForAdmin(db, ctx, productID)

		var message string
		var success bool

		if err != nil {
			message = "Ошибка удаления: " + err.Error()
			success = false
		} else {
			message = "Товар успешно удален"
			success = true
		}

		redirectURL := "/admin?section=products&message=" + message + "&success=" + strconv.FormatBool(success)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

// AdminBanUser - блокировка пользователя
func AdminBanUserHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Проверяем права администратора
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil {
			log.Println("Ошибка проверки статуса пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return

		}

		if !isAdmin {
			log.Println("Пользователь не является администратором")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		err = database.BanUser(db, r.Context(), userID)

		var message string
		var success bool

		if err != nil {
			message = "Ошибка блокировки: " + err.Error()
			success = false
		} else {
			message = "Пользователь заблокирован"
			success = true
		}

		redirectURL := "/admin?section=users&message=" + message + "&success=" + strconv.FormatBool(success)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

// AdminUnbanUser - разблокировка пользователя
func AdminUnbanUserHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Проверяем права администратора
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil {
			log.Println("Ошибка проверки статуса пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			log.Println("Пользователь не является администратором")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		err = database.UnbanUser(db, ctx, userID)

		var message string
		var success bool

		if err != nil {
			message = "Ошибка разблокировки: " + err.Error()
			success = false
		} else {
			message = "Пользователь разблокирован"
			success = true
		}

		redirectURL := "/admin?section=users&message=" + message + "&success=" + strconv.FormatBool(success)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func GenerateReportHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Проверяем права администратора
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil {
			log.Println("Ошибка проверки статуса пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			log.Println("Пользователь не является администратором")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Создаем Excel файл
		f := excelize.NewFile()
		defer f.Close()

		// Создаем лист для данных
		sheetName := "Отчет"
		index, _ := f.NewSheet(sheetName)
		f.SetActiveSheet(index)

		// Заголовки столбцов
		headers := []string{"ID", "Название товара", "Количество на складе"}
		for col, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Устанавливаем ширину столбцов
		widths := []float64{5, 45, 20}
		for col, width := range widths {
			colName, _ := excelize.ColumnNumberToName(col + 1)
			f.SetColWidth(sheetName, colName, colName, width)
		}

		// Получаем данные из PostgreSQL
		query := `
			SELECT product_id, name, stock_qty
			FROM products 
			ORDER BY stock_qty
		`
		rows, err := db.Query(ctx, query)
		if err != nil {
			log.Println("Ошибка работы с БД: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Создаем стили заранее
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"2c5faa"}, Pattern: 1},
		})

		highlightStyle, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#FF918B"},
				Pattern: 1,
			},
		})

		// Заполняем данными и сразу применяем стили
		rowNum := 2
		for rows.Next() {
			var id, stock_qty int
			var name string

			err := rows.Scan(&id, &name, &stock_qty)
			if err != nil {
				log.Println("Ошибка сканирования: ", err)
				continue
			}

			// Записываем данные
			data := []interface{}{id, name, stock_qty}
			for col, value := range data {
				cell, _ := excelize.CoordinatesToCellName(col+1, rowNum)
				f.SetCellValue(sheetName, cell, value)

				// Если это столбец количества и значение < 10, применяем красный стиль
				if col == 2 && stock_qty < 10 {
					f.SetCellStyle(sheetName, cell, cell, highlightStyle)
				}
			}
			rowNum++
		}

		// Применяем стиль заголовков
		f.SetCellStyle(sheetName, "A1", "C1", headerStyle)

		// Генерируем имя файла
		filename := "report_" + time.Now().Format("2006-01-02_15-04-05") + ".xlsx"

		// Устанавливаем заголовки для скачивания
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Pragma", "no-cache")

		// Сохраняем файл в response
		if err := f.Write(w); err != nil {
			log.Println("Ошибка записи Excel файла: ", err)
			http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
			return
		}

		log.Printf("Отчет успешно сгенерирован: %s, строк: %d", filename, rowNum-2)
	}
}

func ExportProductsCSVHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil || !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Создаем CSV файл
		filename := fmt.Sprintf("backup_%s.csv", time.Now().Format("2006-01-02_15-04-05"))
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Заголовки CSV
		headers := []string{"product_id", "name", "description", "price", "stock_qty", "image_url"}
		if err := writer.Write(headers); err != nil {
			log.Println("Ошибка записи заголовков CSV:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Экспортируем данные
		err = database.ExportProductsToCSV(db, ctx, writer)
		if err != nil {
			log.Println("Ошибка экспорта данных:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("CSV экспорт завершен: %s", filename)
	}
}

func ImportProductsCSVHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Redirect(w, r, "/admin?message=Method+not+allowed&success=false", http.StatusSeeOther)
			return
		}

		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/admin?message=Unauthorized&success=false", http.StatusSeeOther)
			return
		}

		user, err := validateToken(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/admin?message=Unauthorized&success=false", http.StatusSeeOther)
			return
		}

		ctx := r.Context()
		isAdmin, err := database.CheckAdminStatus(db, ctx, user.ID)
		if err != nil || !isAdmin {
			http.Redirect(w, r, "/admin?message=Forbidden&success=false", http.StatusSeeOther)
			return
		}

		// Обработка загружаемого файла
		file, header, err := r.FormFile("csv_file")
		if err != nil {
			http.Redirect(w, r, "/admin?message=Ошибка+загрузки+файла&success=false", http.StatusSeeOther)
			return
		}
		defer file.Close()

		// Проверяем тип файла
		if header.Header.Get("Content-Type") != "text/csv" {
			http.Redirect(w, r, "/admin?message=Неверный+формат+файла&success=false", http.StatusSeeOther)
			return
		}

		// Импортируем данные
		imported, err := database.ImportProductsCSV(db, ctx, file)
		if err != nil {
			log.Println("Ошибка импорта данных:", err)
			http.Redirect(w, r, "/admin?message=Ошибка+импорта+данных&success=false", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/admin?message=Успешно+импортировано+%d+записей&success=true", imported), http.StatusSeeOther)
	}
}
