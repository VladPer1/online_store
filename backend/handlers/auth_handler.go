package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/smtp"
	"online_store/database"
	"online_store/models"
	"online_store/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Функция генерации JWT токена
func generateToken(userID int, email string, status string) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		jwtSecret = "temporary-dev-secret-change-in-production"
		log.Println("ВНИМАНИЕ: Используется дефолтный JWT_SECRET!")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["email"] = email
	claims["status"] = status
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["iat"] = time.Now().Unix()

	return token.SignedString(jwtSecret)
}

// validateToken проверяет JWT токен и возвращает данные пользователя

func validateToken(tokenString string) (*models.User, error) {

	jwtSecret := []byte("GainWave-Sports-Supplement-Store-2024-Secure-JWT-Key!")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &models.User{
			ID:     int(claims["user_id"].(float64)), // JWT числа становятся float64
			Email:  claims["email"].(string),
			Status: claims["status"].(string),
		}
		return user, nil
	}

	return nil, fmt.Errorf("невалидный токен")
}

func HashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost) // 10
	return string(bytes), err
}

func CheckPassword(pwd, hash_pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash_pwd), []byte(pwd))
	return err == nil
}

// loginFormHandler обрабатывает отправку формы входа
func loginFormHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := models.TemplateData{
			Success:  false,
			FormType: "login",
		}
		// Парсим форму
		err := r.ParseForm()
		if err != nil {
			data.Message = "Ошибка обработки формы"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		data.Email = r.FormValue("email")
		password := r.FormValue("password")

		if data.Email == "" || password == "" {
			data.Message = "Email и пароль обязательны"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if len(password) < 6 {
			data.Message = "Пароль должен содержать минимум 6 символов"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if !strings.Contains(data.Email, "@") {
			data.Message = "Email должен содержать @"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if len(data.Email) > 255 {
			data.Message = "Email не должен быть длиннее 255 символов"
			utils.RenderTemplate(w, "login.html", data)
			return
		}
		ctx := r.Context()
		// Ищем пользователя
		user, err := database.GetUserByEmail(db, ctx, data.Email)
		if err != nil {
			data.Message = "Неверный email или пароль"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		// Проверяем пароль
		if !CheckPassword(password, user.Password) {
			data.Message = "Неверный email или пароль"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		// Генерируем токен
		token, err := generateToken(user.ID, user.Email, user.Status)
		if err != nil {
			data.Message = "Ошибка входа"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		// Сохраняем токен в cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   24 * 60 * 60, // 24 часа
		})

		switch user.Status {
		case "user":
			http.Redirect(w, r, "/catalog", http.StatusSeeOther)
		case "admin":
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		case "banned":
			data.Message = "Вы забанены по решению администратора, вы не можете зайти на сайт"
			utils.RenderTemplate(w, "login.html", data)
			return
		default:
			// Обработка неизвестного статуса
			data.Message = "Неизвестный статус пользователя"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

	}
}

// registerFormHandler обрабатывает отправку формы регистрации
func registerFormHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := models.TemplateData{
			Success:  false,
			FormType: "register",
		}
		// Парсим форму
		err := r.ParseForm()
		if err != nil {
			data.Message = "Ошибка обработки формы"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		data.Name = r.FormValue("name")
		data.Email = r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Валидация
		if data.Name == "" || data.Email == "" || password == "" || confirmPassword == "" {
			data.Message = "Все поля обязательны для заполнения"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if password != confirmPassword {
			data.Message = "Пароли не совпадают"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if len(password) < 6 {
			data.Message = "Пароль должен содержать минимум 6 символов"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if !strings.Contains(data.Email, "@") {
			data.Message = "Email должен содержать @"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if len(data.Name) > 100 {
			data.Message = "Имя не должно быть длиннее 100 символов"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		if len(data.Email) > 255 {
			data.Message = "Email не должен быть длиннее 255 символов"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		// Проверяем, нет ли уже пользователя
		existingUser, err := database.GetUserByEmail(db, r.Context(), data.Email)
		if err == nil && existingUser.ID != 0 {
			data.Message = "Пользователь с таким email уже существует"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		// Хешируем пароль
		hashedPassword, err := HashPassword(password)
		if err != nil {
			data.Message = "Ошибка обработки пароля"
			utils.RenderTemplate(w, "login.html", data)
			return
		}

		VerificationCode, err := SendCodeOnEmail(data.Email)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Создаем пользователя
		user := &models.PlaceholderUser{
			Email:            data.Email,
			Password:         hashedPassword,
			Name:             data.Name,
			VerificationCode: VerificationCode,
		}

		ctx := r.Context()
		_ = database.DeletePlaceholderUser(db, ctx, user.Email)
		err = database.CreatePlaceholderUser(db, ctx, user)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/verify-code-page?email="+data.Email, http.StatusSeeOther)

	}
}

// handlers/auth.go
func ServeVerifyCodePage(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s", r.Method)
			return
		}

		// Получаем email из URL: /verify-code-page?email=test@example.com
		email := r.URL.Query().Get("email")

		if email == "" {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			log.Println("Пустой email")
			return
		}

		// Рендерим страницу verify-code.html
		data := models.TemplateData{Email: email}
		utils.RenderTemplate(w, "verify-code.html", data)
	}
}

func verifyCodeHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем данные из формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email") // из скрытого поля
		enteredCode := r.FormValue("code")

		data := models.TemplateData{
			Email: email,
		}

		if email == "" || enteredCode == "" {
			data.Message = "Все поля обязательны для заполнения"
			utils.RenderTemplate(w, "verify-code.html", data)
			return
		}
		enteredCodeInt, err := strconv.Atoi(enteredCode)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Ошибка конвертации atoi:", err)
			return
		}

		ctx := r.Context()
		pendingUser, err := database.GetPlaceholderUser(db, ctx, email)
		if err != nil {
			data.Message = "Ошибка получения данных пользователя"
			utils.RenderTemplate(w, "verify-code.html", data)
			return
		}

		if enteredCodeInt != pendingUser.VerificationCode {
			data.Message = "Неверный код"
			utils.RenderTemplate(w, "verify-code.html", data)
			return
		}

		user := models.User{
			Email:    pendingUser.Email,
			Password: pendingUser.Password,
			Name:     pendingUser.Name,
		}

		err = database.CreateUser(db, ctx, &user)
		if err != nil {
			data.Message = "Ошибка создания пользователя"
			utils.RenderTemplate(w, "verify-code.html", data)
			return
		}

		// Генерируем токен
		token, err := generateToken(user.ID, user.Email, user.Status)
		if err != nil {
			data.Message = "Ошибка регистрации"
			utils.RenderTemplate(w, "verify-code.html", data)
			return
		}

		// Сохраняем токен в cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   24 * 60 * 60,
		})

		err = database.DeletePlaceholderUser(db, ctx, email)

		if err != nil {
			log.Println("Ошибка удаления данных пользователя из таблицы pending_registrations")
		}
		// Перенаправляем в личный кабинет
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}

func ServeForgotPasswordPage(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s", r.Method)
			return
		}

		data := models.TemplateData{}
		utils.RenderTemplate(w, "forgot_password.html", data)
	}
}

func ForgotPasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем данные из формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")

		data := models.TemplateData{
			Email: email,
		}

		if email == "" {
			data.Message = "Поле email обязательно для заполнения"
			utils.RenderTemplate(w, "forgot-password.html", data)
			return
		}
		ctx := r.Context()
		_, err = database.GetUserByEmail(db, ctx, email)

		if err != nil {
			data.Message = "Пользователя с таким email не существует"
			utils.RenderTemplate(w, "forgot-password.html", data)
			return
		}

		VerificationCode, err := SendCodeOnEmail(data.Email)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		//Заносим данные чтобы сравнить потом код подтверждения
		user := &models.PlaceholderUser{
			Email:            data.Email,
			VerificationCode: VerificationCode,
		}
		_ = database.DeleteUserFromFPD(db, ctx, data.Email)
		err = database.CreateUserInFPD(db, ctx, user)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Перенаправляем на страницу ввода кода
		http.Redirect(w, r, "/forgot-password-verify-page?email="+email, http.StatusSeeOther)
	}
}

func ServeForgotPasswordVerifyPage(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s", r.Method)
			return
		}

		email := r.URL.Query().Get("email")
		if email == "" {
			http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
			return
		}

		data := models.TemplateData{Email: email}
		utils.RenderTemplate(w, "forgot_password_verify.html", data)
	}
}

func ForgotPasswordVerifyHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем данные из формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email") //скрытое поле
		enteredCode := r.FormValue("code")

		data := models.TemplateData{
			Email: email,
		}

		if enteredCode == "" {
			data.Message = "Все поля обязательны для заполнения"
			utils.RenderTemplate(w, "forgot_password_verify.html", data)
			return
		}
		enteredCodeInt, err := strconv.Atoi(enteredCode)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Ошибка конвертации atoi:", err)
			return
		}

		ctx := r.Context()
		pendingUser, err := database.GetUserFromFPD(db, ctx, email)
		if err != nil {
			data.Message = "Ошибка получения данных пользователя"
			utils.RenderTemplate(w, "forgot_password_verify.html", data)
			return
		}

		if enteredCodeInt != pendingUser.VerificationCode {
			data.Message = "Неверный код"
			utils.RenderTemplate(w, "forgot_password_verify.html", data)
			return
		}

		http.Redirect(w, r, "/forgot-password-update-password-page?email="+email, http.StatusSeeOther)
	}
}

func ServeForgotPasswordUpdatePasswordPage(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Неподдерживаемый метод %s", r.Method)
			return
		}

		email := r.URL.Query().Get("email")
		if email == "" {
			http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
			return
		}

		data := models.TemplateData{Email: email}
		utils.RenderTemplate(w, "forgot_password_update_password.html", data)
	}
}

func ForgotPasswordUpdatePasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем данные из формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email") // скрытое поле
		ConfirmPassword := r.FormValue("confirm_password")
		NewPassword := r.FormValue("new_password")

		data := models.TemplateData{
			Email: email,
		}

		if NewPassword == "" || ConfirmPassword == "" {
			data.Message = "Поля должны быть заполнены"
			utils.RenderTemplate(w, "forgot_password_update_password.html", data)
			return
		}

		if len(NewPassword) < 6 || len(ConfirmPassword) < 6 {
			data.Message = "Пароль должен содержать минимум 6 символов"
			utils.RenderTemplate(w, "forgot_password_update_password.html", data)
			return
		}

		if ConfirmPassword != NewPassword {
			data.Message = "Пароли должны быть одинаковыми"
			utils.RenderTemplate(w, "forgot_password_update_password.html", data)
			return
		}

		ctx := r.Context()

		NewPasswordHash, err := HashPassword(NewPassword)

		if err != nil {
			log.Println("Ошибка хэширования пароля")
			data.Message = "Неверный текущий пароль!"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}
		quary := "UPDATE users SET password = $1 WHERE email = $2"

		_, err = db.Exec(ctx, quary, NewPasswordHash, email)

		if err != nil {
			data.Message = "Серверная ошибка, попробуйте еще раз"
			utils.RenderTemplate(w, "change_password.html", data)
			return
		}

		err = database.DeleteUserFromFPD(db, ctx, email)

		if err != nil {
			log.Println("Ошибка удаления данных из таблицы pending_registrations")
		}

		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	}
}

func SendReceiptOnEmailAsync(userEmail string, orderID int) {

	smtpPort := "587"
	smtpHost := "smtp.gmail.com"
	from := "gainwavegainwave@gmail.com"
	password := os.Getenv("SMTP_PASSWORD")

	message := []byte("Subject: Чек\r\n" +
		"\r\n" + // Пустая строка между заголовками и телом
		"Спасибо за покупку на Gain Wave:\r\n" +
		"Номер вашего заказа: " + strconv.Itoa(orderID) + "\r\n" +
		"Спасибо, что выбрали нас!.\r\n" +
		"С уважением, Gain Wave.\r\n")

	if password == "" {
		log.Printf("SMTP_PASSWORD не установлен для email: %s", userEmail)
		return
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	userEmailArr := []string{userEmail}
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, userEmailArr, message)
	if err != nil {
		log.Printf("Ошибка отправки email на %s: %v", userEmail, err)
	} else {
		log.Printf("Email успешно отправлен на %s", userEmail)
	}

}

func SendCodeOnEmail(userEmail string) (int, error) {

	VerificationСode, err := CreateRandomNum()
	if err != nil {
		return 0, err
	}
	message := []byte("Subject: Код подтверждения\r\n" +
		"\r\n" + // Пустая строка между заголовками и телом
		"Здравствуйте!\r\n" +
		"Ваш код подтверждения: " + VerificationСode + "\r\n" +
		"Используйте его для входа в аккаунт.\r\n" +
		"\r\n" +
		"Если это были не вы, проигнорируйте это сообщение.\r\n" +
		"С уважением, Gain Wave.\r\n")

	go SendCodeOnEmailAsync(userEmail, message)

	VerificationСodeInt, err := strconv.Atoi(VerificationСode)

	if err != nil {
		log.Println("Ошибка конвертации atoi:", err)
		return 0, err
	}

	return VerificationСodeInt, nil
}

func SendCodeOnEmailAsync(email string, message []byte) {
	smtpPort := "587"
	smtpHost := "smtp.gmail.com"
	from := "gainwavegainwave@gmail.com"
	password := os.Getenv("SMTP_PASSWORD")

	if password == "" {
		log.Printf("SMTP_PASSWORD не установлен для email: %s", email)
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	userEmailArr := []string{email}
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, userEmailArr, message)
	if err != nil {
		log.Printf("Ошибка отправки email на %s: %v", email, err)
	} else {
		log.Printf("Email успешно отправлен на %s", email)
	}
}

func CreateRandomNum() (string, error) {
	// Генерируем число от 0 до 89999
	num, err := rand.Int(rand.Reader, big.NewInt(90000))
	if err != nil {
		log.Println("Ошибка генерации случайного числа:", err)
		return "", err
	}

	// Добавляем 10000 чтобы получить диапазон 10000-99999
	randomNum := num.Int64() + 10000
	return fmt.Sprintf("%05d", randomNum), nil
}

// Обработчик выхода
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		utils.ClearCookie(w, r)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Метод Handler
func serveLoginPage(w http.ResponseWriter, r *http.Request) {

	data := models.TemplateData{
		FormType: "login",
	}

	err := utils.RenderTemplate(w, "login.html", data)

	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		log.Println("Ошибка рендеринга login.html")
		return
	}

}
