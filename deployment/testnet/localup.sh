#!/bin/bash
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
project_path=$(git rev-parse --show-toplevel)
bin_name=bdjuno
bin=${project_path}/build/${bin_name}

function stop() {
    docker compose -f "${script_dir}"/docker-compose.yaml down
    rm -rf "${script_dir}"/data
    rm -rf "${script_dir}"/bdjuno.log
    ps -ef | grep "${bin_name}" | awk '{print $2}' | xargs kill
}

function start() {
    export HASURA_GRAPHQL_ADMIN_SECRET="${HASURA_GRAPHQL_ADMIN_SECRET:-moca}"
    docker compose -f "${script_dir}"/docker-compose.yaml up -d
    echo "wait 30s for graphql engine start..."
    for ((i = 30; i > 0; i -= 3)); do
        echo "please wait ${i}s..."
        sleep 3
    done

    echo "Importing the Hasura metadata..."
    hasura metadata apply --project "${project_path}"/hasura --endpoint http://localhost:8080 --admin-secret "${HASURA_GRAPHQL_ADMIN_SECRET}"

    # echo "Initializing the configuration..."
    # ${bin} --home "${script_dir}" parse genesis-file --genesis-file-path ./genesis.json

    # echo "run BDjuno...."
    # nohup "${bin}" start --home "${script_dir}" >"${script_dir}"/bdjuno.log 2>&1 &
    # "${bin}" --home "${script_dir}" start >"${script_dir}"/bdjuno.log
}

cmd=$1
case ${cmd} in
init)
    echo "===== init ===="

    echo "===== end ===="
    ;;
start)
    echo "===== start ===="
    start
    echo "===== end ===="
    ;;
stop)
    echo "===== stop ===="
    stop
    echo "===== end ===="
    ;;
reset)
    echo "===== reset ===="

    echo "===== stop ===="
    stop
    echo "===== end ===="

    echo "===== start ===="
    start
    echo "===== end ===="

    echo "===== end ===="
    ;;
*)
    echo "Usage: localup.sh init | reset | start | stop"
    ;;
esac
