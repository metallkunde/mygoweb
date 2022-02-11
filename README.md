# mygoweb
框架中的很多部分实现的功能都很简单，但是尽可能地体现一个框架核心的设计原则。
例如Router的设计，虽然支持的动态路由规则有限，但为了性能考虑匹配算法是用Trie树实现的，Router最重要的指标之一便是性能。

day1：前置基础知识 http.Handler  
day2：上下文设计 Context  
day3：前缀树路由 Router  
day4：分组控制 Group  
day5：中间件 Middleware  
day6：HTML模板 Template  
day7：错误恢复 Panic Recover
