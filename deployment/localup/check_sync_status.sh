#!/bin/bash

basedir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_section() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

check_bdjuno_process() {
    print_section "BDjuno 进程状态"
    if pgrep -f "bdjuno.*start" > /dev/null; then
        local pids=$(pgrep -f "bdjuno.*start")
        print_info "✓ BDjuno 正在运行 (PID: $pids)"
        return 0
    else
        print_error "✗ BDjuno 进程未运行"
        return 1
    fi
}

check_docker_services() {
    print_section "Docker 服务状态"
    cd "$basedir"
    if docker compose ps 2>/dev/null | grep -q "Up"; then
        docker compose ps
        return 0
    else
        print_error "Docker 服务未运行"
        return 1
    fi
}

check_database_sync() {
    print_section "数据库同步状态"
    
    cd "$basedir"
    
    local block_count=$(docker compose exec -T postgres psql -U postgres -d bdjuno -t -c "SELECT COUNT(*) FROM block;" 2>/dev/null | xargs)
    local latest_height=$(docker compose exec -T postgres psql -U postgres -d bdjuno -t -c "SELECT MAX(height) FROM block;" 2>/dev/null | xargs)
    local bucket_count=$(docker compose exec -T postgres psql -U postgres -d bdjuno -t -c "SELECT COUNT(*) FROM buckets;" 2>/dev/null | xargs)
    local bucket1_count=$(docker compose exec -T postgres psql -U postgres -d bdjuno -t -c "SELECT COUNT(*) FROM buckets WHERE bucket_name LIKE '%bucket1%' OR bucket_name = 'bucket1';" 2>/dev/null | xargs)
    
    if [ -z "$block_count" ] || [ "$block_count" = "0" ]; then
        print_warning "✗ 尚未同步任何区块 (block_count: ${block_count:-0})"
    else
        print_info "✓ 已同步区块数: $block_count"
        if [ -n "$latest_height" ] && [ "$latest_height" != "" ]; then
            print_info "  最新区块高度: $latest_height"
        fi
    fi
    
    echo ""
    if [ -z "$bucket_count" ] || [ "$bucket_count" = "0" ]; then
        print_warning "✗ 尚未同步任何 bucket (bucket_count: ${bucket_count:-0})"
    else
        print_info "✓ 已同步 bucket 数: $bucket_count"
    fi
    
    if [ -n "$bucket1_count" ] && [ "$bucket1_count" != "0" ]; then
        print_info "✓ 找到 bucket1 相关记录: $bucket1_count 条"
        echo ""
        print_info "Bucket1 详细信息:"
        docker compose exec -T postgres psql -U postgres -d bdjuno -c "SELECT bucket_name, owner_address, status, create_at FROM buckets WHERE bucket_name LIKE '%bucket1%' OR bucket_name = 'bucket1' ORDER BY create_at DESC LIMIT 5;" 2>/dev/null
    else
        print_warning "✗ 未找到 bucket1 相关记录"
    fi
}

check_bdjuno_logs() {
    print_section "BDjuno 日志检查"
    
    local log_file="${basedir}/bdjuno.log"
    
    if [ ! -f "$log_file" ]; then
        print_warning "日志文件不存在: $log_file"
        return 1
    fi
    
    local log_size=$(wc -l < "$log_file" 2>/dev/null | xargs)
    print_info "日志文件: $log_file (行数: ${log_size:-0})"
    
    echo ""
    print_info "最近的错误/警告:"
    tail -50 "$log_file" | grep -iE "error|warning|panic|fatal" | tail -10 || echo "  无错误信息"
    
    echo ""
    print_info "最近的同步信息:"
    tail -50 "$log_file" | grep -iE "block|height|sync|bucket" | tail -10 || echo "  无同步信息"
    
    echo ""
    print_info "最后 10 行日志:"
    tail -10 "$log_file"
}

check_port_3000() {
    print_section "端口 3000 检查"
    
    if command -v ss > /dev/null 2>&1; then
        local port_info=$(ss -tlnp | grep :3000 || echo "")
        if [ -n "$port_info" ]; then
            print_warning "端口 3000 已被占用:"
            echo "$port_info"
            return 1
        else
            print_info "✓ 端口 3000 未被占用"
            return 0
        fi
    else
        print_warning "无法检查端口状态 (ss 命令不可用)"
        return 0
    fi
}

main() {
    echo ""
    print_section "BDjuno 同步状态检查"
    echo ""
    
    check_docker_services
    echo ""
    
    check_bdjuno_process
    echo ""
    
    check_port_3000
    echo ""
    
    check_database_sync
    echo ""
    
    check_bdjuno_logs
    echo ""
    
    print_section "检查完成"
    echo ""
    print_info "如果 BDjuno 未运行，请执行:"
    echo "  cd $basedir && ./localup.sh restart"
    echo ""
    print_info "如果端口 3000 被占用，请先释放端口或修改配置"
    echo ""
}

main

