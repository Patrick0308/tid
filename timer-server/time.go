package timer_server

import (
	"sync"
	"time"
)

var second uint32 = 0
var rwMutex = &sync.RWMutex{}

func GetSecond() uint32 {
	var tSecond uint32
	rwMutex.RLock()
	tSecond = second
	rwMutex.RUnlock()
	return tSecond
}
func StartTime() {
	go incrSecond()
}

func incrSecond()  {
	rwMutex.Lock()
	second = second + 1
	rwMutex.Unlock()
	time.AfterFunc(1*time.Second, incrSecond)
}

