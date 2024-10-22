package main

import (
	// "compress/zlib"
	"bufio"
	"bytes"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"io"
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

// TEST: Using compress/zlib
func read_segment_zlib(data []byte, from, to int) ([]byte, error) {
	b := bytes.NewReader(data[from : from+to])
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := io.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// TEST: Using compress/flate
func read_segment_flate(data []byte, from, to int) ([]byte, error) {
	b := bytes.NewReader(data[from : from+to])
	z := flate.NewReader(b)
	defer z.Close()
	p, err := io.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func initiate_logger() log.Logger {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
	})

	return *logger
}

func process_streams(fileContent []byte, logger log.Logger) {
	streamMarker := []byte("stream")
	endstreamMarker := []byte("endstream")

	startIdx := 0
	for {
		startStream := bytes.Index(fileContent[startIdx:], streamMarker)
		if startStream == -1 {
			break
		}
		startStream += startIdx + len(streamMarker)

		for fileContent[startStream] == '\n' || fileContent[startStream] == '\r' || fileContent[startStream] == ' ' {
			startStream++
		}

		endStream := bytes.Index(fileContent[startStream:], endstreamMarker)
		if endStream == -1 {
			logger.Error("No matching 'endstream' found for 'stream'")
			break
		}
		endStream += startStream

		streamData := fileContent[startStream:endStream]

		content, err := read_segment_zlib(streamData, 0, len(streamData))
		if err != nil {
			logger.Warnf("Zlib decompression failed: %v, trying flate...", err)

			content, err = read_segment_flate(streamData, 0, len(streamData))
			if err != nil {
				logger.Error("Flate decompression also failed:", err)
				startIdx = endStream + len(endstreamMarker)
				continue
			}
		}

		fmt.Println("Decompressed content from stream:")
		fmt.Println(string(content))

		startIdx = endStream + len(endstreamMarker)
	}
}

func main() {
	log := initiate_logger()
	log.Info("Hello, World!")

	file_content, err := RetrieveROM("pdf_files/sample1.pdf")
	if err != nil {
		log.Error(err)
		return
	}

	process_streams(file_content, log)

	return
}
