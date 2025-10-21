package controllers

import (
	"net/http"

	"github.com/ADEMOLA200/string-analyzer/cmd/internal/models"
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/services"

	"github.com/gin-gonic/gin"
)

type StringController struct {
	service services.AnalyzerService
}

func NewStringController(service services.AnalyzerService) *StringController {
	return &StringController{service: service}
}

func (c *StringController) CreateAnalyzeString(ctx *gin.Context) {
	var req models.AnalysisRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Value == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'value' field"})
		return
	}

	analysis, err := c.service.AnalyzeString(req.Value)
	if err != nil {
		if err.Error() == "string already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "String already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, analysis)
}

func (c *StringController) GetString(ctx *gin.Context) {
	value := ctx.Param("string_value")

	if value == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "String value is required"})
		return
	}

	analysis, err := c.service.GetStringByValue(value)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "String not found"})
		return
	}

	ctx.JSON(http.StatusOK, analysis)
}

func (c *StringController) GetAllStrings(ctx *gin.Context) {
	var filters models.FilterParams

	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	results, appliedFilters, err := c.service.GetAllStrings(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := models.APIResponse{
		Data:           results,
		Count:          len(results),
		FiltersApplied: appliedFilters,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *StringController) FilterByNaturalLanguage(ctx *gin.Context) {
	var query models.NaturalLanguageQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	results, interpretedQuery, err := c.service.ProcessNaturalLanguage(query.Query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse natural language query"})
		return
	}

	response := models.APIResponse{
		Data:             results,
		Count:            len(results),
		InterpretedQuery: interpretedQuery,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *StringController) DeleteString(ctx *gin.Context) {
	value := ctx.Param("string_value")

	if value == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "String value is required"})
		return
	}

	err := c.service.DeleteString(value)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "String not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
