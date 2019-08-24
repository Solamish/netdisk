# 考核（简略版，晚点再完善）

## 实现的功能
* 上传、下载
* 登录注册
* 分包上传（初始化分块上传信息，上传文件分块，合并分块）(很多问题)
* 秒传

## 接口
* /user/signup 
> 参数：   
用户名：username   
密码：password  
* /user/signin      
> 参数：             
用户名：username       
密码：password                  

## 登录成功之后会返回token，以后所有接口都要把username，token作为表单参数


* /file/upload
参数： file
用户名：username   
token：  token  

* /file/download
 参数: 
filehash

* /file/fast
参数： filehash, filename, filesize


