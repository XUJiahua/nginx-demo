# 测试nginx load balance效果

## 结论
确认会有load balance不成功的情况（502 bad gateway），但nginx足够智能，之后的请求会转发到没问题的进程／服务去的。

低峰时间部署，App接口没有耗时敏感操作，影响面小，用户重试一次就没问题了。所以APP后台没必要切换nginx。

## 测试用例

环境部署如下:

nginx:10080----app1:10000
          |
          |----app2:20000
          

客户端持续读这个连接,进程切换中会不会有问题:
http://192.168.1.196:10080

```
cd client
go run main.go http://192.168.1.196:10080

```

客户端逻辑：
模拟10个用户，发起GET操作，如果发现502就打印，发现10次502就退出。

## 测试进程部署
部署进程到Linux服务器:

```
GOOS=linux GOARCH=amd64 go build
rsync -rcv --progress . blackcat@192.168.1.196:~
```

Linux服务器上启动10000, 20000端口:

```
./nginx-demo -p &
./nginx-demo &
```

修改nginx配置:

sudo vim /etc/nginx/sites-enabled/default

```
upstream demo_nginx_server {
server 127.0.0.1:10000;
server 127.0.0.1:20000;
}

server {
        listen 10080;

        location / {
                proxy_pass http://demo_nginx_server;
        }
}

```

验证:

```
sudo nginx -t
sudo nginx -s reload
```