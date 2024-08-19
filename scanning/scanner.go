// spectre-go/scanning/scanner.go

package scanning

// Scanner defines the interface for network scanning
type Scanner interface {
	Start() error
	Stop() error
}
