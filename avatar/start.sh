#!/bin/sh
set -e

# 打印启动信息（可选）
echo "Starting Go service..."

appName=avatar
workDir=$(cd $(dirname $0) && pwd)
cd $workDir
# 用 exec 启动 Go 进程，替换当前 shell
exec ./$appName "$@"
