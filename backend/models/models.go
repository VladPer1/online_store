package models

type Category struct {
	ID   int
	Name string
}

type Producer struct {
	ID      int
	Name    string
	Country string
}

type User struct {
	ID       int
	Email    string
	Password string
	Name     string
	Status   string
}

type Product struct {
	ProductID     int
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	Image         string
}

type ProductForFilter struct {
	ProductID     int
	Name          string
	Price         float64
	StockQuantity int
	CategoryID    int
	ProducerID    int
	CategoryName  string
	ProducerName  string
	Image         string
}

type FilterParams struct {
	Categories  []int
	Producers   []int
	MinPrice    float64
	MaxPrice    float64
	SortByPrice string // "asc", "desc"
	SearchQuery string
}

type UserCammrt struct {
	ProductID int
	Quantity  int
	Price     float32
}

type CartItem struct {
	ProductID    int
	ProductName  string
	Price        float32
	Quantity     int
	ImageURL     string
	TotalPrice   float32
	Description  string
	CategoryName string
	ProducerName string
}

// Структура для передачи данных в шаблон
type TemplateData struct {
	Message  string
	Success  bool
	FormType string // "login" или "register"
	Email    string
	Name     string
}

// Структура для данных каталога
type CatalogData struct {
	Title              string
	Products           []ProductForFilter
	Categories         []Category // Добавляем категории
	Producers          []Producer // Добавляем производителей
	SearchQuery        string
	SelectedCategories StringSet
	SelectedBrands     StringSet
	PriceMin           string
	PriceMax           string
	SortBy             string
}

type ProfileData struct {
	UserName     string
	UserEmail    string
	RecentOrders []OrdersProfile
}

type OrdersProfile struct {
	OrderID     int
	ProductName string
	Quantity    int
	TotalAmount float32
}

type Order struct {
	UserID      int
	TotalAmount int

	ProductID int
	Quantity  int
	UnitPrice int
}

type StringSet map[string]bool

// MakeStringSet создает StringSet из среза строк
func MakeStringSet(slice []string) StringSet {
	set := make(StringSet)
	for _, item := range slice {
		set[item] = true
	}
	return set
}

// Contains проверяет наличие ключа в StringSet
func (s StringSet) Contains(key string) bool {
	return s[key]
}

type ProductForAdmin struct {
	ProductID     int
	Name          string
	Description   string
	Price         float32
	StockQuantity int
	ImageURL      string
	CategoryID    int
	ProducerID    int
	CategoryName  string
	ProducerName  string
}

type Config struct {
	TemplatePath string
	StaticPath   string
}

type PlaceholderUser struct {
	ID               int
	Email            string
	Password         string
	Name             string
	VerificationCode int
}
