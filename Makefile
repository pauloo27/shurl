BINARY_NAME = shurl
TEST_COMMAND = go test -v

.PHONY: all
all: swagger build

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
	gotestsum --watch -- -cover ./...
else
	gotestsum -- -cover ./...
endif

.PHONY: tidy
tidy:
	go mod tidy

# auto restart
.PHONY: dev
dev:
	air

.PHONY: lint
lint:
	revive -formatter friendly -config revive.toml ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: gosec
gosec:
	gosec -tests ./... 

.PHONY: swagger
swagger:
	swag fmt
	swag init --dir ./internal/server/,./internal/models/ --output ./internal/server/docs --outputTypes go,yaml

.PHONY: inspect
inspect: lint gosec staticcheck

.PHONY: install-inspect-tools
install-inspect-tools:
	go install github.com/mgechev/revive@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: install-swaggo
install-swaggo:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: install-dev-tools
install-dev-tools:
	go install github.com/cosmtrek/air@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/swaggo/swag/cmd/swag@latest

