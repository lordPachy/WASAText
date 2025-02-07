#!/bin/bash

./open-node.sh
# (here you're inside the container)
yarn run build-embed
exit
# (outside the container)
go build -tags webui ./cmd/webapi/
./open-node.sh
# (here you're inside the container)
yarn run dev