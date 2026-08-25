# bateman-chain — Go 放射性衰变链求解 Web 服务与 CLI 工具（Bateman 方程闭式解 + RK4 数值积分）
bateman-chain 是一个 Go 实现的放射性线性衰变链求解器，通过 HTTP JSON API 或命令行输出各核素核数 N 与活度 A 随时间的演化，支持 Bateman 闭式求和与自适应 RK4 数值积分两种求解路径。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . eval example/u238-short.json --t 3600
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
