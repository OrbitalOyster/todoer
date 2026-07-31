package users

type UserRole int

const (
	user UserRole = iota
	moderator
	admin
)

type User struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
	Role string `yaml:"role"`
}
