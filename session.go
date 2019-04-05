package session

import (
	"encoding/json"
	"io"
	"time"
)

type Session struct {
	ID  string
	src string
	upd time.Time
	mgr *Manager
	arg map[string]string
}

func (s Session) Dump(w io.Writer) {
	je := json.NewEncoder(w)
	je.SetIndent("", "    ")
	je.Encode(map[string]interface{}{
		"id":  s.ID,
		"src": s.src,
		"upd": s.upd.Format(time.RFC3339),
		"arg": s.arg,
	})
}

func (s Session) TTL() int {
	ttl := s.mgr.cfg.TTL
	if len(s.arg) == 0 {
		ttl = s.mgr.cfg.VoidTTL
	}
	return ttl - int(time.Now().Sub(s.upd).Seconds())
}

func (s Session) Get(name string) string {
	v := s.arg[name]
	if v != "" && s.mgr.cfg.Refresh {
		s.upd = time.Now()
		s.mgr.Lock()
		s.mgr.reg[s.ID] = s
		s.mgr.Unlock()
	}
	return v
}

func (s *Session) Set(name string, value string) {
	s.upd = time.Now()
	if value == "" {
		delete(s.arg, name)
	} else {
		s.arg[name] = value
	}
	s.mgr.Lock()
	s.mgr.reg[s.ID] = *s
	s.mgr.Unlock()
}
