package output

import (
	"encoding/csv"
	"os"
	"sync"
)

type CSVWriter struct {
    file   *os.File
    writer *csv.Writer
    mu     sync.Mutex // Protects against concurrent writes if you use goroutines later
}