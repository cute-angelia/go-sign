package libary

import (
	"sort"
	"sync"
)

var (
	plugMu sync.RWMutex
	plugs  = make(map[string]Plug)
)

type Rest struct {
	Code    int
	Message string
}

type Plug interface {
	Run(name string) (Rest, error)
}

func Register(name string, plug Plug) {
	plugMu.Lock()
	defer plugMu.Unlock()
	if plug == nil {
		panic("box: Register plug is nil")
	}
	if _, dup := plugs[name]; dup {
		panic("box: Register called twice for plug " + name)
	}
	plugs[name] = plug
}

func unregisterAllDrivers() {
	plugMu.Lock()
	defer plugMu.Unlock()
	// For tests.
	plugs = make(map[string]Plug)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	plugMu.RLock()
	defer plugMu.RUnlock()
	var list []string
	for name := range plugs {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
