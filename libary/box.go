/**
 * box.go
 *
 * @author : Cyw
 * @email  : rose2099.c@gmail.com
 * @created: 2017/6/12 下午6:32
 * @logs   :
 *
 */
package libary

import (
	"sort"
	"sync"
	"log"
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
	Run() (Rest, error)
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

func Open(name string) (rest Rest, e error) {
	if _, dup := plugs[name]; dup {
		// ok
		plug := plugs[name]
		return plug.Run()
	} else {
		log.Print("not found site")
		return rest, e
	}
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
