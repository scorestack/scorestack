#!/bin/sh
# Shoutout to jamielinux for publishing a guide on setting up a certificate
# authority with OpenSSL. The following script was written using the guide
# published at the link below.
# https://jamielinux.com/docs/openssl-certificate-authority/index.html

# Check if PKI directory already exists
if [ -d "pki" ]; then
    echo "pki directory already exists"
    echo "If you would like to regenerate the certificate chain, please delete the pki directory"
    exit 0
fi

# Set up directory structure
mkdir -p pki/intermediate
cp root-ca.conf pki/openssl.conf
cp intermediate-ca.conf pki/intermediate/openssl.conf
cd pki
# cwd: /app/pki
mkdir certs crl newcerts private
touch index.txt
echo 1000 > serial

# Create root key & cert
openssl ecparam -name prime256v1 -genkey -noout -out private/ca.key.pem
openssl req -config openssl.conf -key private/ca.key.pem -new -x509 -days 7300 -sha256 -extensions v3_ca -out certs/ca.cert.pem -subj '/C=US/ST=NewYork/O=ScoreStack/OU=ScoreStack/CN=ScoreStack Root CA'

# Set up intermediate directory structure
cd intermediate
# cwd: /app/pki/intermediate
mkdir certs crl csr newcerts private
touch index.txt
echo 1000 > serial
echo 1000 > crlnumber
cd ..
# cwd: /app/pki

# Create intermediate key, cert, and cert chain
openssl ecparam -name prime256v1 -genkey -noout -out intermediate/private/intermediate.key.pem
openssl req -config intermediate/openssl.conf -new -sha256 -key intermediate/private/intermediate.key.pem -out intermediate/csr/intermediate.csr.pem -subj '/C=US/ST=NewYork/O=ScoreStack/OU=ScoreStack/CN=ScoreStack Intermediate CA'
yes y | openssl ca -config openssl.conf -extensions v3_intermediate_ca -days 3650 -notext -md sha256 -in intermediate/csr/intermediate.csr.pem -out intermediate/certs/intermediate.cert.pem
openssl verify -CAfile certs/ca.cert.pem intermediate/certs/intermediate.cert.pem
cat intermediate/certs/intermediate.cert.pem certs/ca.cert.pem > intermediate/certs/ca-chain.cert.pem

# Create and sign certificates
for name in elas01 elas02 elas03 logs01 localhost
do
  openssl ecparam -genkey -noout -name prime256v1 \
    -out intermediate/private/$name.key.pem
  openssl req -config intermediate/openssl.conf -new -sha256 \
    -key intermediate/private/$name.key.pem \
    -out intermediate/csr/$name.csr.pem \
    -subj "/C=US/ST=NewYork/O=ScoreStack/OU=ScoreStack/CN=$name"
  yes y | openssl ca -notext -md sha256 \
    -config intermediate/openssl.conf \
    -in intermediate/csr/$name.csr.pem \
    -out intermediate/certs/$name.cert.pem
done
openssl verify -CAfile intermediate/certs/ca-chain.cert.pem \
    intermediate/certs/elas01.cert.pem \
    intermediate/certs/elas02.cert.pem \
    intermediate/certs/elas03.cert.pem \
    intermediate/certs/localhost.cert.pem

# Convert the logstash key
openssl pkcs8 -in intermediate/private/logs01.key.pem -topk8 -out intermediate/private/logs01.key.pkcs8