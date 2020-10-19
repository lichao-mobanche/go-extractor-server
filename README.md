# go-extractor-server
该项目提供提链服务，使用http+json与外界交互，灵活配置工作线程数量，使用简单方便。
## 基本使用
#### 编译
make build
#### 运行
extractor-server run -c examples/extractor-server.yaml
#### 配置
建议将工作线程数配置为 cpu核数-1
|参数|含义|类型|是否必需|
| ---- | ---- | ---- | ---- |
|worker.number|工作线程数|int|默认工作线程数为4|
|queue.number|队列中缓存请求数量|int|默认为4000，请求数量过多会返回429错误码|
|http.addr|地址和端口|string|默认端口7890|
## 基础功能
* 正则提链
* css提链
* xpath提链
* 只提取本站
* 自动去除锚点
* 提链黑白名单(domain, pattern)
* 支持xml、html格式文本
* 适配多种编码方式
* 默认提取html、htm、无后缀 链接 

## 请求格式

http+json

|参数|含义|类型|是否必需|
| ---- | ---- | ---- | ---- |
|URL|请求url|string|Y|
|ContentType|http header中的ContentType|string|Y|
|Content|需要提链的页面|string,使用base64编码|Y|
|OnlyHomeSite|是否只提取本站|bool(false)|N|
|IfRegexp|是否使用正则提链|bool(false)|N|
|CSSSelectors|css提链规则|string[]|N|
|XPathQuerys|xpath提链规则|string[]|N|
|AllowedDomains|domain白名单|string[]|N|
|DisallowedDomains|domain黑名单|string[]|N|
|AllowedURLFilters|pattern黑名单|string[]|N|
|DisallowedURLFilters|pattern黑名单|string[]|N|
|AllowedExts|被允许的后缀|string[]|N|

## 返回格式

|参数|含义|类型|
| ---- | ---- | ---- |
|re|正则提链结果|string[]|
|xpath|xpath提链结果|string[]|
|css|css提链结果|string[]|

## 性能(非正则提链)
非正则提链快于正则提链

测试页面 https://www.sohu.com/ 单请求平均耗时与客户端数量/工作线程数成正相关
|客户端数量/请求总量|工作线程数|qps|单请求平均耗时|
| ---- | ---- | ---- | ---- |
|100/10000|1|123.2335|0.8072 secs|
|100/10000|4|316.2894|0.3145 secs|

## 未来方向
* 抽取器

## License
MIT licensed.
