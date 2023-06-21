all: case

case:
		go build -o case main.go common.go config.go request.go

.PHONY: clean
clean:
		rm case 
