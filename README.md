# sponge-vue3-element-admin
sponge admin

# 如何使用

## 管理台

## 初始化 数据库
```mysql
    source server/db/admin.sql
```
## 前端 使用, 如有问题请移驾到 [vue3-element-admin](https://github.com/youlaitech/vue3-element-admin)
```npm
    cd web
    npm install pnpm -g
    pnpm install 
    pnpm dev 
```
## 服务端使用
```shell
    cd server
    go mod tidy
    make docs
    make run
```

## [代码生成器](https://github.com/ice-leng/sponge)
```git
    git clone https://github.com/ice-leng/sponge.git 
    git checkout vue3-element-admin
    cd sponge/cmd/sponge
    go run ./main.go upgrade
```

## 代码生成器使用
```shell
   git clone https://github.com/ice-leng/sponge-vue3-element-admin.git
   cd sponge-vue3-element-admin
   sponge run
```