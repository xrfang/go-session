package session

import (
	"encoding/json"
	"io"
	"reflect"
	"time"
)

type Session struct {
	ID  string
	src string
	upd time.Time
	mgr *Manager
	arg map[string]interface{}
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

func (s Session) Get(name string, value interface{}) bool {
	v := s.arg[name]
	if v == nil {
		return false
	}
	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v))
	if s.mgr.cfg.Refresh {
		s.upd = time.Now()
		s.mgr.Lock()
		s.mgr.reg[s.ID] = s
		s.mgr.Unlock()
	}
	return true
}

func (s *Session) Del(name string) {
	s.upd = time.Now()
	delete(s.arg, name)
	s.mgr.Lock()
	s.mgr.reg[s.ID] = *s
	s.mgr.Unlock()
}

func (s *Session) Set(name string, value interface{}) {
	s.upd = time.Now()
	s.arg[name] = value
	s.mgr.Lock()
	s.mgr.reg[s.ID] = *s
	s.mgr.Unlock()
}
