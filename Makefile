BINARY_NAME=secure-store
 
build:
	CGO_ENABLED=0 go build -o ${BINARY_NAME} -ldflags '-w -s' ./cmd/${BINARY_NAME}/

dist:
	CGO_ENABLED=0 go build -o dist/${BINARY_NAME} -ldflags '-w -s' ./cmd/${BINARY_NAME}/
 
run:
	go build -o ${BINARY_NAME} -ldflags '-w -s' ./cmd/${BINARY_NAME}/
	./${BINARY_NAME}

test:
	echo "No tests"
