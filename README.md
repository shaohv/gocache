# gocache
a cache server written in golang.   
参考《分布式缓存-原理、架构及go语言实现》，用于学习golang编程

项目中用到的技术：

1. git submodule add remoteaddr localdir

2. pipeline技术 

3. rocksdb LSM树

4. 时间轮盘 高性能定时器

5. 一致性协议

6. gossip协议

7. 常用linux 命令: strace, tcpdump
8. 服务端解偶：将处理网络和数据库操作放到不同的goroute中