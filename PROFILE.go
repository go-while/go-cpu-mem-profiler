package prof

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	hpprof "net/http/pprof"
	"runtime/pprof"
	"time"
	"sync"
	"os/signal"
	"syscall"
)

type Profiler struct {
	mux sync.Mutex
	mem sync.Mutex
	cpu sync.Mutex
	CPUProfile bool // set before boot
	MEMProfile bool // set before boot
	CpuProfileFile *os.File
	MemProfileFile *os.File
} // end Profiler struct

func NewProf() *Profiler {
	return &Profiler{
		CPUProfile: true,
		MEMProfile: true,
	}
}

func (p *Profiler) PprofWeb(addr string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	router := mux.NewRouter()
	router.Handle("/debug/pprof/", http.HandlerFunc(hpprof.Index))
	router.Handle("/debug/pprof/cmdline", http.HandlerFunc(hpprof.Cmdline))
	router.Handle("/debug/pprof/profile", http.HandlerFunc(hpprof.Profile))
	router.Handle("/debug/pprof/symbol", http.HandlerFunc(hpprof.Symbol))
	router.Handle("/debug/pprof/trace", http.HandlerFunc(hpprof.Trace))
	router.Handle("/debug/pprof/{cmd}", http.HandlerFunc(hpprof.Index)) // special handling for Gorilla mux
	server := &http.Server{Addr: addr, Handler: router}

	// go launch debug http
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// handle err
			log.Printf("debug_pprof ERROR server.ListenAndServe err='%v'", err)
		}
	}()
} // end func PprofWeb

func (p *Profiler) StartCPUProfile() (*os.File, error) {
	p.cpu.Lock()
	defer p.cpu.Unlock()
	if !p.CPUProfile {
		return nil, fmt.Errorf("ERROR !p.CPUProfile")
	}
	if p.CpuProfileFile != nil {
		return nil, fmt.Errorf("ERROR StartCPUProfile p.CpuProfileFile != nil")
	}
	fn := fmt.Sprintf("cpu.pprof.%d.out", time.Now().Unix())
	cpuProfileFile, err := os.Create(fn)
	if err != nil {
		log.Printf("ERROR startCPUProfile err1='%v'", err)
		return nil, err
	}
	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		log.Printf("ERROR startCPUProfile err2='%v'", err)
		return nil,err
	}
	log.Printf("startCPUProfile: fn='%s'", fn)
	p.CpuProfileFile = cpuProfileFile
	return cpuProfileFile, nil
} // end func startCPUProfile

func (p *Profiler) StopCPUProfile() {
	p.cpu.Lock()
	defer p.cpu.Unlock()
	if p.CpuProfileFile == nil {
		log.Printf("ERROR stopCPUProfile cpuProfileFile=nil")
		return
	}
	log.Printf("stopCPUProfile")
	pprof.StopCPUProfile()
	p.CpuProfileFile.Close()
	p.CpuProfileFile = nil
} // end StopCPUProfile



func (p *Profiler) StartMemoryProfile(duration time.Duration, wait time.Duration) error {
	p.mem.Lock()
	defer p.mem.Unlock()
	if !p.MEMProfile {
		log.Printf("ERROR StartMemoryProfile !p.MEMProfile")
		return fmt.Errorf("ERROR !p.MEMProfile")
	}
	if p.MemProfileFile != nil {
		log.Printf("ERROR StartMemoryProfile p.MemProfileFile != nil")
		return fmt.Errorf("ERROR p.MemProfileFile != nil")
	}
	go p.memoryProfiler(duration, wait)
	return nil
}

func (p *Profiler) memoryProfiler(duration time.Duration, wait time.Duration) {
	p.mem.Lock()
	defer p.mem.Unlock()
	if !p.MEMProfile {
		log.Printf("ERROR memoryProfiler !p.MEMProfile")
		return
	}
	if p.MemProfileFile != nil {
		log.Printf("ERROR memoryProfiler p.MemProfileFile != nil")
		return
	}
	// Generate a unique filename with a timestamp
	fn := fmt.Sprintf("mem.pprof.%d.out", time.Now().Unix())
	time.Sleep(wait)
	log.Printf("capture MemoryProfile duration=(%#v ns) waited=(%#v ns) fn='%s'", duration, wait, fn)
	// Create the profile file
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("ERROR memoryProfiler os.Create(fn='%s') err='%v'", fn, err)
		return
	}
	p.MemProfileFile = f
	// Start memory profiling
	pprof.Lookup("heap").WriteTo(f, 0)
	// Sleep for the specified duration to capture the memory profile
	time.Sleep(duration)
	f.Close()
	p.MemProfileFile = nil
	log.Printf("close MemoryProfile duration=(%#v ns) waited=(%#v ns) fn='%s'", duration, wait, fn)
	return
} // end func MemoryProfile

func (p *Profiler) CatchInterruptSignal(cpu bool, mem bool) {
	if cpu {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGINT)
			<-c
			p.StopCPUProfile()
		}()
	}
	/*
	if mem {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGINT)
			<-c
			p.stopMemProfile(p.MemProfileFile)
		}()
	}
	*/
} // end func CatchInterruptSignal
