#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[deploy]${NC} $1"; }
warn() { echo -e "${YELLOW}[deploy]${NC} $1"; }
err() { echo -e "${RED}[deploy]${NC} $1"; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONTROLLER_COMPOSE="$SCRIPT_DIR/deployments/controller/docker-compose.yml"
WORKER_COMPOSE="$SCRIPT_DIR/deployments/agent-worker/docker-compose.yml"

check_redis() {
    if ! docker ps --format '{{.Names}}' | grep -q "redis"; then
        err "Redis container not found. Please start Redis first."
    fi

    REDIS_NETWORK=$(docker inspect redis --format '{{range $k, $v := .NetworkSettings.Networks}}{{$k}} {{end}}' 2>/dev/null | awk '{print $1}')
    log "Redis found on network: $REDIS_NETWORK"
}

check_network() {
    if ! docker network ls --format '{{.Name}}' | grep -q "^api$"; then
        log "Creating api network..."
        docker network create api
    fi
    log "Network api ready"

    if ! docker network inspect api --format '{{range .Containers}}{{.Name}} {{end}}' 2>/dev/null | grep -q "redis"; then
        log "Connecting Redis to api network..."
        docker network connect api redis 2>/dev/null || true
    fi
    log "Redis connected to api network"
}

build() {
    log "Building controller..."
    docker compose -f "$CONTROLLER_COMPOSE" build

    log "Building agent + worker..."
    docker compose -f "$WORKER_COMPOSE" build

    log "Build complete"
}

start() {
    check_redis
    check_network

    log "Starting controller..."
    docker compose -f "$CONTROLLER_COMPOSE" up -d

    log "Starting worker + agent..."
    docker compose -f "$WORKER_COMPOSE" up -d

    log "All services started"
    echo ""
    status
}

stop() {
    log "Stopping agent + worker..."
    docker compose -f "$WORKER_COMPOSE" down 2>/dev/null || true

    log "Stopping controller..."
    docker compose -f "$CONTROLLER_COMPOSE" down 2>/dev/null || true

    log "All services stopped"
}

restart() {
    stop
    start
}

status() {
    echo "=== Services ==="
    docker compose -f "$CONTROLLER_COMPOSE" ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null
    docker compose -f "$WORKER_COMPOSE" ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null
    echo ""
    echo "=== Endpoints ==="
    echo "  Controller: http://localhost:6001"
    echo "  Worker:     http://localhost:6002"
}

logs() {
    SERVICE=${1:-""}
    if [ -n "$SERVICE" ]; then
        case "$SERVICE" in
            controller) docker compose -f "$CONTROLLER_COMPOSE" logs -f ;;
            worker|agent) docker compose -f "$WORKER_COMPOSE" logs -f "$SERVICE" ;;
            *) err "Unknown service: $SERVICE (use: controller, worker, agent)" ;;
        esac
    else
        docker compose -f "$CONTROLLER_COMPOSE" logs -f &
        docker compose -f "$WORKER_COMPOSE" logs -f &
        wait
    fi
}

test_e2e() {
    log "Running E2E test..."
    echo ""

    log "1. Push config to controller..."
    curl -s -X POST http://localhost:6001/config \
        -H "X-API-Key: default-api-key" \
        -H "Content-Type: application/json" \
        -d '{"data":{"url":"https://httpbin.org/get"},"poll_interval_seconds":10}'
    echo ""

    log "2. Push config to worker..."
    curl -s -X POST http://localhost:6002/config \
        -H "Content-Type: application/json" \
        -d '{"version":1,"data":{"url":"https://httpbin.org/get"},"poll_interval_seconds":10}'
    echo ""

    log "3. Enqueue hit..."
    RESPONSE=$(curl -s http://localhost:6002/hit)
    echo "$RESPONSE"
    TASK_ID=$(echo "$RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['task_id'])" 2>/dev/null)
    echo ""

    if [ -z "$TASK_ID" ]; then
        err "Failed to get task_id"
    fi

    log "4. Waiting for task to process..."
    sleep 3

    log "5. Get result for task: $TASK_ID"
    curl -s "http://localhost:6002/hit/$TASK_ID"
    echo ""
    echo ""

    log "E2E test complete"
}

case "${1:-help}" in
    build)   build ;;
    start)   start ;;
    stop)    stop ;;
    restart) restart ;;
    status)  status ;;
    logs)    logs "$2" ;;
    test)    test_e2e ;;
    *)
        echo "Usage: $0 {build|start|stop|restart|status|logs|test}"
        echo ""
        echo "Commands:"
        echo "  build     Build all Docker images"
        echo "  start     Build + start all services"
        echo "  stop      Stop all services"
        echo "  restart   Stop + start all services"
        echo "  status    Show running services"
        echo "  logs      Tail logs (optional: controller, worker, agent)"
        echo "  test      Run E2E test"
        ;;
esac
