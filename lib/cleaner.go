package lib

import (
	"runtime"
	"runtime/debug"
	"time"
	//"log"
	//"net/http"
	//_ "net/http/pprof"
)

// CleanUp function call GC and freeosmem
func CleanUp() {
	// go func(){
	// 	log.Println(http.ListenAndServe(":6060", nil))

	// }()
	for {
		runtime.GC()
		time.Sleep(10 * time.Second)
		debug.FreeOSMemory()
		//var m runtime.MemStats
		//runtime.ReadMemStats(&m)
		//log.Printf("HeapSys: %d, HeapAlloc: %d, HeapIdle: %d, HeapReleased: %d\n", m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased)
		//log.Println("NumGoroutine:",runtime.NumGoroutine())
		time.Sleep(1 * time.Minute)
	}
}
