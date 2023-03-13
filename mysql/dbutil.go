package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
)

var db *sql.DB

/*
  初始化数据库
*/
func InitDB(host, port, dbname, username, password string) (err error) {
	if host == "" {
		return errors.New("数据库地址参数为空")
	}
	if port == "" {
		return errors.New("数据库端口参数为空")
	}
	if dbname == "" {
		return errors.New("数据库名称参数为空")
	}
	if username == "" {
		return errors.New("数据库用户名参数为空")
	}
	if password == "" {
		return errors.New("数据库密码参数为空")
	}
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

/*
  查询单一数据
*/
func QueryInfo(sql string, src any) (resultInfo any, err error) {
	typeOf := reflect.TypeOf(src)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("QueryInfo查询sql时出错，错误原因:", err, "sql为", sql)
		return nil, err
	} else {
		switch typeOf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for rows.Next() {
				var data int
				err := rows.Scan(&data)
				if err != nil {
					fmt.Println("QueryInfo结果数据遍历出错，错误原因:", err)
				} else {
					return data, nil
				}
			}
		case reflect.String:
			for rows.Next() {
				var data string
				err := rows.Scan(&data)
				if err != nil {
					fmt.Println("QueryList结果数据遍历出错，错误原因:", err)
				} else {
					return data, nil
				}
			}
		case reflect.Ptr:
			//fmt.Println("需要通过反射实例化struct")
			//typeOf = typeOf.Elem()
			//newStruct := reflect.New(typeOf)
			//newStruct.Elem().FieldByName("Id").SetInt(100)
			//newStruct.Elem().FieldByName("Name").SetString("nameTest")
		}
		return nil, nil
	}
}

/*
  查询数据列表
*/
func QueryList[T any](sql string, src T) (resultTableNameArr []T, err error) {
	// 通过反射获取类型
	typeOf := reflect.TypeOf(src)

	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("QueryList查询sql时出错，错误原因:", err, "sql为", sql)
		return nil, err
	} else {
		switch typeOf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var resultArr = []T{}
			for rows.Next() {
				var data T
				err := rows.Scan(&data)
				if err != nil {
					fmt.Println("QueryList结果数据遍历出错，错误原因:", err)
				} else {
					resultArr = append(resultArr, data)
				}
			}
			return resultArr, nil
		case reflect.String:
			var resultArr = []T{}
			for rows.Next() {
				var data T
				err := rows.Scan(&data)
				if err != nil {
					fmt.Println("QueryList结果数据遍历出错，错误原因:", err)
				} else {
					resultArr = append(resultArr, data)
				}
			}
			return resultArr, nil
		case reflect.Ptr:
			var resultArr []T
			// 得到sql查询的所有的列名称
			columns, err := rows.Columns()
			// 先创建一个切片用来保存列的名称
			var columnNames = make([]string, len(columns))
			copy(columnNames, columns)

			// 创建一个interface切片，用来存放columns中每个元素的索引
			columnsIndex := make([]interface{}, len(columns))
			for i := range columnsIndex {
				// 得到columns的指针，然后赋值给columnsIndex，通过指针变更columns的值
				columnsIndex[i] = &columns[i]
			}

			for rows.Next() {

				if err != nil {
					fmt.Println("reflect.Ptr查询数据出错，错误原因:", err)
					os.Exit(1)
				}

				scanError := rows.Scan(columnsIndex...)
				if scanError != nil {
					fmt.Println("reflect.Ptr查询数据遍历出错，错误原因:", err)
				} else {
					// 将数据进行处理，同一行数据，将列名称与列数值对应起来，通过map存放
					rowMap := make(map[string]interface{})

					for i := 0; i < len(columns); i++ {
						rowMap[columnNames[i]] = columns[i]
					}

					// 遍历每一列的数据，将数据组装成struct而后返回
					typeOfElem := typeOf.Elem() // 这里是struct类型路径，例如cmd.TableColumn
					// 实例化新的struct
					var newStructInter T
					newStruct := reflect.New(typeOfElem)
					// 得到传入的struct的字段数量
					numField := typeOfElem.NumField()
					// 得到传入的struct的字段数量，按照查询sql的column的顺序进行拼接
					for i := 0; i < numField; i++ {
						// 得到字段名称，根据名称和规则进行值设置
						columnName := typeOfElem.Field(i).Name
						// 遍历原始的column，按照字段名称设置对应的值，这里如果往后细化，还需要判断Type类型
						// 因为目前使用的结构体全部存放string即可，不需要像比较完善的orm一样，因此就不判断Type了
						value := rowMap[columnName]
						newStruct.Elem().FieldByName(columnName).SetString(value.(string))
						newStructInter = newStruct.Interface().(T)
					}
					// 将创建好的struct放入数组内
					resultArr = append(resultArr, newStructInter)
				}
			}
			return resultArr, nil
		}
		return nil, nil
	}

}
