#!/bin/bash
domain=$1
openssl req -x509 -sha256 -nodes -newkey rsa:2048 \
    -keyout $domain.key -out $domain.crt \
    -days 365 \
    -subj "/CN=*.$domain" \
    -addext "subjectAltName=DNS:*.$domain"
