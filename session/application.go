package session

var application Manager

func init() {
	application = NewManager()
}

func GetApplication() Manager {
	return application
}

// 对外实现
type Application interface {
	// 获取Session管理器
	GetSessionManager() Manager
}