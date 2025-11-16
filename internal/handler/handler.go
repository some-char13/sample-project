package handler

import (
	"net/http"
	"sample_project/internal/conf"
	"sample_project/internal/model/check"
	"sample_project/internal/model/register"
	"sample_project/internal/model/service"
	NewItem "sample_project/internal/service"
	Jwt "sample_project/pkg/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler структура, которая имеет методы для работы с каждым эндпоинтом.
// По желанию можно добавить поля (например различные валидаторы)
type Handler struct {
}

// New конструктор
func New() *Handler {
	return &Handler{}
}

// Create service
// @Summary Создание сервиса
// @Tags service
// @Accept json
// @Produce json
// @Param input body service.ServiceRequest true "Модель сервиса"
// @Success 204
// @Failure 400 {string} string "Service id already exists"
// @Failure      401  {string} string "invalid token"
// @Router /api/service/item [post]
// @Security BearerAuth
func CreateService(ctx *gin.Context) {

	req := &service.Service{}

	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req = service.NewService(req.Id, req.Name, req.Url, req.Interval)
	search := NewItem.SearchServiceItem(req.Id)
	if search == nil {
		// ch := make(chan any, 1)

		// ch <- req
		// NewItem.ProcessItems(ch)

		NewItem.ProcessItems(req)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Service id already exists",
			"id":      req.Id,
		})
	}

	ctx.Status(http.StatusNoContent)

}

// Create result
// @Summary Создание  результата провреки
// @Tags result
// @Accept json
// @Produce json
// @Param input body check.ResultRequest true "Модель результата провреки"
// @Success      204
// @Failure 400 {string} string "Result id already exists"
// @Failure      401  {string} string "invalid token"
// @Router        /api/result/item [post]
// @Security BearerAuth
func CreateResult(ctx *gin.Context) {

	req := &check.Result{}

	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req = check.NewResult(req.Id, req.ServiceId, req.ResponseCode, req.RespDuration)
	search := NewItem.SearchResultItem(req.Id)
	if search == nil {
		// ch := make(chan any, 1)

		// ch <- req
		// NewItem.ProcessItems(ch)

		NewItem.ProcessItems(req)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Result id already exists",
			"id":      req.Id,
		})
	}

	ctx.Status(http.StatusNoContent)

}

// Get services
// @Summary Получить список всех сервисов
// @Tags service
// @Produce json
// @Success 200
// @Failure 400 {string} string "Item not found"
// @Router /api/service/items [get]
func GetService(ctx *gin.Context) {
	services := NewItem.GetServices()

	ctx.JSON(http.StatusOK, services)
}

// @Summary Получить список всех результатов
// @Tags result
// @Produce json
// @Success 200
// @Failure 400 {string} string "Item not found"
// @Router /api/result/items [get]
func GetResult(ctx *gin.Context) {
	results := NewItem.GetResults()

	ctx.JSON(http.StatusOK, results)
}

// @Summary Поиск сервиса по ID
// @Tags service
// @Produce json
// @Param id path int true "ID сервиса"
// @Success 200 {object} service.Service
// @Failure 400 {string} string "Item not found"
// @Router /api/service/item/{id} [get]
func SearchServiceId(ctx *gin.Context) {

	id := ctx.Param("id")

	converted, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	search := NewItem.SearchServiceItem(converted)
	if search == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Item not found",
			"id":      converted,
		})
	} else {
		ctx.JSON(http.StatusOK, search)
	}

}

// @Summary Поиск результата по ID
// @Tags result
// @Produce json
// @Param id path int true "ID результата"
// @Success 200 {object} check.Result
// @Failure 400 {string} string "Item not found"
// @Router /api/result/item/{id} [get]
func SearchResultId(ctx *gin.Context) {

	id := ctx.Param("id")

	converted, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	search := NewItem.SearchResultItem(converted)
	if search == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Item not found",
			"id":      converted,
		})
	} else {
		ctx.JSON(http.StatusOK, search)
	}

}

// @Summary Удаление сервиса по ID
// @Tags service
// @Produce json
// @Param id path int true "ID сервиса"
// @Success 204
// @Failure 400 {string} string "Item not found"
// @Failure 401 {string} string "Invalid token"
// @Router /api/service/item/{id} [delete]
// @Security BearerAuth
func DeleteService(ctx *gin.Context) {

	id := ctx.Param("id")

	converted, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	check := NewItem.SearchServiceItem(converted)
	if check == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Item not found",
			"id":      converted,
		})
	} else {
		NewItem.DeleteItemService(converted)
	}

	ctx.Status(http.StatusNoContent)

}

// @Summary Удаление результата по ID
// @Tags result
// @Produce json
// @Param id path int true "ID результата"
// @Success 204
// @Failure 400 {string} string "Item not found"
// @Failure 401 {string} string "Invalid token"
// @Router /api/result/item/{id} [delete]
// @Security BearerAuth
func DeleteResult(ctx *gin.Context) {

	id := ctx.Param("id")

	converted, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	check := NewItem.SearchResultItem(converted)
	if check == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Item not found",
			"id":      converted,
		})
	} else {
		NewItem.DeleteItemResult(converted)
	}

	ctx.Status(http.StatusNoContent)

}

// @Summary Изменение сервиса по ID
// @Tags service
// @Accept json
// @Produce json
// @Param id path int true "ID сервиса"
// @Param input body service.ServiceRequest true "Модель сервиса"
// @Success 204
// @Failure 400 {string} string "Item not found"
// @Failure 401 {string} string "Invalid token"
// @Router /api/service/item/{id} [put]
// @Security BearerAuth
func ChangeService(ctx *gin.Context) {
	id := ctx.Param("id")

	converted, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	req := &service.Service{}
	req.Id = converted
	err = ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	check := NewItem.SearchServiceItem(converted)
	if check == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Item not found",
			"id":      converted,
		})
	} else {

		NewItem.ChangeItems(converted, req)
	}

	ctx.Status(http.StatusNoContent)

}

// @Summary Изменение результата по ID
// @Tags result
// @Accept json
// @Produce json
// @Param id path int true "ID результата"
// @Param input body check.ResultRequest true "Модель результата"
// @Success 204
// @Failure 400 {string} string "Result item not found"
// @Failure 401 {string} string "Invalid token"
// @Router /api/result/item/{id} [put]
// @Security BearerAuth
func ChangeResult(ctx *gin.Context) {
	id := ctx.Param("id")

	converted, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	req := &check.Result{}
	req.Id = converted
	err = ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	check := NewItem.SearchResultItem(converted)
	if check == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Result item not found",
			"id":      converted,
		})
	} else {

		NewItem.ChangeItems(converted, req)
	}

	ctx.Status(http.StatusNoContent)

}

// SignIn Аутентификация и получение JWT токена
// @Summary Аутентификация
// @Tags auth
// @Accept json
// @Produce json
// @Param input body register.LoginRequest true "Модель аутентификации"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid username or password"
// @Router /api/sign_in [post]
func SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		name := conf.Load().User
		pass := conf.Load().Pass

		var req register.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		// Проверяем учетные данные из конфигурации
		if req.Username != name || req.Password != pass {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
			return
		}

		token, err := Jwt.New().GenerateToken(req.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	}
}

// // Register
// // @Summary Регистрация
// // @Tags auths
// // @Accept			json
// // @Produce		json
// // @Param input body register.RegisterRequest true "Модель которую принимает метод"
// // @Success 200 {string}  string "Registration successful"
// // @Failure 400 {string} string "Invalid request"
// // @Router /api/create_user [post]
// func (h *Handler) Register() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var req register.RegisterRequest

// 		if err := ctx.ShouldBindJSON(&req); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
// 			return
// 		}

// 		// Проверка, что имя пользователя не существует
// 		if _, exists := users[req.Username]; exists {
// 			ctx.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
// 			return
// 		}

// 		// Регистрируем пользователя
// 		users[req.Username] = req.Password
// 		ctx.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
// 	}
// }
