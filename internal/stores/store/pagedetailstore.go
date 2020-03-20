package store

import "github.com/worlve/sp-service/internal/models/pagedetail"

// PageDetailStore defines the required functionality for any associated store.
type PageDetailStore interface {
	UpdatePageDetail(record pagedetail.PageDetail) error
}
