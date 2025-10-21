FROM rocky9.2-go1.24:v3 AS builder

# 工作目录（放源码）
WORKDIR /build

# 设置国内代理
ENV GOPROXY=https://goproxy.cn,direct

# 先复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 清缓存 + 下载依赖
RUN go mod download

# 复制项目所有源代码
COPY . .

RUN ./build.sh

# ---------------------------
# 第2阶段：生成最终运行镜像
# ---------------------------
FROM rocky9.2-go1.24:v3

# 创建运行目录
WORKDIR /data0/app/avatar

# 从builder阶段复制二进制
COPY --from=builder /build/output .

# 启动命令
CMD ["/usr/bin/bash","./start.sh"]
