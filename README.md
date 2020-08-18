# gowebapp
## Documentation TODO
Simple Golang REST Api using CQL
Api provides 3 endpoints:

``- POST /api/message``

Which enables user to store provided messages in Cassandra database. Messages older than 5 minutes are being automatically deleted. 
Usage example:

```
curl -X localhost:8080/api/message -d '{"email":"test@example.com","title":"testTitle","content":"testContent","magic_number":16}'
```

``- POST /api/send``

That sends emails with specified magic_number value, deleting them from database afterwords
Usage example:
```
curl -vv localhost:8080/api/send -d '{"magic_number":11}'
```

``- GET /api/messages/{emailValue}``

GET request that returns all messages with email value specified in URL. Usage example:
```
curl -X GET localhost:8080/api/messages/test@example.com
```
Will return first 4 email values. To get next batch of values, provide URL with page parameter like so:
```
curl -X GET localhost:8080/api/messages/test@example.com?page=2
```

## Deploy
### Docker
First run Cassandra in Docker Container with the detached mode:
```
docker run -p 9042:9042 -d --name cassandra cassandra
```
After Cassandra starts running, build container with an application:
```
docker build -t gowebapp .
```
After building the application run the application with:
```
docker run --env-file=appEnv.env -p 8080:8080 --link cassandra gowebapp /bin/app
```
Then after application starts You can run tests with:
```
docker run --link=cassandra gowebapp /bin/handlers_test
```
### Docker-compose
Another way to deploy whole project is to use docker-compose:
```
docker-compose up --build
```
