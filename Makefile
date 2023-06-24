SRC_DIR := $(abspath ./src)

all: case

case:
		cd ${SRC_DIR} && go build -o ../case main.go common.go config.go request.go

case-windows:
		cd ${SRC_DIR} && GOOS=windows go build -o ../case.exe main.go common.go config.go request.go

.PHONY: clean
clean:
		rm case case.exe
