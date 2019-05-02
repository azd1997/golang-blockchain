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


go mod replace
如果有对应的tag版本则直接在old(比如说golang.org/x/yyy)和new(比如说github/golang/yyy)后追加 v1.0.0形式的版本号
如果要导的那个包没有版本号，那么先去命令行go get new镜像中的这个包，就可以看到虚拟版本号，用以replace
本工程中golang/sys就是用的虚拟版本号
参考链接：在go modules中使用replace替换无法直接获取的package（golang.org/x/...） - apocelipes - 博客园
     https://www.cnblogs.com/apocelipes/archive/2018/09/08/9609895.html