> 接入层接口
```
//查询秒杀接口信息
http://127.0.0.1:8082/sec/info?product_id=3

//查询秒杀商品列表
http://127.0.0.1:8082/sec/list

//秒杀商品接口
http://127.0.0.1:8082/sec/kill
//参数
product_id: 1
user_id: 1
src: 192.168.199.1
auth_code: userauthcode
time: 1530928164
nance: dsdsdjkdjskdjksdjhuieurierei
```

> 管理层接口
```
//添加商品接口
http://127.0.0.1:8081/product/create
//参数
product_name:梨子
product_total:100
status:1

//商品列表接口接口
http://127.0.0.1:8081/product/list

//秒杀活动列表接口
http://127.0.0.1:8081/activity/list

//添加秒杀活动接口
http://127.0.0.1:8081/activity/create
//参数
activity_name: 梨子大甩卖
product_id:4
start_time:1530928052
end_time:1530989052
total:20
status:1
speed:1
buy_limit:1
buy_rate:0.2
```
