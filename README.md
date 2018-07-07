### SecKill
> 这是基于Go语言的一个秒杀系统，这个系统分三层，接入层、逻辑层、管理层。

##### 系统架构图
![image](https://github.com/BlueSimle/SecKill/blob/master/framework.png)

##### 秒杀接入层
1. 从Etcd中加载秒杀活动数据到内存当中。
2. 监听Etcd中的数据变化，实时加载数据到内存中。
3. 从Redis中加载黑名单数据到内存当中。
4. 设置白名单。
5. 对用户请求进行黑名单限制。
6. 对用户请求进行流量限制、秒级限制、分级限制。
7. 将用户数据进行签名校验、检验参数的合法性。
8. 接收逻辑层的结果实时返回给用户。


##### 秒杀逻辑层
1. 从Etcd中加载秒杀活动数据到内存当中。
2. 监听Etcd中的数据变化，实时加载数据到内存中。
3. 处理Redis队列中的请求。
4. 限制用户对商品的购买次数。
5. 对商品的抢购频次进行限制。
5. 对商品的抢购概率进行限制。
6. 对合法的请求给予生成抢购资格Token令牌。

##### 秒杀管理层
1. 添加商品数据。
2. 添加抢购活动数据。
3. 将数据同步到Etcd。
4. 将数据同步到数据库。


##### 目录结构
```
├─sk_admin
│  ├─config
│  ├─controller
│  │  ├─activity
│  │  └─product
│  ├─model
│  ├─service
│  └─setup
├─sk_layer
│  ├─config
│  ├─logic
│  ├─service
│  │  ├─srv_err
│  │  ├─srv_limit
│  │  ├─srv_product
│  │  ├─srv_redis
│  │  └─srv_user
│  └─setup
├─sk_proxy
│  ├─config
│  ├─controller
│  ├─service
│  │  ├─srv_err
│  │  ├─srv_limit
│  │  ├─srv_redis
│  │  └─srv_sec
│  └─setup
└─vendor
    └─github.com
        ├─coreos
        │  └─etcd
        │      └─clientv3
        ├─gin-gonic
        │  └─gin
        ├─go-sql-driver
        │  └─mysql
        ├─gohouse
        │  └─gorose
        ├─spf13
        │  ├─cobra
        │  └─viper
        └─Unknwon
            └─com
```