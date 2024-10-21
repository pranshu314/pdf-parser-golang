.PHONY: run
run:
	go run ./cmd/

.PHONY: build
build:
	make clean; go build -o ./out/parse-pdf ./cmd/ && ./out/parse-pdf

.PHONY: clean
clean:
	rm ./out/parse-pdf
