#!/bin/bash
for ((i=1;i<=1000;i++)); do
  curl -v --header "Connection: keep-alive" "localhost:2000/";
  curl -v --header "Connection: keep-alive" "localhost:3000/";
  curl -v --header "Connection: keep-alive" "localhost:4000/";
  curl -v --header "Connection: keep-alive" "localhost:5000/";
done
