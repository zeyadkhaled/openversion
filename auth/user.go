package auth

type User struct {
	Source string `json:"source,omitempty"`
	ID     string `json:"id,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Name   string `json:"name,omitempty"`
	Roles  []Role `json:"roles,omitempty"`
}

func (u User) HasRole(r Role) bool {
	for i := range u.Roles {
		if u.Roles[i] == r {
			return true
		}
	}
	return false
}

const (
	SourceRiders    = "riders"
	SourceAdmins    = "admins"
	SourceServer    = "server"
	SourceAnonymous = "anonymous"
)

// FullPerm is a user with full permissions.
//
// TODO(ebati): bu blok olmasina gerek yok ihtiyac olan kendi olusturmali
var FullPerm = User{
	Source: SourceServer,
	ID:     "1134ae6f-bd7f-50ff-9dd4-f2c6d79caa3e",
	Phone:  "+908508855467",
	Name:   "server",
	Roles:  FullPermList,
}
