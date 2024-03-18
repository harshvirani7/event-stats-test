# event-stats-test
 An Api Server in Test cluster which can provide event statistics by camera.

## Instructions
    - sudo docker-compose up: Build and run the service using Docker Compose.
    - sudo docker-compose down: Stop the running service managed by Docker Compose.

## Testing
	- hey -n 100 -c 1 -q 2 http://localhost:8080/eventStats/totalEventCountByEventType?eventType=eventType2
	- hey -n 100 -c 1 -q 2 http://localhost:8080/eventStats/totalEventCountByCameraId?cameraId=camera1

