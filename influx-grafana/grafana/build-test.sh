#! /bin/bash

pushd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
docker build -t vish/heapster_grafana:e2e_test .
docker push vish/heapster_grafana:e2e_test
popd
