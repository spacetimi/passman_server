#!/bin/bash
# Run the docker container for passman against $1 environment
# See also: dockerbuild.sh for building the docker image locally
if [ $1 = "Local" ] || [ $1 = "Test" ]; then
   docker run --restart always -d -v $HOME/.aws:/root/.aws -e app_environment=$1 -e app_name=passman -e app_dir_path=/go/src/github.com/spacetimi/passman_server --publish 9000:9000 passman-server | xargs -I containerId docker logs -f containerId
else
   echo "Usage: dockerrun.sh <Local|Test>"
fi

