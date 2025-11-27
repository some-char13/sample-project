package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"sample_project/internal/conf"
	"sample_project/internal/model/register"
	"sample_project/internal/model/service"
	Jwt "sample_project/internal/pkg/jwt"
	srv "sample_project/internal/service"
)

type Handler struct {
	service *srv.Service
	monitor *srv.MonitorService
}

func NewHandler(service *srv.Service, monitor *srv.MonitorService) *Handler {
	return &Handler{
		service: service,
		monitor: monitor,
	}
}

// CreateService создает новый сервис для мониторинга
// @Summary Создание сервиса
// @Tags services
// @Accept json
// @Produce json
// @Param input body service.Request true "Данные сервиса"
// @Success 201 {object} service.Service
// @Failure 400 {string} string "Invalid request data"
// @Failure 409 {string} string "Service with this name already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /api/services [post]
// @Security BearerAuth
func (h *Handler) CreateService(ctx *gin.Context) {
	var req service.Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc, err := h.service.CreateService(ctx.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == fmt.Sprintf("service with name '%s' already exists", req.Name) {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	h.monitor.StartMonitoring(svc)

	ctx.JSON(http.StatusCreated, svc)
}

// @Summary Получение списка сервисов
// @Tags services
// @Produce json
// @Success 200 {array} service.Service
// @Failure 500 {string} string "Internal server error"
// @Router /api/services [get]
func (h *Handler) GetServices(ctx *gin.Context) {
	services, err := h.service.GetServices(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get services"})
		return
	}

	ctx.JSON(http.StatusOK, services)
}

// GetServicesStatus возвращает статус всех сервисов
// @Summary Получение статуса всех сервисов
// @Tags services
// @Produce json
// @Success 200 {array} service.Status
// @Failure 500 {string} string "Internal server error"
// @Router /api/services/status [get]
func (h *Handler) GetServicesStatus(ctx *gin.Context) {
	statuses, err := h.service.GetAllServicesStatus(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get services status"})
		return
	}

	ctx.JSON(http.StatusOK, statuses)
}

// DeleteService удаляет сервис
// @Summary Удаление сервиса
// @Tags services
// @Param id path int true "ID сервиса"
// @Success 204
// @Failure 400 {string} string "Invalid service ID"
// @Failure 404 {string} string "Service not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/services/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	h.monitor.StopMonitoring(id)

	if err := h.service.DeleteService(ctx.Request.Context(), id); err != nil {
		if err.Error() == fmt.Sprintf("service with ID %d not found", id) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetServiceResults возвращает историю проверок для сервиса
// @Summary Получение истории проверок
// @Tags services
// @Param id path int true "ID сервиса"
// @Param limit query int false "Лимит результатов (по умолчанию 20)"
// @Produce json
// @Success 200 {array} check.Result
// @Failure 400 {string} string "Invalid service ID"
// @Failure 500 {string} string "Internal server error"
// @Router /api/services/{id}/results [get]
func (h *Handler) GetServiceResults(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	limit := 20
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	results, err := h.service.GetServiceResults(ctx.Request.Context(), id, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get check results"})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

// GetServiceResultsByStatus возвращает историю проверок по статусу
// @Summary Получение истории проверок по статусу
// @Tags services
// @Param id path int true "ID сервиса"
// @Param respCodeList query string true "Список кодов ответов через запятую"
// @Param limit query int false "Лимит результатов (по умолчанию 20)"
// @Produce json
// @Success 200 {array} check.Result
// @Failure 400 {string} string "Invalid parameters"
// @Failure 500 {string} string "Internal server error"
// @Router /api/services/{id}/results/filter [get]
func (h *Handler) GetServiceResultsByStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	respCodesStr := ctx.Query("respCodeList")
	if respCodesStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'respCodeList' is required"})
		return
	}

	respCodes, err := httpcode(respCodesStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 20
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'limit' must be a positive integer"})
			return
		}
	}

	results, err := h.service.GetServiceResultsByStatus(ctx.Request.Context(), id, respCodes, limit)
	if err != nil {
		log.Printf("Failed to get check results for service %d: %v", id, err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve check results",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"service_id": id,
		"results":    results,
	})
}

// SignIn аутентификация пользователя
// @Summary Аутентификация
// @Tags auth
// @Accept json
// @Produce json
// @Param input body register.LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Router /api/sign_in [post]
func (h *Handler) SignIn(ctx *gin.Context) {
	config := conf.Load()

	var req register.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Username != config.User || req.Password != config.Pass {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := Jwt.New().GenerateToken(req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// HealthCheck проверка API
// @Summary Проверка API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}

func httpcode(codesStr string) ([]int, error) {
	if strings.TrimSpace(codesStr) == "" {
		return nil, fmt.Errorf("response codes list cannot be empty")
	}

	codes := make([]int, 0, len(strings.Split(codesStr, ",")))

	for _, codeStr := range strings.Split(codesStr, ",") {
		codeStr = strings.TrimSpace(codeStr)
		if codeStr == "" {
			continue
		}

		code, err := strconv.Atoi(codeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid response code '%s': must be a number", codeStr)
		}

		if code < 100 || code > 599 {
			return nil, fmt.Errorf("invalid HTTP status code %d", code)
		}

		codes = append(codes, code)
	}

	if len(codes) == 0 {
		return nil, fmt.Errorf("no valid response codes provided")
	}

	return codes, nil
}
