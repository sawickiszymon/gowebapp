# gowebapp
## Dokumentacja TODO

## Deploy
### Docker
First run Cassandra in Docker Container with detached mode:
```
docker run -p 9042:9042 -d --name cassandra cassandra
```
After Cassandra starts running, build container with an application:
```
docker build -t gowebapp .
```
After building the application run the application with:
```
docker run --env-file=appEnv.env -p 8080:8080 --link cassandra gowebapp
```

### Docker-compose
Another way to deploy whole project is to use docker-compose:
```
docker-compose up
```
