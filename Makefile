app:
	PORT=8080 DATABASE_URL="mongodb://localhost:27017/" go run main.go
app-seed:
	PORT=8080 DATABASE_URL="mongodb://localhost:27017/" go run main.go

ab:
	ab -n 100 -c 10 -g out.data http://localhost:8080/users/all > ab.txt