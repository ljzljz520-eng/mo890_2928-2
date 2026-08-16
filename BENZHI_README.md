# 社区长者记忆墙

这是一个单一 Go module 的社区记忆墙示例，使用固定内存数据运行，不需要数据库、网络服务、时钟或随机数。工作人员可以维护长者故事、生活照片、重要年份和家属寄语；家属提交的内容先进入待审列表，工作人员审核后才会出现在长者资料中。每位长者可设置公开、家属可见或私密范围，并可下载确定性 JSON 资料包。

## 运行

需要 Go 1.25.13。模块根目录执行：

```bash
go run ./cmd/memorywall -addr :8080
```

浏览器打开 `http://localhost:8080/` 可查看前端验收界面。服务 API 包括：

- `GET /api/elders`：公开长者列表
- `GET /api/elders/{id}`：长者详情
- `POST /api/import-batches`：提交一批待审导入项目
- `GET /api/submissions/pending`：待审内容
- `POST /api/submissions/{id}/review`：审核内容
- `PUT /api/elders/{id}/visibility`：设置公开范围
- `GET /api/export`：导出整站资料

## 前端

```bash
cd web
npm install
npm run build
```

构建脚本使用 Node.js 20，把静态验收页面写入 `web/dist/`。

## 业务链路验证

```bash
go test -count=1 ./...
```

测试数据固定且不依赖环境状态。导入批次验收覆盖包含格式错误项目的场景。

## Docker 打包

```bash
./build_benzhi_docker.sh community-memory-wall linux/arm64
./build_benzhi_docker.sh community-memory-wall linux/amd64
```
