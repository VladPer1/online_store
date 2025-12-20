package handlers

import (
	"context"
	//"html/template"
	"log"
	"net/http"
	"online_store/database"
	"online_store/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProcessPaymentHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s", r.Method)
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
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		orderID, err := database.CreateOrder(db, ctx, user.ID)

		if err := ctx.Err(); err != nil {
			log.Println("Контекст отменен", err)
			http.Redirect(w, r, "/error-payment", http.StatusSeeOther)
			return
		}

		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				log.Printf("Платеж отменен по таймауту для %s: %v", user.Email, ctxErr)
				http.Error(w, "Payment timeout", http.StatusRequestTimeout)
				return
			} else {
				log.Printf("Ошибка при изменении количсетва товаров в бд %s: %v", user.Email, err)
				http.Redirect(w, r, "/error-payment", http.StatusSeeOther)
				return
			}
		} else {
			http.Redirect(w, r, "/success-payment", http.StatusSeeOther)
			go SendReceiptOnEmailAsync(user.Email, orderID)

		}

	}
}
func SuccessPaymentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := utils.RenderTemplate(w, "success_payment.html", nil)
		if err != nil {
			log.Println("Ошибка рендеринга шаблона", err)
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
			return
		}

	}
}

func ErrorPaymentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := utils.RenderTemplate(w, "payment_error.html", nil)
		if err != nil {
			log.Println("Ошибка рендеринга шаблона", err)
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
		}

	}
}
