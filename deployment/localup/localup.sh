#!/bin/bash
basedir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
project_path=$(git rev-parse --show-toplevel)
bin_name=bdjuno
bin=${project_path}/build/${bin_name}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

function stop() {
    print_info "Stopping services..."
    docker compose -f "${basedir}"/docker-compose.yaml down
    if [ $? -ne 0 ]; then
        print_error "Failed to stop Docker services"
        exit 1
    fi

    # Note: We do NOT remove data directory or logs here to support restart logic.
    # Use 'reset' command if you want to clear data.

    # Kill bdjuno processes if they exist
    pids=$(ps -ef | grep "${bin_name}" | grep -v grep | awk '{print $2}')
    if [ -n "$pids" ]; then
        echo "$pids" | xargs kill -9 2>/dev/null || true
        print_info "BDjuno processes stopped"
    fi

    print_info "Services stopped"
}

function start() {
    # Check if docker is running
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running, please start Docker first"
        exit 1
    fi

    export HASURA_GRAPHQL_ADMIN_SECRET="${HASURA_GRAPHQL_ADMIN_SECRET:-moca}"
    print_info "Starting Docker services..."
    docker compose -f "${basedir}"/docker-compose.yaml up -d
    if [ $? -ne 0 ]; then
        print_error "Failed to start Docker services"
        exit 1
    fi

    print_info "Waiting for GraphQL Engine to start (30 seconds)..."
    for ((i = 10; i > 0; i -= 3)); do
        print_info "Please wait ${i} seconds..."
        sleep 3
    done

    print_info "Importing Hasura metadata..."
    hasura metadata apply --project "${project_path}"/hasura --endpoint http://localhost:8080 --admin-secret "${HASURA_GRAPHQL_ADMIN_SECRET}"
    if [ $? -ne 0 ]; then
        print_warning "Failed to import Hasura metadata, but continuing..."
    fi

    # echo "Initializing the configuration..."
    # ${bin} --home "${basedir}" parse genesis-file --genesis-file-path ./genesis.json

    print_info "Starting BDjuno..."
    if [ ! -f "${bin}" ]; then
        print_error "BDjuno executable not found: ${bin}"
        print_info "Please build the project first: make build"
        exit 1
    fi

    nohup "${bin}" start --home "${basedir}" >"${basedir}"/bdjuno.log 2>&1 &
    if [ $? -eq 0 ]; then
        print_info "BDjuno started (log: ${basedir}/bdjuno.log)"
    else
        print_error "Failed to start BDjuno"
        exit 1
    fi

    echo ""
    print_info "========================================"
    print_info "All services started"
    print_info "========================================"
    echo ""
    print_info "Access URLs:"
    echo "  - Hasura GraphQL: http://localhost:8080"
    echo "  - BDjuno log: ${basedir}/bdjuno.log"
    echo ""
}

function restart() {
    print_info "Restarting services (stop and restart, keep data)..."
    stop
    sleep 2
    print_info "Restarting services..."
    start
}

function reset() {
    print_warning "========================================"
    print_warning "WARNING: Reset will clear all data!"
    print_warning "========================================"
    echo ""
    print_warning "This will delete:"
    echo "  • All data in Docker volumes"
    echo "  • All data in data directory"
    echo "  • bdjuno.log file"
    echo "  • All indexed blocks and transactions"
    echo ""
    read -p "Confirm to continue? (type 'yes' to confirm): " confirm
    if [ "$confirm" != "yes" ]; then
        print_info "Reset operation cancelled"
        exit 0
    fi

    print_info "Stopping services and deleting all data..."
    docker compose -f "${basedir}"/docker-compose.yaml down -v
    if [ $? -ne 0 ]; then
        print_error "Failed to stop services"
        exit 1
    fi

    # Remove data directory using Docker to handle permissions
    if [ -d "${basedir}"/data ]; then
        print_info "Removing data directory..."
        # Use a temporary docker container to remove the data directory
        # This bypasses the permission issue since the container runs as root
        docker run --rm -v "${basedir}:/workspace" alpine rm -rf /workspace/data

        # Double check if it's gone, fallback to sudo if needed
        if [ -d "${basedir}"/data ]; then
            print_warning "Docker removal failed, trying sudo..."
            sudo rm -rf "${basedir}"/data
        fi
    fi

    # Remove log file
    if [ -f "${basedir}"/bdjuno.log ]; then
        rm -f "${basedir}"/bdjuno.log 2>/dev/null || {
            sudo rm -f "${basedir}"/bdjuno.log 2>/dev/null || true
        }
    fi

    # Kill bdjuno processes
    pids=$(ps -ef | grep "${bin_name}" | grep -v grep | awk '{print $2}')
    if [ -n "$pids" ]; then
        echo "$pids" | xargs kill -9 2>/dev/null || true
    fi

    print_info "Services stopped, all data cleared"
    sleep 2
    print_info "Restarting services..."
    start
}

function status() {
    print_info "Service status:"
    docker compose -f "${basedir}"/docker-compose.yaml ps
    echo ""
    print_info "BDjuno process:"
    pids=$(ps -ef | grep "${bin_name}" | grep -v grep | awk '{print $2}')
    if [ -n "$pids" ]; then
        echo "  BDjuno is running (PID: $pids)"
    else
        echo "  BDjuno is not running"
    fi
}

cmd=$1
case ${cmd} in
init)
    print_info "===== Initialize ====="
    if [ ! -f "${basedir}"/docker-compose.yaml ]; then
        print_error "docker-compose.yaml file not found"
        exit 1
    fi
    if [ ! -f "${bin}" ]; then
        print_error "BDjuno executable not found: ${bin}"
        print_info "Please build the project first: make build"
        exit 1
    fi
    print_info "Initialization completed"
    print_info "===== End ====="
    ;;
start)
    print_info "===== Start ====="
    start
    print_info "===== End ====="
    ;;
stop)
    print_info "===== Stop ====="
    stop
    print_info "===== End ====="
    ;;
restart)
    print_info "===== Restart ====="
    restart
    print_info "===== End ====="
    ;;
reset)
    print_info "===== Reset ====="
    reset
    print_info "===== End ====="
    ;;
status)
    status
    ;;
*)
    echo "Usage: localup.sh {init|start|stop|restart|reset|status}"
    echo ""
    echo "Commands:"
    echo "  init      - Check configuration and dependencies"
    echo "  start     - Start all services"
    echo "  stop      - Stop all services"
    echo "  restart   - Stop and restart all services (keep data)"
    echo "  reset     - Stop, clear all data and restart all services"
    echo "  status    - Show service status"
    echo ""
    echo "Examples:"
    echo "  ./localup.sh start"
    echo "  ./localup.sh restart"
    echo "  ./localup.sh reset"
    echo "  ./localup.sh status"
    ;;
esac
