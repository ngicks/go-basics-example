#! /bin/sh

docker buildx build\
    --build-arg HTTP_PROXY=${HTTP_PROXY}\
    --build-arg HTTPS_PROXY=${HTTPS_PROXY}\
    --build-arg MAIN_PKG_PATH=${MAIN_PKG_PATH:-./}\
    --secret id=certs,src=/etc/ssl/certs/ca-certificates.crt\
    --secret id=.netrc,src=${DOTNETRC_PATH}\
    --secret id=goenv,src=$(go env GOENV)\
    -t $1\
    -f Dockerfile\
    .