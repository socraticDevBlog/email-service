# email-service

prototype for containerized k8s native email service

## routes

```bash
GET /sending
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

## logging

basic logging has been implemented with different severity levels:

- ERROR: server is erroring, possibly down
- WARNING: a http request was made to a non-existing endpoint
- INFO: a http request to a valid endpoint
- DEBUG: (future use) for tracing

### logging - Docker

```bash
# print current email-service container ID
docker ps

docker logs <CONTAINER ID>
```

## minikube - running `email-service` on local kubernetes cluster

I like to use Minikube to develop k8s native service on local machine:
<https://minikube.sigs.k8s.io/docs/start/>

using Minikube requires some extra configs to make your container image
available to k8s runtime- see local file ([minikube-setup.sh](minikube-setup.sh)):

0. start Minikube
1. set Docker environment to be running within Minikube
2. build your container image within Minikube

```bash
kubectl apply -f deployment.yml

kubectl get pods

# expect
# > NAME                                       READY   STATUS    RESTARTS   AGE
#   email-service-deployment-d9b4c95d9-h9bh9   1/1     Running   0          8s

```
