package cmd

import (
	"db2doc/mysql"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"reflect"
)

// 数据库host地址
var host string

// 数据库端口
var port string

// 数据库用户名
var user string

// 数据库密码
var password string

// 数据库名称
var dbname string

// 表名称，支持使用下划线或百分号进行模糊匹配
var tablename string

// 输出的文件格式，为空则默认为md(markdown)
var fileType string

// 输出的文件名称
var fileName string

type TableColumn struct {
	// 列名
	COLUMN_NAME string
	// 字段类型
	COLUMN_TYPE string
	// 是否为空
	IS_NULLABLE string
	// 键
	COLUMN_KEY string
	// 默认值
	COLUMN_DEFAULT string
	// 备注
	COLUMN_COMMENT string
}

// 定义实现mysql2doc子命令
var mysql2docCmd = &cobra.Command{
	Use:   "mysql2doc",
	Short: "mysql数据库表结构文档生成工具",
	Long:  `该工具可以根据mysql表结构，生成对应的数据库文档，支持输出markdown格式`,
	Run: func(cmd *cobra.Command, args []string) {

		if password == "" {
			fmt.Println("密码不能为空")
			os.Exit(1)
		}

		if dbname == "" {
			fmt.Println("数据库名称不能为空")
			os.Exit(1)
		}

		if fileName == "" {
			fmt.Println("输出文件名称不能为空")
			os.Exit(1)
		}

		err := mysql.InitDB(host, port, dbname, user, password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 定义查询表结构的数据sql
		var findTableNameSql = "SELECT DISTINCT table_name FROM INFORMATION_SCHEMA.COLUMNS " +
			" WHERE table_schema = '" + dbname + "'"
		// 如果表名称不为空，则根据指定的表名称进行查询
		if tablename != "" {
			findTableNameSql = findTableNameSql + "AND table_name LIKE '" + tablename + "'"
		}

		resultTableNameArr, err := mysql.QueryList(findTableNameSql, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			// 遍历所有表名称，这里通过反射直接获取到切片长度，而后获取值
			tables := reflect.ValueOf(resultTableNameArr)
			// 创建一个切片，用来存储所有表的表名称
			tablesNames := make([]string, tables.Len())
			// 创建一个切片，用来存储所有表的表注释
			tablesComments := make([]string, tables.Len())
			// 创建一个切片，用来存储所有表的描述信息
			tablesColumnsComments := make([][]*TableColumn, tables.Len())
			// 开始遍历查询到的数据
			for i := 0; i < tables.Len(); i++ {
				// 获取每个表的表名称
				itemTableName := tables.Index(i).String()
				tablesNames[i] = itemTableName
				// 查询该表对应的注释内容
				var tableCommentSql = "SELECT table_comment FROM information_schema.TABLES " +
					" WHERE table_schema = '" + dbname + "'" + " AND table_name = '" + itemTableName + "'"
				tableCommentInfo, err := mysql.QueryInfo(tableCommentSql, "")
				if err != nil {
					fmt.Println("查询数据表注释出错", err)
					os.Exit(1)
				} else {
					// 获取到表注释
					tableComment := reflect.ValueOf(tableCommentInfo).String()
					tablesComments[i] = tableComment
				}

				// 开始查询该表对应的字段说明
				var columnNameSql = "SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, IF(COLUMN_DEFAULT IS NULL , 'NULL', COLUMN_DEFAULT) AS COLUMN_DEFAULT, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS " +
					"WHERE table_schema = '%s' AND table_name = '%s'"
				queryColumnNameSql := fmt.Sprintf(columnNameSql, dbname, itemTableName)
				resultArr, err := mysql.QueryList(queryColumnNameSql, (*TableColumn)(nil))
				tablesColumnsComments[i] = resultArr
				if err != nil {
					fmt.Println("查询数据表结构出错", err)
					os.Exit(1)
				}
			}
			// 如果需要写入到markdown文件中
			if fileType == "md" {
				writeToMd(tablesNames, tablesComments, tablesColumnsComments)
			}

		}
	},
}

/*
  将数据写入到markdown文件
*/
func writeToMd(tablesName []string, tablesComments []string, tablesColumnsComments [][]*TableColumn) {
	// 先创建要输出的文件
	outFile, err := os.Create(fileName + "." + fileType)
	if err != nil {
		fmt.Println("文件创建出错", err)
		os.Exit(1)
	}
	// 循环所有的表
	for i := 0; i < len(tablesColumnsComments); i++ {
		// 先输出当前表名称
		outFile.WriteString(tablesName[i] + " ")
		// 再输出当前表注释
		outFile.WriteString("(" + tablesComments[i] + ")")
		// 最后输出当前表字段描述
		outFile.WriteString("\r\n")
		outFile.WriteString("| 列名称 | 字段类型 | 是否为空 | 键 | 默认值 | 备注 |")
		outFile.WriteString("\r\n")
		outFile.WriteString("| ----- | -------- | ------- | -- | ----- | ---- |")
		outFile.WriteString("\r\n")
		for j := 0; j < len(tablesColumnsComments); j++ {
			columnName := tablesColumnsComments[i][j].COLUMN_NAME
			columnType := tablesColumnsComments[i][j].COLUMN_TYPE
			nullable := tablesColumnsComments[i][j].IS_NULLABLE
			columnKey := tablesColumnsComments[i][j].COLUMN_KEY
			columnDefault := tablesColumnsComments[i][j].COLUMN_DEFAULT
			comment := tablesColumnsComments[i][j].COLUMN_COMMENT
			columnLine := fmt.Sprintf("| %s | %s | %s | %s | %s | %s |", columnName, columnType, nullable, columnKey, columnDefault, comment)
			outFile.WriteString(columnLine)
			outFile.WriteString("\r\n")
			defer outFile.Close()
		}
		// 每输出完毕一张表，都进行下换行
		outFile.WriteString("\r\n")
	}

}

func init() {
	rootCmd.AddCommand(mysql2docCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mysql2docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mysql2docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	mysql2docCmd.Flags().StringVarP(&host, "host", "H", "localhost", "数据库host地址 默认为本机地址")
	mysql2docCmd.Flags().StringVarP(&port, "port", "P", "3306", "数据库端口 默认为3306")
	mysql2docCmd.Flags().StringVarP(&user, "user", "u", "root", "数据库用户名 默认为root")
	mysql2docCmd.Flags().StringVarP(&password, "password", "p", "", "数据库密码")
	mysql2docCmd.Flags().StringVarP(&dbname, "dbname", "d", "", "数据库名称")
	mysql2docCmd.Flags().StringVarP(&tablename, "tablename", "t", "", "表名称，支持使用下划线或百分号进行模糊匹配，不传入则生成所有表")
	mysql2docCmd.Flags().StringVarP(&fileType, "filetype", "T", "md", "输出的文件格式，md为markdown，不输入则默认为markdown格式")
	mysql2docCmd.Flags().StringVarP(&fileName, "filename", "f", "", "输出的文件名称")
}
