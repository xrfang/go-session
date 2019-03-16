package session

import (
	"fmt"
	"os"
	"path"
)

type Config struct {
	TTL     int    //session life time in seconds
	VoidTTL int    //life time for empty session in seconds
	Refresh bool   //reset session life time on Get
	Persist string //path to save session data
}

func (sm *Manager) setConfig(cfg *Config) {
	if cfg == nil {
		sm.cfg.Refresh = true
	} else {
		sm.cfg = *cfg
	}
	if sm.cfg.TTL <= 0 {
		sm.cfg.TTL = 86400
	}
	if sm.cfg.VoidTTL <= 0 {
		sm.cfg.VoidTTL = 60
	}
	if sm.cfg.Persist != "" {
		err := os.MkdirAll(path.Dir(sm.cfg.Persist), 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, "session.Manager.setConfig:", err)
		}
	}
}
