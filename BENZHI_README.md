# bateman-chain

bateman-chain 是放射性线性衰变链求解器：给定各核素衰变常数 λ（1/s）与初始核数，命令行沿时间输出每一核的核数 N 与活度 A（A = λ·N）。

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
