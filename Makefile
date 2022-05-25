create_todo:
	restcli --show body ./test/create_todos.http

kill:
	lsof -i :8081 | awk 'NR > 1 {print $$2}' | xargs kill -SIGTERM
