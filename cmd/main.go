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

func process_streams(fileContent []byte, logger log.Logger) string {
	streamMarker := []byte("stream")
	endstreamMarker := []byte("endstream")

	text_content := ""
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
			logger.Warnf("Zlib decompression failed: %d-%d:%v, trying flate...", startIdx, endStream, err)

			content, err = read_segment_flate(streamData, 0, len(streamData))
			if err != nil {
				logger.Errorf("Flate decompression also failed: %d-%d:%v", startIdx, endStream, err)
				startIdx = endStream + len(endstreamMarker)
				continue
			}
		}

		// fmt.Println("Decompressed content from stream:")
		// fmt.Println(string(content))

		content2 := string(content)
		for j := 0; j < len(content2); j++ {
			if content2[j] == '(' {
				j++
				for ; j < len(content2); j++ {
					if content2[j] == ')' {
						break
					}
					text_content += string(content2[j])
					// append to a dummy string
				}
			}
		}

		startIdx = endStream + len(endstreamMarker)
	}

	return text_content
}

func main() {
	log := initiate_logger()
	log.Info("Hello, World!")

	file_content, err := RetrieveROM("pdf_files/sample1.pdf")
	if err != nil {
		log.Error(err)
		return
	}

	text := process_streams(file_content, log)
	fmt.Println(text)

	return
}
