package main

import (
	// "compress/zlib"
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func RetrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}

func initiate_logger() log.Logger {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
	})

	return *logger
}

func main() {
	log := initiate_logger()
	log.Info("Hello, World!")

	file_content, err := RetrieveROM("pdf_files/sample1.pdf")
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println()
	fmt.Println("File Contents Below")
	fmt.Println(string(file_content))

	return
}
