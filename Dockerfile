# --- STAGE 1: Compile GO ---
FROM golang:1.24-alpine AS builder

# Set the working directory for Go modules
WORKDIR /go/src/app/server

# Copy the Go module files and download dependencies from root project
COPY server/go.mod ./
RUN go mod download

# Copy the entire 'server' directory and build the application
COPY server/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o /go/bin/server .

# Start from Alpine base image
FROM alpine
LABEL maintainer="Mr.Philipp <d3vilh@github.com>"
LABEL version="0.5.5"

# Set the working directory to /opt/app
WORKDIR /opt/app

RUN apk --no-cache --no-progress upgrade && apk --no-cache --no-progress add bash bind-tools oath-toolkit-oathtool curl ip6tables iptables openvpn easy-rsa

#Install Latest RasyRSA Version
RUN chmod 755 /usr/share/easy-rsa/*

# Copy all files in the current directory to the /opt/app directory in the container
COPY bin /opt/app/bin
COPY docker-entrypoint.sh /opt/app/docker-entrypoint.sh

COPY fortivpn/forticonfig /opt/app/forticonfig
COPY fortivpn/openfortivpn /usr/bin/openfortivpn

COPY --from=builder /go/bin/server /opt/app/bin/

RUN chmod +x /usr/bin/openfortivpn
RUN mkdir -p /opt/app/clients \
    /opt/app/db \
    /opt/app/log \
    /opt/app/pki \
    /opt/app/staticclients \
    /opt/app/config

# Add the openssl-easyrsa.cnf file to the easy-rsa directory
ADD openssl-easyrsa.cnf /opt/app/easy-rsa/

# Make all files in the bin directory executable
RUN chmod +x bin/*; chmod +x docker-entrypoint.sh

# Expose the OpenVPN port (1194/udp)
EXPOSE 1194/udp

# Set the entrypoint to the docker-entrypoint.sh script, passing in the following arguments:
# $REQ_COUNTRY $REQ_PROVINCE $REQ_CITY $REQ_ORG $REQ_OU $REQ_CN
ENTRYPOINT ./docker-entrypoint.sh