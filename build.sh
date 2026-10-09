#!/bin/bash

set -e 

rm -rf output

mkdir output
mkdir output/frontend

cd frontend
npm install
npm run build
cd ..

cp -r ./docker ./output
cp -r ./frontend/dist ./output/frontend
cp ./start.sh ./output
cp ./cleanup.sh ./output
cp ./README.md ./output

go build -o ./output/vnc ./main.go
