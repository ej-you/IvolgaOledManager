package i2c

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	host "periph.io/x/host/v3"
)

const (
	addrFrom = 0x08
	addrTo   = 0x77
)

var _ Scanner = (*scanner)(nil)

type Scanner interface {
	ScanAll() (*ScanResult, error)
}

type scanner struct {
	logger *log.Logger
}

func NewScanner(logOutput io.Writer) Scanner {
	return &scanner{
		logger: log.New(logOutput, "[INFO]", log.Ldate|log.Ltime),
	}
}

// ScanAll scans all busses and all addresses of each bus.
// It returns map of found buses with slices of available addresses for each bus.
func (s *scanner) ScanAll() (*ScanResult, error) {
	// initialise all relevant drivers
	if _, err := host.Init(); err != nil {
		s.logger.Println(err)
		return nil, err
	}

	// get all bus refs
	foundRefs := i2creg.All()
	foundRefsAmount := len(foundRefs)
	// init scan result struct
	scanResult := newScanResult(foundRefsAmount)

	var wgAll sync.WaitGroup
	wgAll.Add(foundRefsAmount)
	for _, foundRef := range foundRefs {
		go func() {
			defer wgAll.Done()
			s.scanBus(foundRef, scanResult)
		}()
	}
	wgAll.Wait()
	return scanResult, nil
}

// scanAddr scan all address of a given bus concurrently.
func (s *scanner) scanBus(ref *i2creg.Ref, scanResult *ScanResult) {
	s.logger.Printf("start check bus %s...\n", ref.Name)

	var wgBus sync.WaitGroup
	wgBus.Add(int(addrTo-addrFrom) + 1)

	var addr uint16
	for addr = 0x08; addr <= 0x77; addr++ {
		go func() {
			defer wgBus.Done()
			s.scanAddr(ref, addr, scanResult)
		}()
	}
	wgBus.Wait()
}

// scanAddr checks given address of a bus via attempt to send zero byte to device.
func (s *scanner) scanAddr(ref *i2creg.Ref, addr uint16, scanResult *ScanResult) {
	busNumber := ref.Number
	// open connection with bus
	bus, err := ref.Open()
	if err != nil {
		s.logger.Printf("bus %d: addr %x: open connection with bus error: %v\n", busNumber, addr, err)
		return
	}
	// open connection with device
	dev := i2c.Dev{Bus: bus, Addr: addr}
	if err := dev.Tx([]byte{0}, nil); err != nil {
		s.logger.Printf("bus %d: addr %x: open connection with device error: %v\n", busNumber, addr, err)
		return
	}
	s.logger.Printf("bus %d: addr %x: checked successfully\n", busNumber, addr)
	scanResult.add(ref.Name, addr)
}

// ScanResult store full i2c scan results.
// It implements Stringer interface and provides thread-safe method add for appending results.
type ScanResult struct {
	result map[string][]uint16
	mu     sync.Mutex
}

func newScanResult(resultLen int) *ScanResult {
	return &ScanResult{
		result: make(map[string][]uint16, resultLen),
	}
}

// add is a thread-safe method for appending results to result map.
func (r *ScanResult) add(key string, value uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// if slice for given key is not exists
	if _, ok := r.result[key]; !ok {
		r.result[key] = []uint16{value}
		return
	}
	r.result[key] = append(r.result[key], value)
}

// Implements Stringer interface.
func (r *ScanResult) String() string {
	builder := strings.Builder{}
	for key, values := range r.result {
		builder.WriteString("\n")
		builder.WriteString(key)
		builder.WriteString(":\t")
		for _, value := range values {
			builder.WriteString(fmt.Sprintf("%x", value))
			builder.WriteString(" ")
		}
	}
	return builder.String()
}
