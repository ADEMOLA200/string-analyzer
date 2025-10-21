package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/ADEMOLA200/string-analyzer/cmd/internal/models"
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/repository"
)

type AnalyzerService interface {
	AnalyzeString(value string) (*models.StringAnalysis, error)
	GetStringByValue(value string) (*models.StringAnalysis, error)
	GetAllStrings(filters models.FilterParams) ([]*models.StringAnalysis, map[string]interface{}, error)
	ProcessNaturalLanguage(query string) ([]*models.StringAnalysis, *models.InterpretedQuery, error)
	DeleteString(value string) error
}

type analyzerService struct {
	repo repository.StringRepository
}

func NewAnalyzerService(repo repository.StringRepository) AnalyzerService {
	return &analyzerService{repo: repo}
}

func (s *analyzerService) AnalyzeString(value string) (*models.StringAnalysis, error) {
	if existing, _ := s.repo.FindByValue(value); existing != nil {
		return nil, repository.ErrDuplicateString
	}

	properties := s.analyzeStringProperties(value)

	analysis := &models.StringAnalysis{
		ID:         properties.SHA256Hash,
		Value:      value,
		Properties: properties,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.Create(analysis); err != nil {
		return nil, err
	}

	return analysis, nil
}

func (s *analyzerService) analyzeStringProperties(value string) models.StringProperties {
	hash := sha256.Sum256([]byte(value))
	sha256Hash := hex.EncodeToString(hash[:])

	cleanStr := strings.ToLower(value)
	cleanStr = strings.ReplaceAll(cleanStr, " ", "")

	charFreq := make(map[string]int)
	for _, char := range value {
		charStr := string(char)
		charFreq[charStr]++
	}

	words := strings.Fields(value)
	wordCount := len(words)

	uniqueChars := len(charFreq)

	return models.StringProperties{
		Length:             len(value),
		IsPalindrome:       isPalindrome(cleanStr),
		UniqueCharacters:   uniqueChars,
		WordCount:          wordCount,
		SHA256Hash:         sha256Hash,
		CharacterFrequency: charFreq,
	}
}

func isPalindrome(s string) bool {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

func (s *analyzerService) GetStringByValue(value string) (*models.StringAnalysis, error) {
	return s.repo.FindByValue(value)
}

func (s *analyzerService) GetAllStrings(filters models.FilterParams) ([]*models.StringAnalysis, map[string]interface{}, error) {
	filterMap := make(map[string]interface{})

	if filters.IsPalindrome != nil {
		filterMap["is_palindrome"] = *filters.IsPalindrome
	}
	if filters.MinLength != nil {
		filterMap["min_length"] = *filters.MinLength
	}
	if filters.MaxLength != nil {
		filterMap["max_length"] = *filters.MaxLength
	}
	if filters.WordCount != nil {
		filterMap["word_count"] = *filters.WordCount
	}
	if filters.ContainsCharacter != nil && len(*filters.ContainsCharacter) == 1 {
		filterMap["contains_character"] = *filters.ContainsCharacter
	}

	results, err := s.repo.FindWithFilters(filterMap)
	return results, filterMap, err
}

func (s *analyzerService) ProcessNaturalLanguage(query string) ([]*models.StringAnalysis, *models.InterpretedQuery, error) {
	parsedFilters := s.parseNaturalLanguage(query)

	if len(parsedFilters) == 0 {
		return nil, nil, fmt.Errorf("unable to parse natural language query")
	}

	results, err := s.repo.FindWithFilters(parsedFilters)
	if err != nil {
		return nil, nil, err
	}

	interpretedQuery := &models.InterpretedQuery{
		Original:      query,
		ParsedFilters: parsedFilters,
	}

	return results, interpretedQuery, nil
}

func (s *analyzerService) parseNaturalLanguage(query string) map[string]interface{} {
	query = strings.ToLower(query)
	filters := make(map[string]interface{})

	// Parse for palindrome
	if strings.Contains(query, "palindrom") {
		filters["is_palindrome"] = true
	}

	// Parse for word count
	if strings.Contains(query, "single word") || strings.Contains(query, "one word") {
		filters["word_count"] = 1
	} else if strings.Contains(query, "two words") || strings.Contains(query, "2 words") {
		filters["word_count"] = 2
	} else if strings.Contains(query, "three words") || strings.Contains(query, "3 words") {
		filters["word_count"] = 3
	}

	// Parse for length
	if strings.Contains(query, "longer than") || strings.Contains(query, "greater than") {
		if num := extractNumber(query); num > 0 {
			filters["min_length"] = num + 1
		}
	} else if strings.Contains(query, "shorter than") || strings.Contains(query, "less than") {
		if num := extractNumber(query); num > 0 {
			filters["max_length"] = num - 1
		}
	}

	// Parse for character containment
	for _, char := range query {
		if unicode.IsLetter(char) && strings.Contains(query, "contains") ||
			strings.Contains(query, "with") || strings.Contains(query, "having") {
			charStr := string(char)
			if len(charStr) == 1 && unicode.IsLetter(char) {
				filters["contains_character"] = charStr
				break
			}
		}
	}

	return filters
}

func extractNumber(s string) int {
	// Simple number extraction - you might want to enhance this
	words := strings.Fields(s)
	for i, word := range words {
		if word == "than" && i+1 < len(words) {
			switch words[i+1] {
			case "10", "ten":
				return 10
			case "20", "twenty":
				return 20
			case "5", "five":
				return 5
			}
		}
	}
	return 0
}

func (s *analyzerService) DeleteString(value string) error {
	return s.repo.DeleteByValue(value)
}
