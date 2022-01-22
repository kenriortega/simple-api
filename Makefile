git_hash := $(shell git rev-parse --short HEAD)

app:
	PORT=8080 DATABASE_URL="mongodb://localhost:27017/challenge" go run main.go

ab:
	ab -n 100 -c 10 -g out.data http://localhost:3000/all > ab.txt

docker-mongo-seed:
	docker exec -i mongo sh -c 'mongoimport -c users -d challenge --drop' < ./dataset/users.json

# Use the dockerfile principal and create an image with ~25mb
docker-build:
	docker build -t kenriortega/challenge-go:${git_hash} \
		-f Dockerfile .
docker-api:
	docker run --name importer --rm -m 128m --cpus="0.25" -p 3000:3000 --network host \
	 -e DATABASE_URL="mongodb://localhost:27017/challenge" \
	 -e PORT=3000 kenriortega/challenge-go:${git_hash}