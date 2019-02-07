package session

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	cfg Config
	reg map[string]*Session
	sync.RWMutex
}

func NewManager(cfg *Config) *Manager {
	sm := Manager{reg: make(map[string]*Session)}
	sm.setConfig(cfg)
	sm.loadSessions()
	go func() {
		for {
			time.Sleep(time.Minute)
			for n, s := range sm.reg {
				if s.ttl() <= 0 {
					sm.Lock()
					delete(sm.reg, n)
					sm.Unlock()
				}
			}
			sm.saveSessions()
		}
	}()
	return &sm
}

func (sm *Manager) Get(r *http.Request, w http.ResponseWriter) (s *Session) {
	from, _, _ := net.SplitHostPort(r.RemoteAddr)
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		c, err := r.Cookie("session")
		if err == nil {
			sessionID = c.Value
		}
	}
	if sessionID != "" {
		sm.RLock()
		s = sm.reg[sessionID]
		sm.RUnlock()
	}
	if s == nil || s.src != from || s.ttl() <= 0 {
		s = &Session{
			ID:  uuid(16),
			src: from,
			upd: time.Now(),
			mgr: sm,
			arg: make(map[string]interface{}),
		}
		sm.Lock()
		sm.reg[s.ID] = s
		sm.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  s.ID,
		Path:   "/",
		MaxAge: s.ttl(),
		Secure: r.TLS != nil,
	})
	return
}
