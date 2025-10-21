package utils

import (
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/models"
)

func CreateSuccessResponse(data interface{}, count int, filters interface{}) models.APIResponse {
	return models.APIResponse{
		Data:           data,
		Count:          count,
		FiltersApplied: filters,
	}
}

func CreateNaturalLanguageResponse(data interface{}, count int, interpretedQuery *models.InterpretedQuery) models.APIResponse {
	return models.APIResponse{
		Data:             data,
		Count:            count,
		InterpretedQuery: interpretedQuery,
	}
}

func CreateErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"error": message,
	}
}
