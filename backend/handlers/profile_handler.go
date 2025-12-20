package handlers

import (
	"context"
	"log"
	"net/http"
	"online_store/database"
	"online_store/models"
	"online_store/utils"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProfileHandler(db *pgxpool.Pool) http.HandlerFunc {
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
		user, err = database.GetUser(db, ctx, user.ID) // Добавление имени в структуру user
		if err != nil {
			log.Println("Ошибка получения пользователя", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		Orders, err := database.GetOrders(db, ctx, user.ID)

		if err != nil {
			log.Println("Ошибка получения заказов", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		profile_data := models.ProfileData{
			UserName:     user.Name,
			UserEmail:    user.Email,
			RecentOrders: Orders,
		}
		err = utils.RenderTemplate(w, "profile.html", profile_data)

		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

	}

}

func ServeChangeProfilePage(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s ", r.Method)
			return
		}

		// Получаем текущего пользователя для отображения текущих данных
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
		user, err = database.GetUser(db, ctx, user.ID) // Добавление имени в структуру users

		if err != nil {
			log.Println("Ошибка получения данных пользователя")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// Передаем текущие данные пользователя в шаблон
		data := struct {
			CurrentName  string
			CurrentEmail string
			Message      string
			Success      bool
		}{
			CurrentName:  user.Name,
			CurrentEmail: user.Email,
		}

		err = utils.RenderTemplate(w, "update_profile.html", data)

		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

	}
}

func UpdateProfileHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s для ", r.Method)
			return
		}

		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value == "" {
			log.Println("Ошибка получения Cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		User, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseForm()

		if err != nil {
			log.Println("ошибка парсинга формы", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		NewName := r.FormValue("new_name")
		NewEmail := r.FormValue("new_email")
		Password := r.FormValue("password")
		ctx := r.Context()
		if NewName == "" || NewEmail == "" || Password == "" {
			showUpdateProfileError(w, db, ctx, User.ID, "Поля должны быть заполнены")
			return
		}

		if len(NewName) > 100 {
			showUpdateProfileError(w, db, ctx, User.ID, "Имя не должно быть длиннее 100 символов")
			return
		}

		if !strings.Contains(NewEmail, "@") {
			showUpdateProfileError(w, db, ctx, User.ID, "Email должен содержать @")
			return
		}

		if len(NewEmail) > 255 {
			showUpdateProfileError(w, db, ctx, User.ID, "Email не должен быть длиннее 255 символов")
			return
		}

		if len(Password) < 6 {
			showUpdateProfileError(w, db, ctx, User.ID, "Пароль должен содержать минимум 6 символов")
			return
		}

		// Проверяем, нет ли уже пользователя
		if NewEmail != User.Email {
			existingUser, err := database.GetUserByEmail(db, ctx, NewEmail)
			if err != nil {
				log.Print("Ошибка получения пользователя по Email", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if existingUser.Email != NewEmail {
				showUpdateProfileError(w, db, ctx, User.ID, "Пользователь с таким email уже существует")
				return
			}
		}

		user, err := database.GetUser(db, ctx, User.ID)

		if err != nil {
			showUpdateProfileError(w, db, ctx, user.ID, "Серверная ошибка, попробуйте еще раз")
			log.Println("Ошибка получения данных пользователя: ", err)
			return
		}

		if !CheckPassword(Password, user.Password) {
			showUpdateProfileError(w, db, ctx, user.ID, "Неверный пароль!")
			return
		}

		quary := "UPDATE users SET email = $1, name = $2 WHERE id = $3"

		_, err = db.Exec(ctx, quary, NewEmail, NewName, user.ID)

		if err != nil {
			showUpdateProfileError(w, db, ctx, user.ID, "Серверная ошибка, попробуйте еще раз")
			log.Println("Ошибка обновления данных профиля: ", err)
			return
		}

		if NewEmail != User.Email {
			newToken, err := generateToken(user.ID, NewEmail, user.Status)
			if err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "auth_token",
					Value:    newToken,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   24 * 60 * 60,
				})
			}
		}

		showUpdateProfileSuccess(w, db, ctx, user.ID, "Данные были успешно изменены!")
	}
}

func showUpdateProfileError(w http.ResponseWriter, db *pgxpool.Pool, ctx context.Context, userID int, message string) {

	user, err := database.GetUser(db, ctx, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentName  string
		CurrentEmail string
		Message      string
		Success      bool
	}{
		CurrentName:  user.Name,
		CurrentEmail: user.Email,
		Message:      message,
		Success:      false,
	}

	err = utils.RenderTemplate(w, "update_profile.html", data)

	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
}

func showUpdateProfileSuccess(w http.ResponseWriter, db *pgxpool.Pool, ctx context.Context, userID int, message string) {

	user, err := database.GetUser(db, ctx, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentName  string
		CurrentEmail string
		Message      string
		Success      bool
	}{
		CurrentName:  user.Name,
		CurrentEmail: user.Email,
		Message:      message,
		Success:      true,
	}

	err = utils.RenderTemplate(w, "update_profile.html", data)

	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		log.Println("Ошибка рендеринга шаблона update_profile.html")
		return
	}
}

func ServeChangePassword(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s ", r.Method)
			return
		}

		// Получаем текущего пользователя для отображения текущих данных
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

		data := struct {
			Message string
			Success bool
		}{
			Message: "",
			Success: false,
		}

		err = utils.RenderTemplate(w, "change_password.html", data)

		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			log.Println("Ошибка рендеринга шаблона update_profile.html")
			return
		}
	}
}

func ChangePasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s ", r.Method)
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
			log.Println("ошибка парсинга формы", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		data := struct {
			Message string
			Success bool
		}{
			Success: false,
		}

		CurrentPassword := r.FormValue("current_password")
		NewPassword := r.FormValue("new_password")
		ConfirmPassword := r.FormValue("confirm_password")
		ctx := r.Context()
		if CurrentPassword == "" || NewPassword == "" || ConfirmPassword == "" {
			data.Message = "Поля должны быть заполнены"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		if len(CurrentPassword) < 6 || len(NewPassword) < 6 || len(ConfirmPassword) < 6 {
			data.Message = "Пароль должен содержать минимум 6 символов"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		if ConfirmPassword != NewPassword {
			data.Message = "Пароли должны быть одинаковыми"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		user, err = database.GetUser(db, ctx, user.ID) // Дополнение user для получения password

		if err != nil {
			data.Message = "Серверная ошибка, попробуйте еще раз"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		if !CheckPassword(CurrentPassword, user.Password) {
			data.Message = "Неверный текущий пароль!"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		NewPasswordHash, err := HashPassword(NewPassword)

		if err != nil {
			log.Println("Ошибка хэширования пароля")
			data.Message = "Неверный текущий пароль!"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}
		quary := "UPDATE users SET password = $1 WHERE id = $2"

		_, err = db.Exec(ctx, quary, NewPasswordHash, user.ID)

		if err != nil {
			data.Message = "Серверная ошибка, попробуйте еще раз"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		data.Success = true
		data.Message = "Пароль был успешно изменен!"
		utils.ClearCookie(w, r)
		utils.RenderTemplate(w, "change_password.html", data)

	}
}

func DeleteAccountPage(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s ", r.Method)
			return
		}

		// Получаем текущего пользователя для отображения текущих данных
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

		data := struct {
			UserEmail string
			Message   string
			Success   bool
		}{
			UserEmail: user.Email,
			Message:   "",
			Success:   false,
		}
		err = utils.RenderTemplate(w, "delete_account.html", data)
		if err != nil {
			log.Println("Ошибка парсинга html файла удаления аккаунта пользователя")
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
			return
		}

	}
}

func DeleteAccountHandler(db *pgxpool.Pool) http.HandlerFunc {
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

		User, err := validateToken(cookie.Value)

		if err != nil {
			log.Println("Ошибка валидации токена", err)
			utils.ClearCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseForm()

		if err != nil {
			log.Println("ошибка парсинга формы", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		password := r.FormValue("password")
		email := r.FormValue("email")

		data := struct {
			UserEmail string
			Message   string
			Success   bool
		}{
			UserEmail: email,
			Success:   false,
		}
		ctx := r.Context()
		if password == "" || email == "" {
			data.Message = "Поля должны быть заполнены"
			utils.RenderTemplate(w, "delete_account.html", data)
			return
		}

		if len(password) < 6 {
			data.Message = "Пароль должен содержать минимум 6 символов"
			utils.RenderTemplate(w, "delete_account.html", data)
			return
		}

		if !strings.Contains(email, "@") {
			data.Message = "Email должен содержать @"
			utils.RenderTemplate(w, "delete_account.html", data)
			return
		}

		if len(email) > 255 {
			data.Message = "Email не должен быть длиннее 255 символов"
			utils.RenderTemplate(w, "delete_account.html", data)
			return
		}

		user, err := database.GetUser(db, ctx, User.ID)

		if err != nil {
			data.Message = "Серверная ошибка, попробуйте еще раз"
			utils.RenderTemplate(w, "delete_account.html", data)
			return
		}

		if !CheckPassword(password, user.Password) { //Пароль введеный пользователем и его пароль из БД
			data.Message = "Неверный текущий пароль!"
			utils.RenderTemplate(w, "delete_account.html", data)
			return
		}

		//Каскадное удаление с users, cart, cart_items, orders, order_items
		quary := "Delete FROM users WHERE id = $1"

		_, err = db.Exec(ctx, quary, user.ID)

		if err != nil {
			data.Message = "Серверная ошибка, попробуйте еще раз"
			utils.RenderTemplate(w, "delete_account.html", data)
			log.Println("Ошибка обновления данных профиля: ", err)
			return
		}

		utils.ClearCookie(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
