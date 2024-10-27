#!/bin/bash

{
  echo "building proxy ontainer"
  docker-compose -f src/doc/proxy/docker-compose.yml -p vim build
}