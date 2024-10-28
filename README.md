# sponge-vue3-element-admin
sponge admin

# 如何使用

## 管理台

## 初始化 数据库
```mysql
    source server/deployments/sql/data.sql
```
## 前端 使用, 如有问题请移驾到 [vue3-element-admin](https://github.com/youlaitech/vue3-element-admin)
```npm
    cd web
    pnpm install 
    npm dev
```
## [sponge 安装](https://github.com/ice-leng/sponge)
目前只实现 通过web sql 生成代码。
```git
    git clone https://github.com/ice-leng/sponge.git 
    git checkout vue3-element-admin
    cd sponge/cmd/sponge
    go run ./main.go init
```