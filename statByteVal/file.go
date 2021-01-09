package statByteVal

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func isFile(filePath string) (bool, string) {
	info, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		errStr := fmt.Sprintf("Does not exist: %s\n\n", filePath)
		return false, errStr
	}

	if info.IsDir() {
		errStr := fmt.Sprintf("Not a valid file: %s\n\n", filePath)
		return false, errStr
	}

	return true, ""
}

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
