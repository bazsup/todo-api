create_todo:
	restcli --show body ./test/create_todos.http

kill:
	lsof -i :8081 | awk 'NR > 1 {print $$2}' | xargs kill -SIGTERM

build:
	go build \
			-ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
			-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%SZ:00"`" \
			-o app

liveness:
	@cat /tmp/live 2> /dev/null
	@echo $$?

install_vegeta:
	go install github.com/tsenart/vegeta@latest

attack:
	echo "GET http://:8081/limitz" | vegeta attack -rate=10/s -duration=1s | vegeta report
