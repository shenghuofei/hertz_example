### 这个是使用hertz框架的实例
主要包含
1. 使用Viper 根据环境信息读取对应的配置文件
2. 日志，个性化日志配置及rotate
3. 中间件，使用中间件打印access log
4. GORM集成，支持多db及主从分离，及GORM使用
5. 路由分组
6. 按照controller(handler)，service，model分层组织项目结构
7. 统一响应结构体，并使用中间件捕获异常，某个方法panic不会导致进程退出
8. 支持cronjob，自动捕获cronjob中的panic，单个cronjob panic不会导致进程退出


### 相关链接
[hertz官方文档](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/middleware/basic-auth/)  
[hertz github](https://github.com/cloudwego/hertz)  
[hertz官方提供的一些示例github](https://github.com/cloudwego/hertz-examples/tree/main)
