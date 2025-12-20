#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING='user=postgres password=TestPass!1234 host=35.188.83.64 port=5432 dbname=notely sslmode=require'

goose up
