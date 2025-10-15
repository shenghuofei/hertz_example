# hertz-mysql 示例

## 依赖
- Go 1.20+
- MySQL

## 运行
1. 修改 config.yaml/prod.yaml/test.yaml（或通过环境变量覆盖）
2. 初始化 go module: go mod tidy
3. 运行
   go run ./cmd/main.go
4. goland run 配置  
   run kind: package  
   package path: avatar/cmd  
   working directory: /xxx.../avatar  
   env: DATABASE_DEFAULT_PASSWORD=xxxxx;DATABASE_DEFAULT_USER=xxx;DATABASE_K8S-IAC_PASSWORD=xxx;DATABASE_K8S-IAC_USER=xxx;ENV=test

服务：
- HTTP: http://localhost:8080
  - / -> simple landing
  - /healthz -> health check pings write DB
- Prometheus metrics: http://localhost:9090/metrics

