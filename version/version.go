package version

import (
	"time"
)

type Application struct {
	// ID is made from appname_platform
	ID string `json:"id,omitempty"`
	// Holds the minimal version for this application
	MinVersion string `json:"code,omitempty"`
	// It is the package details i.e com.zeyad.thisappname
	Package string `json:"package,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
