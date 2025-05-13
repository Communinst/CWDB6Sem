package entities

type Role struct {
	RoleId            int    `json:"role_id" db:"role_id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	SignificanceOrder int    `json:"significance_order" db:"significance_order"`
}

var (
	Bot     Role = Role{0, "Bot", "", -1}
	Player  Role = Role{0, "User", "User", 0}
	Stuff   Role = Role{0, "Stuff", "", 1}
	Manager Role = Role{0, "Manager", "", 2}
	Chief   Role = Role{0, "Chief", "", 3}
)
