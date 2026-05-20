#!/bin/bash
# ============================================================
# PostgreSQL 多数据库初始化脚本
# 在 PostgreSQL 容器首次启动时自动执行，创建多个数据库
# 使用方式：由 docker-entrypoint-initdb.d 自动调用
# 环境变量 POSTGRES_MULTIPLE_DATABASES 传入逗号分隔的数据库名列表
# ============================================================

set -e
set -u

# 创建多个数据库的函数
# 遍历逗号分隔的数据库名列表，为每个数据库执行创建操作
function create_multiple_databases() {
    local dblist="${POSTGRES_MULTIPLE_DATABASES:-}"
    echo "正在创建多个数据库: ${dblist}"

    local IFS=','
    for db in ${dblist}; do
        # 去除前后空格
        db=$(echo "${db}" | xargs)
        echo "检查数据库 ${db} 是否存在..."
        if psql -U "${POSTGRES_USER}" -lqt | cut -d \| -f 1 | grep -qw "${db}"; then
            echo "数据库 ${db} 已存在，跳过创建"
        else
            echo "创建数据库 ${db}..."
            createdb -U "${POSTGRES_USER}" "${db}"
            echo "数据库 ${db} 创建成功"
        fi
    done
}

# 执行数据库创建
create_multiple_databases
