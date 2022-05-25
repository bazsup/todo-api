create_todo:
	restcli --show body ./test/create_todos.http

kill:
	lsof -i :8081 | awk 'NR > 1 {print $$2}' | xargs kill -SIGTERM

build:
	go build \
			-ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
			-X main.buildtime=`date "+%Y-%m-%dT%H:%M:%SZ:00"`" \
			-o app