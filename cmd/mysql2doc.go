package cmd

import (
	"github.com/spf13/cobra"
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

// 定义实现mysql2doc子命令
var mysql2docCmd = &cobra.Command{
	Use:   "mysql2doc",
	Short: "mysql数据库表结构文档生成工具",
	Long:  "该工具可以根据mysql表结构，生成对应的数据库文档，支持输出markdown和word两种格式",
	Run: func(cmd *cobra.Command, args []string) {

	},
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
