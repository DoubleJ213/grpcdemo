#! /bin/bash

\rm ca.key
\rm ca.crt
\rm server.key
\rm server.csr
\rm server.crt

\rm client.key
\rm client.crt

openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-server/CN=localhost" \
    -key ca.key -out ca.crt

openssl genrsa -out server.key 2048
openssl req -new \
    -subj "/C=GB/L=China/O=server/CN=localhost" \
    -key server.key \
    -out server.csr

#openssl x509 -req -sha256 \
#    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
#    -in server.csr \
#    -out server.crt

openssl x509 -req -sha512 -days 3650 -extfile openssl.cnf \
        -CA ca.crt -CAkey ca.key -CAcreateserial \
        -in server.csr -out server.crt

openssl genrsa -out client.key 2048
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=client/CN=localhost" \
    -key client.key -out client.crt

