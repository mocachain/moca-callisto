#!/bin/bash

# Bucket 创建同步测试脚本
# 用于验证 bdjuno 能否正确索引 bucket 创建事件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
MOCA_CMD_PATH="${MOCA_CMD_PATH:-./build/moca-cmd}"
MOCA_CMD_HOME="${MOCA_CMD_HOME:-./deployment/localup}"
MOCA_CMD_PASSWORD_FILE="${MOCA_CMD_PASSWORD_FILE:-./deployment/localup/testkey/password.txt}"
POSTGRES_CONTAINER="${POSTGRES_CONTAINER:-localup-postgres-1}"
DB_NAME="${DB_NAME:-bdjuno}"
WAIT_TIME="${WAIT_TIME:-10}"

# 函数：打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 函数：检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        print_error "命令 '$1' 未找到，请先安装"
        exit 1
    fi
}

# 函数：检查 moca-cmd 是否存在
check_moca_cmd() {
    if [ ! -f "$MOCA_CMD_PATH" ]; then
        print_error "moca-cmd 未找到: $MOCA_CMD_PATH"
        print_info "请先编译 moca-cmd: cd <moca-cmd-root> && make build"
        exit 1
    fi
}

# 函数：检查 postgres 容器是否运行
check_postgres() {
    if ! docker ps | grep -q "$POSTGRES_CONTAINER"; then
        print_error "Postgres 容器未运行: $POSTGRES_CONTAINER"
        print_info "请先启动服务: cd <moca-callisto-root> && bash deployment/localup/localup.sh start"
        exit 1
    fi
}

# 函数：检查 bdjuno 是否运行
check_bdjuno() {
    if ! pgrep -f "bdjuno start" > /dev/null; then
        print_error "bdjuno 未运行"
        print_info "请先启动 bdjuno: cd <moca-callisto-root> && ./build/bdjuno start --home deployment/localup"
        exit 1
    fi
}

# 函数：查询数据库
query_db() {
    docker exec -i "$POSTGRES_CONTAINER" psql -U postgres -d "$DB_NAME" -t -c "$1" 2>/dev/null | xargs
}

# 函数：创建 bucket
create_bucket() {
    local bucket_name="$1"
    print_info "创建 bucket: $bucket_name"

    local output
    output=$("$MOCA_CMD_PATH" --home "$MOCA_CMD_HOME" \
        --passwordfile "$MOCA_CMD_PASSWORD_FILE" \
        bucket create "$bucket_name" 2>&1)

    if echo "$output" | grep -q "transaction hash"; then
        local tx_hash=$(echo "$output" | grep "transaction hash" | awk '{print $NF}')
        print_info "Bucket 创建交易已提交: $tx_hash"
        echo "$tx_hash"
        return 0
    else
        print_error "Bucket 创建失败"
        echo "$output"
        return 1
    fi
}

# 函数：等待 bucket 同步
wait_for_bucket_sync() {
    local bucket_name="$1"
    local max_wait="$2"
    local elapsed=0

    print_info "等待 bucket 同步到数据库 (最多等待 ${max_wait}s)..."

    while [ $elapsed -lt $max_wait ]; do
        local count=$(query_db "SELECT COUNT(*) FROM buckets WHERE bucket_name = '$bucket_name';")

        if [ "$count" = "1" ]; then
            print_info "✓ Bucket 已同步到数据库 (耗时: ${elapsed}s)"
            return 0
        fi

        sleep 1
        elapsed=$((elapsed + 1))

        # 每 2 秒显示一次进度
        if [ $((elapsed % 2)) -eq 0 ]; then
            echo -n "."
        fi
    done

    echo ""
    print_error "✗ Bucket 同步超时 (${max_wait}s)"
    return 1
}

# 函数：验证 bucket 数据
verify_bucket_data() {
    local bucket_name="$1"

    print_info "验证 bucket 数据..."

    local query="SELECT bucket_name, owner_address, status, global_virtual_group_family_id 
                 FROM buckets 
                 WHERE bucket_name = '$bucket_name';"

    local result=$(docker exec -i "$POSTGRES_CONTAINER" psql -U postgres -d "$DB_NAME" -c "$query" 2>/dev/null)

    if echo "$result" | grep -q "$bucket_name"; then
        print_info "✓ Bucket 数据验证成功"
        echo "$result"
        return 0
    else
        print_error "✗ Bucket 数据验证失败"
        return 1
    fi
}

# 函数：清理测试数据
cleanup_test_bucket() {
    local bucket_name="$1"

    print_info "清理测试 bucket: $bucket_name"

    # 尝试删除 bucket（可能失败，因为需要先删除对象等）
    "$MOCA_CMD_PATH" --home "$MOCA_CMD_HOME" \
        --passwordfile "$MOCA_CMD_PASSWORD_FILE" \
        bucket rm "$bucket_name" 2>/dev/null || true
}

# 函数：显示统计信息
show_statistics() {
    print_info "数据库统计信息:"

    local sp_count=$(query_db "SELECT COUNT(*) FROM storage_providers;")
    local gvgf_count=$(query_db "SELECT COUNT(*) FROM global_virtual_group_families;")
    local bucket_count=$(query_db "SELECT COUNT(*) FROM buckets;")
    local object_count=$(query_db "SELECT COUNT(*) FROM objects;")

    echo "  - Storage Providers: $sp_count"
    echo "  - Virtual Group Families: $gvgf_count"
    echo "  - Buckets: $bucket_count"
    echo "  - Objects: $object_count"
}

# 主测试流程
main() {
    print_info "========================================"
    print_info "Bucket 创建同步测试"
    print_info "========================================"
    echo ""

    # 检查依赖
    print_info "检查依赖..."
    check_command docker
    check_moca_cmd
    check_postgres
    check_bdjuno
    print_info "✓ 所有依赖检查通过"
    echo ""

    # 显示初始统计
    show_statistics
    echo ""

    # 生成唯一的 bucket 名称
    local bucket_name="test-sync-$(date +%s)"
    local tx_hash=""

    # 创建 bucket
    if ! tx_hash=$(create_bucket "$bucket_name"); then
        print_error "测试失败: Bucket 创建失败"
        exit 1
    fi
    echo ""

    # 等待同步
    if ! wait_for_bucket_sync "$bucket_name" "$WAIT_TIME"; then
        print_error "测试失败: Bucket 同步超时"

        # 显示调试信息
        print_warning "检查 bdjuno 日志以获取更多信息:"
        print_warning "  tail -50 deployment/localup/bdjuno.log"

        cleanup_test_bucket "$bucket_name"
        exit 1
    fi
    echo ""

    # 验证数据
    if ! verify_bucket_data "$bucket_name"; then
        print_error "测试失败: Bucket 数据验证失败"
        cleanup_test_bucket "$bucket_name"
        exit 1
    fi
    echo ""

    # 显示最终统计
    show_statistics
    echo ""

    # 清理（可选）
    if [ "${CLEANUP:-true}" = "true" ]; then
        cleanup_test_bucket "$bucket_name"
    else
        print_warning "跳过清理，测试 bucket 保留: $bucket_name"
    fi

    echo ""
    print_info "========================================"
    print_info "✓ 测试通过！"
    print_info "========================================"
}

# 显示使用说明
show_usage() {
    cat << EOF
用法: $0 [选项]

选项:
  -h, --help              显示此帮助信息
  -w, --wait-time SECONDS 设置等待同步的最大时间（默认: 10 秒）
  --no-cleanup            测试后不清理 bucket
  --moca-cmd PATH         指定 moca-cmd 路径
  --postgres CONTAINER    指定 postgres 容器名称

环境变量:
  MOCA_CMD_PATH           moca-cmd 可执行文件路径
  MOCA_CMD_HOME           moca-cmd home 目录
  MOCA_CMD_PASSWORD_FILE  密码文件路径
  POSTGRES_CONTAINER      postgres 容器名称
  WAIT_TIME               等待同步的最大时间（秒）
  CLEANUP                 是否清理测试数据（true/false）

示例:
  # 基本用法
  $0

  # 自定义等待时间
  $0 --wait-time 20

  # 保留测试数据
  $0 --no-cleanup

  # 使用环境变量
  WAIT_TIME=15 CLEANUP=false $0

EOF
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -w|--wait-time)
            WAIT_TIME="$2"
            shift 2
            ;;
        --no-cleanup)
            CLEANUP=false
            shift
            ;;
        --moca-cmd)
            MOCA_CMD_PATH="$2"
            shift 2
            ;;
        --postgres)
            POSTGRES_CONTAINER="$2"
            shift 2
            ;;
        *)
            print_error "未知选项: $1"
            show_usage
            exit 1
            ;;
    esac
done

# 运行主测试
main

