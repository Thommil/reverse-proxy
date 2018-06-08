#!/bin/sh

if [ $# -ne 1 ] ; then
	echo "Usage : docker-rebuild.sh all|CONTAINER_NAME"
	exit 1
fi

if [ $1 == 'all' ] ; then
	docker-compose build && docker-compose down && docker-compose up -d
else
	docker-compose build $1 && docker-compose rm -f $1 && docker-compose up -d $1
fi