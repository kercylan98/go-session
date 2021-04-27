package session

import "testing"

func TestSession_GetId(t *testing.T) {
	session, _ := sm.RegisterSession("session-0001")
	t.Log(session.GetId())
}

func TestSession_Close(t *testing.T) {
	session, _ := sm.RegisterSession("session-0001")
	sm.RegisterSession("session-0002")
	sm.RegisterSession("session-0003")

	session.Close()

	for id, session := range sm.sessions {
		t.Log(id, session)
	}

}

func TestSession_Store(t *testing.T) {
	session, _ := sm.RegisterSession("session-0001")
	session.Store("prop", 1)
}

func TestSession_Load(t *testing.T) {
	session, _ := sm.RegisterSession("session-0001")
	session.Store("prop", 1)
	t.Log(session.Load("prop"))
}

func TestSession_Del(t *testing.T) {
	session, _ := sm.RegisterSession("session-0001")
	session.Store("prop", 1)
	t.Log(session.Load("prop"))
	session.Del("prop")
	t.Log(session.Load("prop"))
}
