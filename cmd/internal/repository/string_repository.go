package repository

import (
	"fmt"
	"sync"

	"github.com/ADEMOLA200/string-analyzer/cmd/internal/models"
)

var (
	ErrNotFound        = fmt.Errorf("string not found")
	ErrDuplicateString = fmt.Errorf("string already exists")
)

type StringRepository interface {
	Create(analysis *models.StringAnalysis) error
	FindByID(id string) (*models.StringAnalysis, error)
	FindByValue(value string) (*models.StringAnalysis, error)
	FindAll() ([]*models.StringAnalysis, error)
	FindWithFilters(filters map[string]interface{}) ([]*models.StringAnalysis, error)
	DeleteByValue(value string) error
}

type inMemoryRepository struct {
	mu      sync.RWMutex
	byID    map[string]*models.StringAnalysis
	byValue map[string]*models.StringAnalysis
}

func NewInMemoryRepository() StringRepository {
	return &inMemoryRepository{
		byID:    make(map[string]*models.StringAnalysis),
		byValue: make(map[string]*models.StringAnalysis),
	}
}

func (r *inMemoryRepository) Create(analysis *models.StringAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byValue[analysis.Value]; exists {
		return ErrDuplicateString
	}

	r.byID[analysis.ID] = analysis
	r.byValue[analysis.Value] = analysis
	return nil
}

func (r *inMemoryRepository) FindByID(id string) (*models.StringAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	analysis, exists := r.byID[id]
	if !exists {
		return nil, ErrNotFound
	}
	return analysis, nil
}

func (r *inMemoryRepository) FindByValue(value string) (*models.StringAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	analysis, exists := r.byValue[value]
	if !exists {
		return nil, ErrNotFound
	}
	return analysis, nil
}

func (r *inMemoryRepository) FindAll() ([]*models.StringAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	analyses := make([]*models.StringAnalysis, 0, len(r.byID))
	for _, analysis := range r.byID {
		analyses = append(analyses, analysis)
	}
	return analyses, nil
}

func (r *inMemoryRepository) FindWithFilters(filters map[string]interface{}) ([]*models.StringAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.StringAnalysis

	for _, analysis := range r.byID {
		if matchesFilters(analysis, filters) {
			results = append(results, analysis)
		}
	}

	return results, nil
}

func (r *inMemoryRepository) DeleteByValue(value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	analysis, exists := r.byValue[value]
	if !exists {
		return ErrNotFound
	}

	delete(r.byID, analysis.ID)
	delete(r.byValue, value)
	return nil
}

func matchesFilters(analysis *models.StringAnalysis, filters map[string]interface{}) bool {
	for key, value := range filters {
		switch key {
		case "is_palindrome":
			if expected, ok := value.(bool); ok && analysis.Properties.IsPalindrome != expected {
				return false
			}
		case "min_length":
			if min, ok := value.(int); ok && analysis.Properties.Length < min {
				return false
			}
		case "max_length":
			if max, ok := value.(int); ok && analysis.Properties.Length > max {
				return false
			}
		case "word_count":
			if count, ok := value.(int); ok && analysis.Properties.WordCount != count {
				return false
			}
		case "contains_character":
			if char, ok := value.(string); ok && len(char) == 1 {
				charStr := string(char[0])
				if _, exists := analysis.Properties.CharacterFrequency[charStr]; !exists {
					return false
				}
			}
		}
	}
	return true
}
