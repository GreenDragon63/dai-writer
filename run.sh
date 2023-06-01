#!/bin/bash

if [ ! -f .env ]; then
	cp env.default .env
fi

. .env
./dai-writer

