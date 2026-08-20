package model

type (
	Template struct {
		ID          int64        `json:"id"`
		Type        TemplateType `json:"type"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		VideoLinks  []string     `json:"video_links"`
	}

	TemplateType string
)

const (
	TemplateTypePower TemplateType = "power"
	TemplateTypeCross TemplateType = "cross"
)
