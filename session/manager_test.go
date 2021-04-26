package session

import "testing"

var sm = NewManager()

func TestManager_RegisterSession(t *testing.T) {
	sm.RegisterSession("test-0001")
	for id, s := range sm.sessions {
		t.Log(id, s)
	}
}

func TestManager_UnRegisterSession(t *testing.T) {
	sm.RegisterSession("test-0001")
	sm.RegisterSession("test-0002")
	sm.RegisterSession("test-0003")
	t.Log("register ==========================> all session")
	for id, s := range sm.sessions {
		t.Log(id, s)
	}

	sm.UnRegisterSession(sm.sessions["test-0001"])

	t.Log("un register =======================> all session")
	for id, s := range sm.sessions {
		t.Log(id, s)
	}
}

func TestManager_GetAllSession(t *testing.T) {
	sm.RegisterSession("test-0001")
	sm.RegisterSession("test-0002")
	sm.RegisterSession("test-0003")

	for id, s := range sm.GetAllSession() {
		t.Log(id, s)
	}
}

func TestManager_GetSession(t *testing.T) {
	sm.RegisterSession("test-0001")

	t.Log(sm.GetSession("test-0001"))
}