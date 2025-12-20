package main

import (
	"net/http"

	"online_store/database"
	"online_store/server"
	"online_store/utils"
)

func main() {
	db := database.Connect()
	defer db.Close()
	mux := http.NewServeMux()
	utils.InitConfig()
	utils.PutFiles(mux)

	// Запуск сервера с graceful shutdown
	server.RunServerWithShutdown(db, mux)
}
