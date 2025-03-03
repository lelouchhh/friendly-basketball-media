package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// CORSWithConfig - конфигурируемая версия middleware для CORS
func CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Разрешаем доступ с любых доменов
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		// Указываем, какие методы разрешены
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// Указываем, какие заголовки разрешены
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Разрешаем использовать credentials (если нужно)
		// c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

		// Обработка Preflight запроса (OPTIONS)
		if c.Request().Method == http.MethodOptions {
			// Для OPTIONS запроса возвращаем успешный ответ с нужными заголовками
			return c.NoContent(http.StatusNoContent)
		}

		// Продолжаем выполнение запроса
		return next(c)
	}
}
