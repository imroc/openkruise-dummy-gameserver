#!/bin/bash

result=$(curl http://127.0.0.1:8080/api/idle)

if [ "$result" == "true" ]; then
  exit 0
else
  exit 1
fi
