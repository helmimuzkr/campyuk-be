.PHONY: all
all: mysql_start run

.PHONY: mysql_start
mysql_start: ## run mysql service
	sudo service mysql start

.PHONY: test
test: ## running unit test
	go test ./... --coverprofile cover.out

.PHONY: cover
cover: ## check the coverage test
	go tool cover -func cover.out

.PHONY: testcover
testcover: test cover


.PHONY: run
run: ## run main.go
	go run main.go


