# Container image that runs your code
FROM golang:1.16
WORKDIR /go/src/asana_github_action
COPY . .

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]
