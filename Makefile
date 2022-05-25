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

maria:
	docker run -p 127.0.0.1:3306:3306 --name some-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:10.7.3

image:
	docker build -t todo:test -f Dockerfile .

container:
	docker run -p:8081:8081 --env-file ./test.env \
	--link some-mariadb:db --name runningtodo todo:test
