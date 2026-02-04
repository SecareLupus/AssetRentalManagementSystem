package domain

import (
	"time"
)

type FieldType string

const (
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeString  FieldType = "text"
	FieldTypeImage   FieldType = "image"
)

type InspectionField struct {
	ID           int64     `json:"id"`
	TemplateID   int64     `json:"template_id"`
	Label        string    `json:"label"`
	Type         FieldType `json:"type"`
	Required     bool      `json:"required"`
	DisplayOrder int       `json:"display_order"`
}

type InspectionTemplate struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Fields      []InspectionField `json:"fields,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type InspectionSubmission struct {
	ID          int64                `json:"id"`
	AssetID     int64                `json:"asset_id"`
	TemplateID  int64                `json:"template_id"`
	PerformedBy string               `json:"performed_by"`
	Responses   []InspectionResponse `json:"responses"`
	CreatedAt   time.Time            `json:"created_at"`
}

type InspectionResponse struct {
	ID           int64  `json:"id"`
	SubmissionID int64  `json:"submission_id"`
	FieldID      int64  `json:"field_id"`
	Value        string `json:"value"` // Stores bool, text, or image URL
}
