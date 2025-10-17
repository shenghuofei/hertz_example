FROM rocky9.2-go1.24:v2 AS builder

# 工作目录（放源码）
WORKDIR /build

# 复制项目所有源代码
COPY . .

# 缓存依赖
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go mod download

RUN /usr/bin/bash -c "./build.sh"

# ---------------------------
# 第2阶段：生成最终运行镜像
# ---------------------------
FROM rocky9.2-go1.24:v2

# 创建运行目录
WORKDIR /data0/app/avatar

# 从builder阶段复制二进制
COPY --from=builder /build/output .

# 启动命令
CMD ["/usr/bin/bash","./start.sh"]
