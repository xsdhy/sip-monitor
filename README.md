# SIP 信令收集和监控系统

支持Mysql、MongoDB、sqlite
启动一个Hep Server监听9060端口，用于收集SIP信令
启动一个Http Server监听9059端口，用于页面调用API

## 数据库抽象层

该系统实现了一个灵活的数据库抽象层，使用了Repository设计模式和工厂模式，支持以下数据库：

1. **MySQL** - 适用于大规模生产环境
2. **SQLite** - 适用于轻量级部署或开发环境



## 目录结构
后端项目在src目录，
前端项目在web目录，前端使用react+ant design