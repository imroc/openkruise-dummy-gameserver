#!/bin/bash

# 发送请求并获取状态码
status_code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/api/idle)

# 检查状态码是否为200
if [ "$status_code" -eq 200 ]; then
  exit 0
else
  exit 1
fi
