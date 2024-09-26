package pkg

import (
	"log"
	"os"
	"syscall"
)

func LoadFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}

	return f
}

func CloseFile(f *os.File) {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	f.Close()
}
