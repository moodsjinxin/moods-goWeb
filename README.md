# entrytask


- 简介

  ​		本项目是实习入职第一周的entrytask，主要目标是实现go web开发中的用户登录、注册、鉴权、个人信息修改（头像和昵称）等基本功能。主要目的是为了熟悉go web开发的方法、tcp socket通信、http请求的处理、go-redis的使用、go-mysql的使用以及熟悉go的语法和代码规范。

  ​		本项目共分为http server和tcp server。httpserver主要是接受前端页面的请求，将请求通过tcp socket发送给tcp服务器，由tcp服务器实现功能逻辑。

  ​	![image](https://github.com/moodsjinxin/moods-goWeb/blob/master/images/%E5%9F%BA%E6%9C%AC%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

- git项目结构

  .

  ├── README.md.             

  ├── entrytasktest                     //测试程序项目

  │  ├── benchmark.go            //压测代码

  │  └── sqlinsert.go			   //千万数据插入sql代码

  ├── images							 

  ├── myhttpsever				  // http server 项目

  │  ├── WebPages				//web html静态资源

  │  ├── controllers				//请求处理功能

  │  ├── main.go			

  │  └── userpictures			//头像静态资源

  ├── mytcpsever					// tcp server 项目

  │  ├── controllers				//功能实现目录

  │  └── main.go					

  ├── 	测试报告.md

  └──   设计文档.md





- 项目架构

  ![架构图](https://github.com/moodsjinxin/moods-goWeb/blob/master/images/架构图.png)

- 相关文档

  - [设计文档](./设计文档.md "设计文档")
  - [测试报告](./测试报告.md "测试报告")
  - [总结文档](./总结文档.md "总结文档")
  - 部署见下

  

- 部署

  - git clone 当前项目

  - 配置import

    - 项目中引入了2个三方库：redigo（redis）和 go-sql-driver（mysql）根据import内容可以自动进行下载
    
    ![image](https://github.com/moodsjinxin/moods-goWeb/blob/master/images/import配置.png)
    
  - 修改mytcpsever中redis和mysql的配置信息（按照自己的redis 和 mysql信息进行配置）

    ![image](https://github.com/moodsjinxin/moods-goWeb/blob/master/images/%E4%BF%AE%E6%94%B9mysql%E9%85%8D%E7%BD%AE.png)

  - 修改文件上传存储位置
    
    -  在tcpserver中的profilescontroller中用户头像的存储位置应该放置在自己git下来的httpserver的userpictures文件夹中，本项目用的是绝对路径存储。
