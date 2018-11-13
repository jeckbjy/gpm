# gpm

简单的go项目管理,不依赖全局的gopath,大部分代码源自glide,但去除了复杂的依赖查找, 仅有添加,删除依赖,下载依赖到vendor目录

- 实现目标
  - 1:不依赖GOPATH,仅依赖gpm.yaml配置所在目录,与vendor在同一级目录
  - 2:支持创建，安装，删除，更新操作
  - 3:支持编译,自动检测GOPATH
- TODO:
  - version管理
  - lock文件
  - 依赖树管理