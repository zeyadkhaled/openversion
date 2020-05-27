package version

import (
	"time"
)

type Application struct {
	// ID is made from appname_platform
	ID string `json:"id,omitempty"`
	// version isn't used in case we need to do DB related operations
	MinVersion string `json:"code,omitempty"`
	Package    string `json:"package,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
