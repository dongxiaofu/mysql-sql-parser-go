# mysql-sql-parser-go
解析mysql sql 文件，生成markdown数据词典

main/main.go 是实现功能的代码

main/main_test.go 是单元测试代码

stack 中的代码是实现栈的代码，来源于网络

bin/mysql-sql-parser 是编译后的可执行程序

用法：mysql-sql-parser --sql=example.sql --doc=example.md
