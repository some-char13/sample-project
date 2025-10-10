package handler

import (
	"net/http"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	NewItem "sample_project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	// req := &service.Service{}

	// err := ctx.ShouldBindJSON(req)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// req = service.NewService(req.Id, req.Name, req.Url, req.Interval)
	// search := NewItem.SearchItem(req.Id)
	// if search == nil {
	// 	ch := make(chan any, 1)

	// 	ch <- req
	// 	NewItem.ProcessItems(ch)
	// } else {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Item already exists",
	// 		"id":      req.Id,
	// 	})
	// }

	ctx.Status(http.StatusNoContent)

}

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

func GetService(ctx *gin.Context) {
	//services := repository.GetServices()
	services := NewItem.GetServices()

	ctx.JSON(http.StatusOK, services)
}

func GetResult(ctx *gin.Context) {
	results := NewItem.GetResults()

	ctx.JSON(http.StatusOK, results)
}

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
		// ch := make(chan any, 1)
		// ch <- req
		// NewItem.ChangeItems(converted, ch)

		NewItem.ChangeItems(converted, req)
	}

	ctx.Status(http.StatusNoContent)

}

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
		// ch := make(chan any, 1)
		// ch <- req
		// NewItem.ChangeItems(converted, ch)

		NewItem.ChangeItems(converted, req)
	}

	ctx.Status(http.StatusNoContent)

}
