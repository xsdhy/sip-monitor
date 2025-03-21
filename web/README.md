# SIP监控平台前端

## 技术栈

- Vite - 构建工具
- React - UI框架
- TypeScript - 类型安全
- React Router - 路由管理
- Axios - HTTP请求
- Ant Design - UI组件库
- Day.js - 日期处理

## 开发

```bash
# 安装依赖
npm install
# 或
yarn

# 启动开发服务器
npm run dev
# 或
yarn dev
```

## 构建

```bash
# 构建生产版本
npm run build
# 或
yarn build

# 预览生产构建
npm run preview
# 或
yarn preview
```

## 项目结构

```
web/
├── public/          # 静态资源
├── src/
│   ├── @types/      # TypeScript类型定义
│   ├── components/  # 公共组件
│   ├── utils/       # 工具函数
│   ├── views/       # 页面组件
│   ├── main.tsx     # 入口文件
│   └── index.css    # 全局样式
├── .env             # 环境变量
├── index.html       # HTML模板
├── package.json     # 项目依赖和脚本
├── tsconfig.json    # TypeScript配置
├── vite.config.ts   # Vite配置
└── README.md        # 项目说明
```

## 最佳实践

- 使用类型安全的 TypeScript
- 模块化组件设计
- 使用环境变量管理配置
- 统一的错误处理
- 路由集中管理 