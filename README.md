# SIP 信令收集和监控系统

支持Mysql、MongoDB、sqlite
启动一个Hep Server监听9060端口，用于收集SIP信令
启动一个Http Server监听9059端口，用于页面调用API

## 数据库抽象层

该系统实现了一个灵活的数据库抽象层，使用了Repository设计模式和工厂模式，支持以下数据库：

1. **MySQL** - 适用于大规模生产环境
2. **MongoDB** - 适用于需要灵活数据结构的环境
3. **SQLite** - 适用于轻量级部署或开发环境

### 配置数据库

在环境配置中指定数据库类型：

```
# MySQL配置示例
DB_TYPE=mysql
DB_USER=username
DB_PASSWORD=password
DB_ADDR=localhost:3306
DB_NAME=sip_monitor

# MongoDB配置示例
DB_TYPE=mongodb
DSN_URL=mongodb://username:password@localhost:27017
DB_NAME=sip_monitor

# SQLite配置示例
DB_TYPE=sqlite
DB_PATH=./sip_monitor.db
```

### 数据库操作接口

所有数据库操作都通过Repository接口进行抽象



## 目录结构
-resources
-src 后端项目
--entity 实体定义
--model 数据库操作
--pkg 第三方包
--services 服务层
-web 前端项目