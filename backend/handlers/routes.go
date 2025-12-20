package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(mux *http.ServeMux, db *pgxpool.Pool) {
	// Маршруты

	mux.HandleFunc("/", serveLoginPage)
	mux.HandleFunc("/login", loginFormHandler(db))
	mux.HandleFunc("/register", registerFormHandler(db))

	mux.HandleFunc("/forgot-password", ServeForgotPasswordPage(db))                                    // GET - форма email
	mux.HandleFunc("/forgot-password-send", ForgotPasswordHandler(db))                                 // POST - отправить код
	mux.HandleFunc("/forgot-password-verify-page", ServeForgotPasswordVerifyPage(db))                  // GET - форма кода
	mux.HandleFunc("/forgot-password-verify", ForgotPasswordVerifyHandler(db))                         // POST - проверить код
	mux.HandleFunc("/forgot-password-update-password-page", ServeForgotPasswordUpdatePasswordPage(db)) // GET - форма пароля
	mux.HandleFunc("/forgot-password-update-password", ForgotPasswordUpdatePasswordHandler(db))        // POST - сохранить пароль

	mux.HandleFunc("/verify-code-page", ServeVerifyCodePage(db))
	mux.HandleFunc("/verify-code", verifyCodeHandler(db))
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/catalog", serveCatalog(db)) // Передаем db в обработчик
	mux.HandleFunc("/add-to-cart", AddProductToCartHandler(db))
	mux.HandleFunc("/cart", CartHandler(db))
	mux.HandleFunc("/delete-from-cart", DeleteProductFromCartHandler(db))
	mux.HandleFunc("/profile", ProfileHandler(db))

	mux.HandleFunc("/update-profile-page", ServeChangeProfilePage(db))
	mux.HandleFunc("/update-profile", UpdateProfileHandler(db))

	mux.HandleFunc("/change-password-page", ServeChangePassword(db))
	mux.HandleFunc("/change-password", ChangePasswordHandler(db))

	mux.HandleFunc("/delete-account-page", DeleteAccountPage(db))
	mux.HandleFunc("/delete-account", DeleteAccountHandler(db))

	mux.HandleFunc("/process-payment", ProcessPaymentHandler(db))
	mux.HandleFunc("/success-payment", SuccessPaymentHandler())
	mux.HandleFunc("/error-payment", ErrorPaymentHandler())

	mux.HandleFunc("/admin", AdminPanelHandler(db))
	mux.HandleFunc("/admin/products/save", AdminSaveProductHandler(db))
	mux.HandleFunc("/admin/products/delete", AdminDeleteProductHandler(db))
	mux.HandleFunc("/admin/users/ban", AdminBanUserHandler(db))
	mux.HandleFunc("/admin/users/unban", AdminUnbanUserHandler(db))
	mux.HandleFunc("/admin/report", GenerateReportHandler(db))

	mux.HandleFunc("/admin/export-csv", ExportProductsCSVHandler(db))
	mux.HandleFunc("/admin/import-csv", ImportProductsCSVHandler(db))

}
