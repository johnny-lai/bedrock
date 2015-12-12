# Based on https://blog.golang.org/docker

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM scratch

# Copy configuration
COPY ./config/production.yml /opt/etc/{{.Env "APP_NAME"}}.yml

# Copy executables
COPY ./tmp/dist/{{.Env "APP_NAME"}} /go/bin/{{.Env "APP_NAME"}}

# Entrypoint
ENTRYPOINT ["/go/bin/{{.Env "APP_NAME"}}", "-c", "/opt/etc/{{.Env "APP_NAME"}}.yml"]
CMD ["server"]

# Service listens on port 8080
EXPOSE 8080
