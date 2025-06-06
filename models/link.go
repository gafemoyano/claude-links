package models

import "time"

type Link struct {
	ID           string    `json:"id" db:"id"`
	DeepLink     string    `json:"deep_link" db:"deep_link"`
	IOSStore     string    `json:"ios_store" db:"ios_store"`
	AndroidStore string    `json:"android_store" db:"android_store"`
	Title        string    `json:"title,omitempty" db:"title"`
	Description  string    `json:"description,omitempty" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ClickCount   int       `json:"click_count" db:"click_count"`
}

type CreateLinkRequest struct {
	DeepLink     string `json:"deep_link" validate:"required"`
	IOSStore     string `json:"ios_store" validate:"required,url"`
	AndroidStore string `json:"android_store" validate:"required,url"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
}