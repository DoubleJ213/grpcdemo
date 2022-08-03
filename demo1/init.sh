#! /bin/bash


\rm server.crt
\rm server.key

openssl genrsa -out server.key 2048

openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-server/CN=localhost" \
    -key server.key -out server.crt
