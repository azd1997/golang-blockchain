1.进入工程目录下（注意该目录不能在$GOPATH/src下级）
go mod init nameOfMod
2.手动编辑go.mod添加require，具体要求自行查阅。
通用为：
require （
        github.com/xxx/yyyy v1.1.1
）
3.执行
go mod tidy [-v]
下载相关require以及移除没用的require，生成go.sum（记录各require的哈希）
-v表示打印细节到控制台
4.验证
go mod verify
验证所有需要的模块是否已下载以及是否被篡改
5. go mod vendor [-v]
在工程根目录生成vendor文件夹，放置go.mod描述的依赖包
6.在工程中的代码中import对应的依赖包
