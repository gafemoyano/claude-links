package models

import "time"

type Link struct {
	ID            string    `json:"id" db:"id"`
	UniversalLink string    `json:"universal_link" db:"universal_link"`
	DeepLink      string    `json:"deep_link" db:"deep_link"`
	IOSStore      string    `json:"ios_store" db:"ios_store"`
	AndroidStore  string    `json:"android_store" db:"android_store"`
	Title         string    `json:"title,omitempty" db:"title"`
	Description   string    `json:"description,omitempty" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	ClickCount    int       `json:"click_count" db:"click_count"`
}

type CreateLinkRequest struct {
	UniversalLink string `json:"universal_link" validate:"required"`
	DeepLink      string `json:"deep_link,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
}