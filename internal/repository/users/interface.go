package users

type Users interface {
	AddUserID(user string) error
	CheckUserID(user string) (bool, error)
}
