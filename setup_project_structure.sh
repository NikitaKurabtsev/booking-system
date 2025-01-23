#!/bin/bash

# Create project directories
mkdir -p cmd/{api,worker}
mkdir -p internal/{handlers,services,repositories,models,notifications,auth,config}
mkdir -p migrations
mkdir -p pkg/{cache,db,mq}
mkdir -p tests/{unit,integration,load}

# Create configuration files
touch Dockerfile
touch docker-compose.yml
touch README.md

echo "Project structure created successfully"
