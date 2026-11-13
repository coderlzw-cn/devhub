# DevHub

Vue 3、TypeScript、Vite、Tailwind CSS 4 与 TanStack Vue Query 构建的站点与分类管理应用。

## 启动与构建

```bash
pnpm install
pnpm dev
pnpm build
```

后端默认运行于 `http://127.0.0.1:8080`，开发环境通过 Vite 代理 `/api`。生产环境需配置同源反向代理或构建时设置 `VITE_API_BASE_URL`。

## 功能与数据

- 分类管理：真实分类查询、新增、编辑、删除、名称搜索与两级分类。新增需填写所属用户 UUID 和 slug，新增子分类默认沿用父分类用户；有子分类时需先删除子分类。
- 站点管理：真实站点查询、新增、编辑、删除，分类筛选、标题/网址/描述搜索、名称排序、分页、图标地址、排序值和置顶。查询读取全部后端分页后进行界面筛选和分页。新增需填写真实所属用户 UUID；所属分类使用后端 ID。
- 编辑站点时可切换分类；当前 OpenAPI 的 PATCH 不支持 `null` 分类 ID，因此已有分类的站点无法在编辑时清空分类。
- 首页总数、分类数量和置顶站点均来自真实接口。
- 当前 OpenAPI 没有当前用户资料接口。个人中心为用户自行填写的本地资料，首次为空；偏好与本地资料仍保存在浏览器，JSON 导出包含本地资料、真实站点和偏好。
- 登录页尚未接入认证服务。

已移除应用 mock 数据、MSW handlers、worker 文件及启动逻辑，保留 `package.json` 中的 `msw` 依赖。启动时注销旧版本遗留的 MSW worker；旧的 `devhub-resources` 和 `devhub-profile` 数据不会再被读取。

## API 生成与验证

```bash
pnpm api:generate  # 从 OpenAPI 快照生成分类与站点钩子和类型
pnpm api:refresh   # 从本机 8080 更新 OpenAPI 快照并生成
pnpm exec vitest run tests
pnpm lint
```

生成配置位于 `orval.config.ts`，快照位于 `openapi/backend.json`。`src/api/generated` 由 Orval 自动生成，不应手工修改。测试中的 HTTP 替身仅用于隔离验证，不参与应用运行。
