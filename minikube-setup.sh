minikube start

# sets Docker to be running within local Minikube instance
eval $(minikube docker-env)

# build Docker image within Minikube
docker build -t email-service:latest .
