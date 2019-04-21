package main

import "testing"

// 运行单元测试 go test

func Test_parseTableName(t *testing.T)  {
	var sql string = "CREATE TABLE `cg_passage` ("
	if name := parseTableName(sql); name != "cg_passage" {
		t.Error("parseTableName测试没通过")
	}else{
		t.Log("parseTableName测试通过")
	}
}

func Test_parseTableComment(t *testing.T)  {
	var sql string = ") ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='友情链接';"
	if comment := parseTableComment(sql); comment != "友情链接" {
		t.Error("parseTableComment测试没通过")
	}else{
		t.Log("parseTableComment测试通过")
	}

	var sql2 string = ") ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;"
	if comment := parseTableComment(sql2); comment != "" {
		t.Error("parseTableComment测试没通过")
	}else{
		t.Log("parseTableComment测试通过")
	}
}
