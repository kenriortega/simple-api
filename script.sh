curl -X POST -H "Content-Type: application/json" \
    -d '{"first_name": "linuxize","last_name": "linuxize","phone": "53555555","job_title": "DEV", "email": "linuxize@example.com"}' \
    localhost:8080/users/create

curl localhost:8080/users/all