# 考核（简略版）

## 实现的功能
* 上传、下载
* 登录注册
* 分包上传（初始化分块上传信息，上传文件分块，合并分块）(很多问题)
* 秒传

## 接口
* POST：/user/signup
> 参数：   
用户名：username   
密码：password  
* POST：/user/signin
> 参数：             
用户名：username       
密码：password                  

#### 登录成功之后会返回token，以下所有接口都要把username，token作为表单参数
#### 文件哈希生成
```go
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}
```

* POST：/file/upload
> 参数： file  
用户名：username     
token：  token    

* GET: /file/downloadPOST  
参数: 
用户名： username    
token: token  
文件哈希： filehash

* POST: /file/fast（尝试秒传）
参数：   
用户名： username  
token: token  
文件哈希：filehash 
文件名：filename                    
文件大小： filesize                

* GET: /file/qury(查询）                             
参数：                        
用户名： username                           
token: token                        
行数： limit                        

* POST: /multiPart/init (初始化分块信息）                       
参数：                    
用户名： username             
文件哈希： filehash                 
文件大小: filesize                      
                         
             
* POST: /multiPart/upload (上传文件分块）               
参数:      
分块信息id： uploadid   
序号:  index  

                                 
* POST: /multiPart/complete (合并分块)                            
参数：
用户名： username     
文件哈希： filehash     
文件大小: filesize         
文件名：filename   
分块信息id： uploadid               
