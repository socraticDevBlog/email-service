# email-service

prototype for containerized k8s native email service

## routes

```bash
GET /sending

curl localhost:4000/sending
```

## test locally

### build Docker image, run it, and ping /sending endpoint

```bash
docker build . -t email-service:local

docker run -p 4000:4000 -d -t email-service:local

curl -i localhost:4000/sending
```

### expect

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: {now}
Content-Length: 109

{"id":1,"title":"Very Important Message!!!","email":"test@socratic.dev","content":"Lorem ipsum dipsum more"}
```
