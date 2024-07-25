IMAGE_NAME=my-go-app

MAIN_FILE=cmd/main.go

PORT=8080

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run -p $(PORT):$(PORT) $(IMAGE_NAME)

go-run:
	go run $(MAIN_FILE)

docker-stop:
	-docker stop 99b21cdbaa54

docker-rm:
	-docker rm $(shell docker ps -a -q --filter ancestor=$(IMAGE_NAME))

docker-clean: docker-stop docker-rm
	-docker rmi -f $(IMAGE_NAME)
	-docker system prune -f

clean:
	go clean

build:
	go build -o main $(MAIN_FILE)

all: docker-build docker-run

test-request:
	curl -X GET http://localhost:$(PORT)/most-changed-address
