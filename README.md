## 统一定时任务中心
灵感来自`beego`的task

##### 功能：
+ 支持定时任务注册
+ 支持定时任务修改
+ 支持定时任务删除
+ 支持定时任务暂停运行
+ 支持定时任务恢复运行
+ 支持定时任务运行(跨节点远程rpc调用)

##### 逻辑：
+ 多个节点运行时，从数据库加载同一份定时任务列表，每个任务均会加redis分布式锁
+ 后续在节点：注册新的定时任务 / 原定时任务修改、删除、暂停、恢复时，会通过redis发布/订阅同步到所有节点，并且任务也会加上redis分布式锁

##### 特点：
+ 多节点部署
+ redis分布式锁防止：同一定时任务在多个节点同时运行
+ redis publish && subscribe 机制，任务增加、修改、删除、暂停、恢复等动作，同步通知到所有节点

##### 不足：
`handle/conversion.go`实现了**跨节点远程rpc调用**关键功能：**方法寻址**。
抱歉，此功能是由内部框架实现，内部框架暂时不能开源。请自行实现寻址，后期将从框架切分，加入不依赖于框架的寻址

`log`请替换自己的log包

##### 截图：
###### 节点首次运行，从数据库加载所有任务
![Alt text](https://github.com/shenping1916/cronservice/blob/master/images/1537935279389.jpg)

![Alt text](https://github.com/shenping1916/cronservice/blob/master/images/1537934583881.jpg)

###### 节点收到redis订阅消息
![Alt text](https://github.com/shenping1916/cronservice/blob/master/images/1537934876826.jpg)

###### 任务争抢`redis`锁，拿到锁的节点运行定时任务
![Alt text](https://github.com/shenping1916/cronservice/blob/master/images/1537933646132.jpg)
  