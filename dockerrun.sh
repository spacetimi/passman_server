#!/bin/bash
# Run the docker container for passman against $1 environment
# See also: dockerbuild.sh for building the docker image locally
if [ $1 = "Local" ]; then
   docker run --restart always -d -v $HOME/.aws:/root/.aws -e app_environment=Local -e app_name=passman --publish 9000:9000 passman-server | xargs -I containerId docker logs -f containerId
elif [ $1 = "Production" ]; then
   docker run --restart always -d -v $HOME/.aws:/root/.aws -e app_environment=Production -e app_name=passman -e aws_profile=passman-dev --publish 80:80 passman-server | xargs -I containerId docker logs -f containerId
else
   echo "Usage: dockerrun.sh <Local|Production>"
fi

