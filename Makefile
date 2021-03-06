configure:
	go get -u mvdan.cc/xurls github.com/go-redis/redis
	go get -u github.com/stretchr/testify/assert
	go get -u github.com/justincampbell/timeago
	go get -u github.com/sirupsen/logrus
	go get -u golang.org/x/tools/cmd/cover
	go get -u golang.org/x/tools/cmd/goimports

build:
	go build -o bot.bin -ldflags="-s -w" cmd/novitoll_daemon_bot/main.go

run:
	@make build
	./bot.bin

docker-compose-local:
	docker-compose -f deployments/docker-compose-local.yml up

docker-compose:
	docker-compose -f deployments/docker-compose.yml up -d

docker-compose-green:
	docker-compose -f deployments/docker-compose-green.yml up

docker-compose-stop:
	docker-compose -f deployments/docker-compose.yml stop

goimports:
	goimports -w `find . -type f -name '*.go' -not -path "./vendor/*"`

debug:
	dlv debug --output bot.bin cmd/novitoll_daemon_bot/main.go

test:
	go test -v -cover -coverprofile=coverage.out github.com/novitoll/novitoll_daemon_bot/internal/bot
	go tool cover -html=coverage.out -o coverage.html 
