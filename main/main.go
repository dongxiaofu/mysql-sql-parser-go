package main

import(
	"fmt"
	"os"
	"io"
	"bufio"
	"mysql-sql-parser/stack"
	"reflect"
	"strings"
	"regexp"
	"errors"
	"bytes"
	"flag"
)

func main(){

	sqlFile := flag.String("sql", "", "sql文件名，必填")
	docFile := flag.String("doc", "example.md", "markdown文件名")

	flag.Parse()

	if *sqlFile == "" {
		fmt.Println("请输入要解析的SQL文件的文件名")
		return
	}

	fileName := *sqlFile
	file,err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error", err)
		return
	}


	if err != nil {
		panic(err)
	}

	var in bool = false
	var out bool = false

	var tables []Table

	buf := bufio.NewReader(file)
	myStack := Algorithm.NewStack(reflect.TypeOf(""))
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")

				break
			}else{
				fmt.Println("Read file error!", err)
			}
		}

		if(!in && isStart(line)){
			in = true
		}

		if(!out && isEnd(line)){
			out = true
		}

		if in {
			myStack.Push(line)

			if out {
				// 出栈
				myTableStack := Algorithm.NewStack(reflect.TypeOf(""))
				for myStack.Len() > 0 {
					line2, _ := myStack.Pop()
					if _,ok2 := line2.(string); ok2 {
						myTableStack.Push(line2)
					}
				}

				in = false
				out = false

				var sql []string
				// 再出栈

				for myTableStack.Len() > 0 {
					str, _ := myTableStack.Pop()
					if v,ok := str.(string);ok {
						sql = append(sql, v)
					}
				}

				table, _ := parseTable(sql)
				tables = append(tables, table)
			}
		}
	}

	// 生成文档
	doc := createDocument(tables)

	//  保存到文件中
	saveToFile(doc, *docFile)

	fmt.Println("解析 ", fileName, " 完成")
	fmt.Println("生成的markdown文件已经保存到 ", *docFile, " 中")
}

func isStart(str string) bool {
	var pattern string = `CREATE TABLE .*? \(`
	reg := regexp.MustCompile(pattern)
	res := reg.FindAllString(str, -1)

	return (res != nil)
}

func isEnd(str string) bool  {
	var search string = ") ENGINE="

	return strings.Contains(str, search)
}

func parseTableName(sql string) string  {

	reg := regexp.MustCompile("`(.*?)`")
	match := reg.FindStringSubmatch(sql)

	var len int = len(match)
	if len >= 1 {
		return match[1]
	}else{
		return ""
	}
}

func parseTableComment(sql string) string  {


	reg := regexp.MustCompile("COMMENT='(.*?)'")
	match := reg.FindStringSubmatch(sql)

	var len int = len(match)

	if(len >= 2) {
		return match[1]
	}else{
		return ""
	}
}

type TableColumn struct {
	name string
	comment string
}

func parseTableColumn(sql string) TableColumn  {
	s := strings.Split(sql, " ")

	var column TableColumn

	len := len(s)
	if len > 0 {
		column.name = strings.Replace(s[0], "`", "", -1)
		if (len - 2 >= 0) && (s[len - 2] == "COMMENT") {
			if len -1 >= 0 {
				column.comment = strings.Replace(s[len - 1], "'", "", -1)
			}else{
				column.comment = ""
			}
		}else{
			column.comment = ""
		}
	}else{
		column.name = ""
		column.comment = ""
	}

	return column
}

type Table struct {
	name string
	comment string
	column []TableColumn
}

func parseTable(sql []string) (Table, error)  {
	len := len(sql)
	if len <= 0 {
		return Table{}, errors.New("empty sql")
	}

	var table Table
	var column []TableColumn
	var tableName string
	var comment string

	// 遍历切片
	for i, str := range sql {
		if i == 0 {
			tableName = parseTableName(str)
		}else if i < len - 2 {
			if checkSqlIsIndex(str) == false {
				column = append(column, parseTableColumn(str))
			}
		}else{
			comment = parseTableComment(str)
		}
	}

	table.name = tableName

	table.comment = comment

	table.column = column

	return table, nil
}

func checkSqlIsIndex(sql string) bool  {
	var pattern string = "KEY.*?(.*?)"
	reg := regexp.MustCompile(pattern)
	res := reg.FindAllString(sql, -1)

	return res != nil
}

func createDocumentSegment(table Table) string {
	var buffer bytes.Buffer

	buffer.WriteString(table.name)
	buffer.WriteString("(")
	buffer.WriteString(table.comment)
	buffer.WriteString(")")
	buffer.WriteString("\n")

	buffer.WriteString("字段|描述")
	buffer.WriteString("\n")
	buffer.WriteString(":---|:---")
	buffer.WriteString("\n")

	columns := table.column

	for _, column := range columns {
		buffer.WriteString(column.name)
		buffer.WriteString("|")
		buffer.WriteString(column.comment)
		buffer.WriteString("\n")
	}

	doc := buffer.String()

	return doc
}

func createDocument(tables []Table) string {
	var buffer bytes.Buffer

	for _, table := range tables {
		docSegment := createDocumentSegment(table)
		buffer.WriteString(docSegment)
	}

	doc := buffer.String()

	return doc
}

func saveToFile(doc string, file string)  {
	dstFile, err := os.Create(file)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer dstFile.Close()

	dstFile.WriteString(doc)
}