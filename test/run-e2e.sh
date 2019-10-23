#!/usr/bin/env bash

export S3_ENDPOINT=https://oss2-gw.example.org

# TODO Replace access/secret keys of user1 and user2
export U1_ACCESS_KEY='pion-1234'
export U1_SECRET_KEY='....'

export U2_ACCESS_KEY='pion-5678'
export U2_SECRET_KEY='....'

go test -v ./e2e/