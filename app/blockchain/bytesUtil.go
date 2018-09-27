package blockchain

import "time"

// Append bytes into byte slice with i size
func FitBytesInto(d []byte, size int) []byte {
	if len(d) < size {

		dif := size - len(d)

		return append(SliceOfBytes(dif, 0), d...)
	}

	// Return the slice with length equal to the given size
	return d[:size]
}

// Cretes a slice of bytes for the "i" size
func SliceOfBytes(i int, b byte) (p []byte) {
	for i != 0 {
		p = append(p, b)
		i--
	}
	return
}

func StripByte(d []byte, b byte) (p []byte) {
	for i, bb := range d {

		if bb != b {
			return d[i:]
		}
	}

	return nil
}

func Timeout(i time.Duration) chan bool {

	// Creates a new pointer type channel
	// to pauses the current goroutine
	t := make(chan bool)
	go func() {
		time.Sleep(i)
		t <- true
	}()

	// return the time
	return t
}
