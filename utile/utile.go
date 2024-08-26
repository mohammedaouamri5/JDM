package utile

import (
	"errors"
	fmt "fmt"
	os "os"
	"runtime"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"

	"time"
)

func PathIsExist(p_path string) (bool, error) {

	log.Infoln(fmt.Sprintf("\n\tTesting If the File %s exist ", p_path))
	_, err := os.Stat(p_path)
	if err == nil {
		log.Infoln("\n\tFile does not exist")
		return true, nil

	} else if os.IsNotExist(err) {
		log.Infoln("\n\tFile does not exist")
		return false, nil
	}

	return false, err
}

func SplitSlice(slice []int, numChunks int) ([][]int, error) {

	if len(slice) == 0 {
		err := fmt.Sprintln("The slice is empty")
		log.Error(err)
		return [][]int{{}}, errors.New(err)
	}

	var chunks [][]int
	chunkSize := (len(slice) + numChunks - 1) / numChunks // Ceiling of len(slice) / numChunks

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks, nil
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func RandomIntByRange(TheMin int, TheMax int) int {
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate a random number within the specified range
	return int(rand.Intn((TheMax - TheMin + 1))) + TheMin
}

func PrintCallStack() {
	// Create a slice to hold the program counters
	var pcs [50]uintptr

	// Capture the stack frames
	n := runtime.Callers(2, pcs[:]) // Skip 2 frames: the call to Callers itself and the function that called it

	// Iterate over the captured frames
	for _, pc := range pcs[:n] {
		f := runtime.FuncForPC(pc)
		file, line := f.FileLine(pc)
		fmt.Printf("%s:%d\n", file, line)
	}
}
