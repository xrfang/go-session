package session

import (
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

func (s *Session) ttl() int {
	ttl := s.mgr.cfg.TTL
	if len(s.arg) == 0 {
		ttl = s.mgr.cfg.VoidTTL
	}
	return int(time.Now().Sub(s.upd).Seconds()) - ttl
}

func (s *Session) Get(name string, value interface{}) bool {
	v := s.arg[name]
	if v == nil {
		return false
	}
	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v))
	if s.mgr.cfg.Refresh {
		s.upd = time.Now()
	}
	return true
}

func (s *Session) Set(name string, value interface{}) {
	s.upd = time.Now()
	s.arg[name] = value
}
