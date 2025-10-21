package models

import "time"

type StringAnalysis struct {
	ID         string           `json:"id"`
	Value      string           `json:"value"`
	Properties StringProperties `json:"properties"`
	CreatedAt  time.Time        `json:"created_at"`
}

type StringProperties struct {
	Length             int            `json:"length"`
	IsPalindrome       bool           `json:"is_palindrome"`
	UniqueCharacters   int            `json:"unique_characters"`
	WordCount          int            `json:"word_count"`
	SHA256Hash         string         `json:"sha256_hash"`
	CharacterFrequency map[string]int `json:"character_frequency_map"`
}

type AnalysisRequest struct {
	Value string `json:"value" binding:"required"`
}

type FilterParams struct {
	IsPalindrome      *bool   `form:"is_palindrome"`
	MinLength         *int    `form:"min_length"`
	MaxLength         *int    `form:"max_length"`
	WordCount         *int    `form:"word_count"`
	ContainsCharacter *string `form:"contains_character"`
}

type NaturalLanguageQuery struct {
	Query string `form:"query" binding:"required"`
}

type APIResponse struct {
	Data             interface{}       `json:"data,omitempty"`
	Count            int               `json:"count,omitempty"`
	FiltersApplied   interface{}       `json:"filters_applied,omitempty"`
	InterpretedQuery *InterpretedQuery `json:"interpreted_query,omitempty"`
}

type InterpretedQuery struct {
	Original      string                 `json:"original"`
	ParsedFilters map[string]interface{} `json:"parsed_filters"`
}
