build:
	go build -o ./bin/dot-user-git-util
build-watch:
	while inotifywait -r -e modify,move,create,delete ./; do make build; done
