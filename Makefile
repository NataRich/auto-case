SRC_DIR := $(abspath ./src)

all: case

case:
		cd ${SRC_DIR} && go build -o ../case main.go common.go config.go request.go

.PHONY: clean
clean:
		rm case 
