package session

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (s *Session) pack() map[string]interface{} {
	return map[string]interface{}{
		"id":  s.ID,
		"src": s.src,
		"upd": s.upd,
		"arg": s.arg,
	}
}

func (s *Session) unpack(p map[string]interface{}) {
	defer func() {
		if e := recover(); e != nil {
			s.ID = ""
			fmt.Fprintln(os.Stderr, "session.Session.unpack:", e)
		}
	}()
	var ok bool
	chkType := func(tag string) {
		if !ok {
			panic(fmt.Errorf("invalid type for '%s'", tag))
		}
	}
	s.ID, ok = p["id"].(string)
	chkType("id")
	s.src, ok = p["src"].(string)
	chkType("src")
	upd, ok := p["upd"].(string)
	chkType("upd")
	if ok {
		t, err := time.Parse(time.RFC3339Nano, upd)
		if err != nil {
			panic(fmt.Errorf("invalid 'upd' (%v)", err))
		} else {
			s.upd = t
		}
	}
	s.arg, ok = p["arg"].(map[string]interface{})
	chkType("arg")
}

func (sm *Manager) loadSessions() {
	sm.Lock()
	defer func() {
		sm.Unlock()
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stderr, "session.Manager.loadSession:", e)
		}
	}()
	f, err := os.Open(sm.cfg.Persist)
	assert(err)
	defer f.Close()
	var data []map[string]interface{}
	assert(json.NewDecoder(f).Decode(&data))
	for _, d := range data {
		var s Session
		s.unpack(d)
		if s.ID != "" {
			s.mgr = sm
			sm.reg[s.ID] = &s
		}
	}
}

func (sm *Manager) saveSessions() {
	sm.RLock()
	defer func() {
		sm.RUnlock()
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stderr, "session.Manager.saveSession:", e)
		}
	}()
	f, err := os.Create(sm.cfg.Persist)
	assert(err)
	defer func() {
		err := f.Close()
		if e := recover(); e != nil {
			panic(e)
		}
		assert(err)
	}()
	var data []map[string]interface{}
	for _, s := range sm.reg {
		data = append(data, s.pack())
	}
	je := json.NewEncoder(f)
	je.SetIndent("", "    ")
	assert(je.Encode(data))
}
