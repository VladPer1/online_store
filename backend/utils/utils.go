package utils

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"online_store/models"
)

var AppConfig models.Config

func ClearCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// InitConfig инициализирует конфигурацию
func InitConfig() {
	// Путь к шаблонам
	AppConfig.TemplatePath = os.Getenv("TEMPLATE_PATH")
	if AppConfig.TemplatePath == "" {
		AppConfig.TemplatePath = "templates" // значение по умолчанию
	}

	// Путь к статическим файлам
	AppConfig.StaticPath = os.Getenv("STATIC_PATH")
	if AppConfig.StaticPath == "" {
		AppConfig.StaticPath = "static" // значение по умолчанию
	}

	log.Printf("Конфигурация: templates=%s, static=%s",
		AppConfig.TemplatePath, AppConfig.StaticPath)
}

func PutFiles(mux *http.ServeMux) {
	// Статические файлы
	fileSystem := http.Dir(AppConfig.StaticPath)
	fileServer := http.FileServer(fileSystem)
	staticHandler := http.StripPrefix("/static/", fileServer)

	mux.Handle("/static/", staticHandler)

}

// Вспомогательная функция для рендеринга шаблонов
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	// Путь к папке с шаблонами

	path := filepath.Join(AppConfig.TemplatePath, tmpl)

	// Логируем путь для отладки
	log.Printf("Загрузка шаблона: %s", path)

	// Парсим шаблон
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Ошибка парсинга шаблона %s: %v", tmpl, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return err
	}

	// Выполняем шаблон с данными
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона %s: %v", tmpl, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return err
	}

	return nil
}
