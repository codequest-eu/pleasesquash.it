#!/bin/sh

trap 'kill -9 $(jobs -p)' EXIT
echo -e $SSL_CERT > cert.pem
echo -e $SSL_KEY > key.pem
echo -e $DATASTORE_KEY > datastore.key
$@
