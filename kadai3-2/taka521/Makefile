NAME := pdown

all: build

build: clean
	go build -o ${NAME}.exe ./cmd/main.go

clean:
	@rm -f ./${NAME}.exe
