# SIMPLE API

Simple golang app using std lib net/http package to serve two endpoints REST API

### Endpoints GetAllUsers

```bash
curl localhost:8080/all
```

### Bench

```bash

# simple command two use 
# The values can be change. 
ab -n 1000 -c 100 -g out.data <http://localhost:3000/all> > ab.txt
```