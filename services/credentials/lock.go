package credentials

import "sync"

var userTokenLocks sync.Map

func tokenLockForUser(userID uint) *sync.Mutex {
	stored, _ := userTokenLocks.LoadOrStore(userID, &sync.Mutex{})
	return stored.(*sync.Mutex)
}
