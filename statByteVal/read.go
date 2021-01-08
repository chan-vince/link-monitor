package statByteVal

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func ReadFromFile(filePath string) uint64 {

	fmt.Println(filePath)


	// Open cmd
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Get the size from the FileInfo object
	fileInfo, err := fd.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	data := make([]byte, fileSize)

	count, err := fd.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes: %q\n", count, data[:count])

	final, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(final)

	fd.Close()
	return final
}

func HandleReadValue(value uint64) int {
	fmt.Printf("Value: %d\n", value)
	return 0
}