package users

var list []User

func Load(newList []User) {
	list = newList
}

func GetAllUsers() []User {
	return list;
}
