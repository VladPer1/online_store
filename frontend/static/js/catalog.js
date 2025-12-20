// Функции для работы с каталогом

function applyFilters() {
    const search = document.getElementById('search').value;
    const priceMin = document.querySelector('input[name="price_min"]').value;
    const priceMax = document.querySelector('input[name="price_max"]').value;
    const sortBy = document.getElementById('sortBy').value;
    
    const params = new URLSearchParams();
    if (search) params.append('search', search);
    if (priceMin) params.append('price_min', priceMin);
    if (priceMax) params.append('price_max', priceMax);
    if (sortBy) params.append('sortBy', sortBy);
    
    document.querySelectorAll('input[name="category"]:checked').forEach(cb => {
        params.append('category', cb.value);
    });
    
    document.querySelectorAll('input[name="producer"]:checked').forEach(cb => {
        params.append('producer', cb.value);
    });
    
    window.location.href = '/catalog?' + params.toString();
}

function clearFilters() {
    window.location.href = '/catalog';
}

// Функции для управления количеством
function changeQuantity(productId, change) {
    const input = document.getElementById('qty-' + productId);
    const currentValue = parseInt(input.value) || 1;
    const stockQuantity = parseInt(input.max) || 999;
    
    let newValue = currentValue + change;
    
    // Проверяем границы
    if (newValue < 1) newValue = 1;
    if (newValue > stockQuantity) newValue = stockQuantity;
    
    input.value = newValue;
}

function validateQuantity(productId, stockQuantity) {
    const input = document.getElementById('qty-' + productId);
    let value = parseInt(input.value) || 1;
    
    if (value < 1) value = 1;
    if (value > stockQuantity) value = stockQuantity;
    
    input.value = value;
}

// Инициализация при загрузке страницы
function initCatalog() {
    // Обработчики событий для кнопок +/-
    document.querySelectorAll('.qty-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            const productId = this.getAttribute('data-product-id');
            const change = parseInt(this.getAttribute('data-change'));
            changeQuantity(productId, change);
        });
    });

    // Валидация при ручном вводе
    document.querySelectorAll('.qty-input').forEach(input => {
        input.addEventListener('change', function() {
            const productId = this.id.replace('qty-', '');
            const stockQuantity = parseInt(this.max) || 999;
            validateQuantity(productId, stockQuantity);
        });
    });

    // Автопоиск с задержкой
    const searchInput = document.getElementById('search');
    if (searchInput) {
        let searchTimeout;
        searchInput.addEventListener('input', function(e) {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(applyFilters, 500);
        });
    }

    // Обработчик сортировки
    const sortSelect = document.getElementById('sortBy');
    if (sortSelect) {
        sortSelect.addEventListener('change', applyFilters);
    }
}

// Запуск инициализации при загрузке документа
document.addEventListener('DOMContentLoaded', initCatalog);