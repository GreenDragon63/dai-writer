package models

import (
	"log"
	"sync"
	"time"
)

type userMutexMap struct {
	global    sync.Mutex
	lastCheck time.Time
	users     map[int]*sync.Mutex
	lastUsed  map[int]time.Time
}

var um = userMutexMap{
	lastCheck: time.Now(),
	users:     make(map[int]*sync.Mutex),
	lastUsed:  make(map[int]time.Time),
}

func (um *userMutexMap) lock(uid int) {
	um.global.Lock()
	defer um.global.Unlock()

	if um.users[uid] == nil {
		um.users[uid] = &sync.Mutex{}
	}

	um.users[uid].Lock()
}

func (um *userMutexMap) unlock(uid int) {
	um.global.Lock()
	defer um.global.Unlock()

	if um.users[uid] != nil {
		um.users[uid].Unlock()
		um.lastUsed[uid] = time.Now()
	} else {
		log.Printf("No mutex to unlock for user: %d\n", uid)
	}

	if time.Since(um.lastCheck) >= 1*time.Hour {
		for id, lastUsed := range um.lastUsed {
			if time.Since(lastUsed) >= 1*time.Hour {
				delete(um.users, id)
				delete(um.lastUsed, id)
			}
		}
	}
}
