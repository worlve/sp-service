package store

import (
	"github.com/worlve/sp-service/internal/models/pagetemplate"
)

// PageTemplateStore defines the required functionality for any associated store.
type PageTemplateStore interface {
	GetPageTemplate(pageTemplateGUID string) (pagetemplate.PageTemplate, error)
}
