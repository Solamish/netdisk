# 考核（简略版，晚点再完善）

## 实现的功能
* 上传、下载
* 登录注册
* 分包上传（初始化分块上传信息，上传文件分块，合并分块）(bug)

## 接口
* /user/signup 
> 参数： username 
         password
* /user/signin
> 参数： username  password  

## 登录成功之后会返回token，以后所有接口都要传username，token作为参数

* /file/upload
参数： 选择文件


* /file/download
参数: filehash

* /file/fast
参数： filehash, filename, filesize


