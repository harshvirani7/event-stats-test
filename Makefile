# Define variables
IMAGE_NAME=event_stats_test
CONTAINER_NAME=event_stats_test_container

build:
	docker build -t $(IMAGE_NAME) .

# Run Docker container
run:
	docker run -d -p 8080:8080 --name $(CONTAINER_NAME) $(IMAGE_NAME)

# Stop Docker container
stop:
	docker stop $(CONTAINER_NAME)

# Remove Docker container
rm:
	docker rm $(CONTAINER_NAME)

# Remove Docker image
rmi:
	docker rmi $(IMAGE_NAME)
