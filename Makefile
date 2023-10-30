BINARY_NAME = shurl

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

# (build but with a smaller binary)
.PHONY: dist
dist:
	CGO_ENABLED=0 go build -gcflags=all=-l -v -ldflags="-w -s" -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: run
run: build
	./$(BINARY_NAME) 

.PHONY: test
test: 
	$(TEST_COMMAND) -cover -parallel 5 -failfast -count=1 ./... 

# human readable test output
.PHONY: love
love:
ifeq ($(filter watch,$(MAKECMDGOALS)),watch)
	gotestsum --watch ./...
else
	gotestsum ./...
endif

.PHONY: tidy
tidy:
	go mod tidy

# auto restart
.PHONY: dev
dev:
	air
