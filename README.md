# String Analyzer Service

A RESTful API built with Go for analyzing and managing strings. This service computes various string properties and stores them for retrieval, forming part of the **HNG Internship Stage 1 Task**.

## Features

* Analyze and store string properties
* Retrieve and filter analyzed strings
* Query using natural language
* Identify strings uniquely via SHA256 hash

## Computed Properties

* **length** – total number of characters
* **is_palindrome** – whether the string reads the same backward and forward
* **unique_characters** – count of distinct characters
* **word_count** – total number of words
* **sha256_hash** – SHA-256 hash for unique identification
* **character_frequency_map** – frequency of each character

## API Endpoints

### 1. Analyze a String

```http
POST /api/v1/strings
Content-Type: application/json

{
  "value": "string to analyze"
}
```

### 2. Get a Specific String

```http
GET /api/v1/strings/{string_value}
```

### 3. Get All Strings with Filters

```http
GET /api/v1/strings?is_palindrome=true&min_length=5&max_length=20&word_count=2&contains_character=a
```

### 4. Natural Language Filtering

```http
GET /api/v1/strings/filter-by-natural-language?query=all%20single%20word%20palindromic%20strings
```

### 5. Delete a String

```http
DELETE /api/v1/strings/{string_value}
```

## Setup and Installation

### Prerequisites

* Go **1.21+**

### Local Development

Clone the repository:

```bash
git clone <https://github.com/ADEMOLA200/string-analyzer.git>
cd string-analyzer
```

Install dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run cmd/api/main.go
```

Server starts at **[http://localhost:8080](http://localhost:8080)**

### Environment Variables

* `PORT` – Port number for the server (default: `8080`)
