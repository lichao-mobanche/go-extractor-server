# go-extractor-server
## 基本使用
#### 编译
make build
#### 运行
extractor-server run -c examples/extractor-server.yaml
## 基础功能
* 正则提链
* css提链
* xpath提链
* 只提取本站
* 提链黑白名单(domain, pattern)
* 支持xml、html格式文本
* 适配多种编码方式

## 请求格式

|参数|含义|类型|是否必需|
| ---- | ---- | ---- | ---- |
|URL|请求url|string|Y|
|ContentType|http header中的ContentTyoe|string|Y|
|Content|需要提链的页面|string,使用base64编码|Y|
|OnlyHomeSite|是否只提取本站|bool(false)|N|
|IfRegexp|是否使用正则提链|bool(false)|N|
|CSSSelectors|css提链规则|string[]|N|
|XPathQuerys|xpath提链规则|string[]|N|
|AllowedDomains|domain白名单|string[]|N|
|DisallowedDomains|domain黑名单|string[]|N|
|AllowedURLFilters|pattern黑名单|string[]|N|
|DisallowedURLFilters|pattern黑名单|string[]|N|

## 返回格式

|参数|含义|类型|
| ---- | ---- | ---- |
|re|正则提链结果|string[]|
|xpath|xpath提链结果|string[]|
|css|css提链结果|string[]|
|re|正则提链结果|string[]|
|re|正则提链结果|string[]|

## 性能
测试中

## 未来方向
* 抽取器
