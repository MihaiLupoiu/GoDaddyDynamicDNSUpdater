FROM golang:alpine as builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/MihaiLupoiu/GoDaddyDynamicDNSUpdater
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /GoDaddyDynamicDNSUpdater .

FROM alpine:latest
RUN apk update && apk upgrade
RUN apk add --no-cache ca-certificates netcat-openbsd

COPY --from=builder /GoDaddyDynamicDNSUpdater ./
RUN chmod +x /GoDaddyDynamicDNSUpdater
    
ENTRYPOINT ["./GoDaddyDynamicDNSUpdater"]