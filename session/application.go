package session

var application Manager

func init() {
	application = NewManager()
}

func GetManager() Manager {
	return application
}
