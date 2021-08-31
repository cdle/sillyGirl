
目前支持 `windows amd64`/`liunx amd64`/`darwin arm64`二进制双击运行安装

## Docker
```bash
# 在你要存放数据的目录下手动新建sillyGirl文件夹
# (以root目录为例)
# 警告！群晖用户请勿在root下存放任何文件！修改成你的硬盘目录！
mkdir /root/sillyGirl    #在root目录新建sillyGirl文件夹

# 拉取并运行容器 并进入交互控制台
docker run -itd \
-v /root/sillyGirl/data:/etc/sillyplus \
-v /root/sillyGirl/plugins:/usr/local/sillyplus/plugins \
-p 8080:8080 \
--name sillyGirl \
--restart always \
jackytj/sillyplus && docker attach sillyGirl

```
进入容器交互控制台
```bash
#进入
docker attach sillyGirl
# 退出交互控制台
Ctrl+p Ctrl+q
```

查看日志
```bash
docker logs sillyGirl
```

进入容器命令行(一般用不到)
```bash
docker exec -it sillyGirl /bin/sh
```


## Windows
下载双击运行

 - [sillyGirl_windows_amd64.exe](https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_windows_amd64.exe)


## MacOS

敬请期待

## Liunx

下载双击运行

 - [sillyGirl_linux_amd64](https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_amd64)

 - [sillyGirl_linux_arm64](https://github.com/cdle/sillyGirl/releases/download/main/sillyGirl_linux_arm64)

## Other

 - 如果你的运行环境不受支持

 - 如果你在使用上有任何困难

 - 可以加入本项目的[知识星球](https://wx.zsxq.com/dweb2/index/group/28885424215821)，以获取手拉手教程、特殊资源共享、定制需求和同作者实时交流