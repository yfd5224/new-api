# New API Fork

🍥 **新一代大模型网关与AI资产管理系统 - Fork 版本**

> ⚠️ 此版本为 Fork 版本，基于 [QuantumNous/new-api](https://github.com/QuantumNous/new-api)

---

## 总结：你已修改了什么，还有什么未改

### 已完成的修改

| 模块 | 文件 | 做了什么 |
|------|------|---------|
| 数据模型 | `model/token_channel_override.go` | 新增 `TokenChannelOverride` 表和全套 CRUD，按 `token_id` + `channel_id` 联合唯一 |
| 数据库迁移 | `model/main.go` | 在 `AutoMigrate` 和 `migrateDBFast` 中注册新表 |
| 转发核心 | `middleware/distributor.go` | 在 `SetupContextForSelectedChannel` 里，获取渠道默认 key 后，查覆盖表；如有则替换，并往 context 打 `key_override` 标记 |
| 管理接口 | `controller/token_channel_override.go` | 5 个 API：查单 Token 覆盖、新增、修改、删除单条、查当前用户全部覆盖 |
| 路由 | `router/api-router.go` | 注册上述接口到 `/api/token/:id/overrides` 等 |
| 级联清理 | `model/token.go` | 删 Token 时自动清掉它的覆盖记录 |
| 级联清理 | `model/channel.go` | 删 Channel 时自动清掉关联的覆盖记录 |
| 前端编辑页 | `EditTokenModal.jsx` | 新增"渠道密钥覆盖"卡片：展示已配置列表、下拉选渠道、输入覆盖 key、逐条删除 |
| 日志标记 | `service/log_info_generate.go` | 在日志 `admin_info` 里写入 `key_override: true`，方便后台区分 |
| Bug 修复 | `EditTokenModal.jsx` | 修复 `channels.map is not a function`（分页对象未取 `.items`） |
| UI 修复 | `EditTokenModal.jsx` | 修复 Select 下拉截断和 Popconfirm 溢出 |

### 还未修改的遗留点

| 位置 | 问题 | 影响范围 |
|------|------|---------|
| `relay/relay_task.go:440` | 直接读 `channelModel.Key`，未走 context 覆盖 | Gemini / VertexAI 异步任务轮询时，覆盖 key 不生效 |
| `relay/mjproxy_handler.go:300,475` | 直接读 `channel.Key`，未走 context 覆盖 | Midjourney 任务（放大/变换/重绘）时，覆盖 key 不生效 |

如果不用 Gemini 视频/图片轮询和 Midjourney，这两处可以暂时不管；如果要用，需要把它们也改成从 context 取 key（像主链路那样）。

### 计费 & 用量统计是否需要改？

**不需要**。目前改动只替换了向上游发请求的 `Authorization` key，渠道ID、模型名、渠道类型、BaseURL 全都没变。系统计费只认"哪个渠道、哪个模型、多少 token"，不认 key 内容。

所以：
- Token 额度：照旧，各自独立
- 用量日志：`Logs` 表仍记录同一个 `channel_id`
- 价格/表达式：计算依据不变
- 后台统计：不受影响

---

## � 快速启动

### 方式一：分别启动前后端（推荐开发）

#### 1. 启动后端服务

```bash
# 在项目根目录执行
go run main.go
```

后端将运行在：`http://localhost:3000`

#### 2. 启动前端开发服务器

```bash
# 进入 web 目录
cd web

# 启动开发服务器
bun run dev
```

前端将运行在：`http://localhost:5173`

#### 3. 访问应用

打开浏览器访问：**http://localhost:5173**

---

### 方式二：Docker 启动

```bash
# 拉取最新镜像
docker pull calciumion/new-api:latest

# 使用 SQLite（默认）
docker run --name new-api -d --restart always \
  -p 3000:3000 \
  -e TZ=Asia/Shanghai \
  -v ./data:/data \
  calciumion/new-api:latest

# 使用 MySQL
docker run --name new-api -d --restart always \
  -p 3000:3000 \
  -e SQL_DSN="root:123456@tcp(localhost:3306)/oneapi" \
  -e TZ=Asia/Shanghai \
  -v ./data:/data \
  calciumion/new-api:latest
```

> **💡 提示：** `-v ./data:/data` 会将数据保存在当前目录的 `data` 文件夹中

---

##  许可证

本项目采用 [GNU Affero 通用公共许可证 v3.0 (AGPLv3)](./LICENSE) 授权。
