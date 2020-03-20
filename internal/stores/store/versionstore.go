package store

import "github.com/worlve/sp-service/internal/models/version"

// VersionStore defines the required functionality for any associated store.
type VersionStore interface {
	GetVersion(versionGUID string) (version.Version, error)
}
