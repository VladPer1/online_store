package handlers

import (
	"log"
	"net/http"
	"strconv"

	"online_store/database"
	"online_store/models"
	"online_store/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CartHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			log.Printf("Неподдерживаемый метод %s для /cart", r.Method)
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

		cart_id, err := database.GetCartIDUser(db, ctx, user.ID)
		if err != nil {
			cart_id, err = database.CreateCartForUser(db, ctx, user.ID)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Println("Ошибка создания корзины пользователя", err)
				return
			}
		}

		UserCart, err := database.GetUserCart(db, ctx, cart_id)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Ошибка получения корзины пользователя", err)
			return
		}

		var SumCartProducts float32
		var TotalItems = 0

		for i := range UserCart {
			UserCart[i].TotalPrice = UserCart[i].Price * float32(UserCart[i].Quantity)
			SumCartProducts += UserCart[i].TotalPrice
			TotalItems += UserCart[i].Quantity
		}

		// Подготавливаем данные для шаблона
		data := struct {
			Title      string
			UserCart   []models.CartItem
			Subtotal   float32
			TotalItems int
		}{
			Title:      "Корзина - GainWave",
			UserCart:   UserCart,
			Subtotal:   SumCartProducts,
			TotalItems: TotalItems,
		}
		err = utils.RenderTemplate(w, "cart.html", data)

		if err != nil {
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
			log.Println("Ошибка рендеринга cart.html", err)
			return
		}
	}
}

func AddProductToCartHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s ", r.Method)
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

		err = r.ParseForm()
		if err != nil {
			log.Println("Ошибка парсинга формы", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		productIDStr := r.FormValue("product_id")
		quantityStr := r.FormValue("quantity")

		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			log.Println("Ошибка конвертации atoi", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(quantityStr)
		if err != nil || quantity < 1 {
			log.Println("Ошибка конвертации atoi", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		cartID, err := database.GetCartIDUser(db, ctx, user.ID)
		if err != nil {
			// Если корзины нет - создаем
			cartID, err = database.CreateCartForUser(db, ctx, user.ID)
			if err != nil {
				log.Println("Ошибка создания корзины", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		err = database.AddProductToCart(db, ctx, cartID, productID, quantity)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Ошибка, товар не добавлен в корзину: ", err)
			return
		} else {
			log.Println("Товар добавлен в корзину: ")
			http.Redirect(w, r, "/catalog", http.StatusSeeOther)
		}

	}

}

func DeleteProductFromCartHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s ", r.Method)
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

		err = r.ParseForm()
		if err != nil {
			log.Print("Ошибка парсинга формы", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		product_id_str := r.FormValue("product_id")

		product_id, err := strconv.Atoi(product_id_str)
		if err != nil {
			log.Println("Ошибка конвертации product_id atoi", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		cartID, err := database.GetCartIDUser(db, ctx, user.ID)
		if err != nil {
			log.Print("Ошибка получения ID пользователя", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = database.DeleteProductFromCart(db, ctx, cartID, product_id)
		if err != nil {
			log.Print("Ошибка удаления элемента корзины", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		log.Println("Товар удален из корзины")
	}
}
