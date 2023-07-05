.PHONY: build

BINARY_NAME=sillyGirl

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin -ldflags "-s -w" main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux -ldflags "-s -w" main.go


release:
	GOARCH=amd64 GOOS=darwin garble build -o ${BINARY_NAME}-darwin -a -ldflags "-s -w" -trimpath -buildmode=pie main.go
	GOARCH=amd64 GOOS=linux garble build -o ${BINARY_NAME}-linux -a -ldflags "-s -w" -trimpath -buildmode=pie main.go

pre-release:
	go install mvdan.cc/garble@latest
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux