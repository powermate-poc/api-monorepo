#!/bin/bash
paths=$1
for str in ${paths[@]}; do

  IFS=' ' read -ra ADDR <<< "$str"
  APP_NAME="${ADDR[0]}"

  echo "Building >${APP_NAME}<..."

  go build -o $APP_NAME lambda/$APP_NAME/main.go
  zip -r function-$APP_NAME.zip $APP_NAME
done