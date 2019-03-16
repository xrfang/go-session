package session

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	cfg Config
	reg map[string]Session
	sync.RWMutex
}

func NewManager(cfg *Config) *Manager {
	sm := Manager{reg: make(map[string]Session)}
	sm.setConfig(cfg)
	sm.loadSessions()
	go func() {
		for {
			time.Sleep(time.Minute)
			sm.Lock()
			for n, s := range sm.reg {
				if s.TTL() <= 0 {
					delete(sm.reg, n)
				}
			}
			sm.Unlock()
			if sm.cfg.Persist != "" {
				sm.saveSessions()
			}
		}
	}()
	return &sm
}

func (sm *Manager) Get(w http.ResponseWriter, r *http.Request) (s Session) {
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
	if s.ID == "" || s.src != from || s.TTL() <= 0 {
		s = Session{
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
		Name:     "session",
		Value:    s.ID,
		Path:     "/",
		MaxAge:   s.TTL(),
		Secure:   r.TLS != nil,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return
}
