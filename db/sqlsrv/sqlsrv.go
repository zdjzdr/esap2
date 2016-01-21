// db/sqlsrv for ESAP
// By woylin 2015.10.20
package sqlsrv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	_ "github.com/denisenkom/go-mssqldb"
)

//数据库连接池
var (
	db *sql.DB
	dc *DbConf
)

type DbConf struct {
	UserId string
	Pwd    string
	Server string
	DbName string
}

//checkSingle
func CheckBool(sql string, cond ...interface{}) bool {
	checkDB()
	rs, err := db.Query(sql, cond...)
	checkErr(err)
	if !rs.Next() {
		return false
	}
	return true
}

//fetchOne
func FetchOne(query string, cond ...interface{}) *[]interface{} {
	checkDB()
	row := db.QueryRow(query, cond...)
	result := make([]interface{}, 0)
	onerow := make([]interface{}, 2)
	err := row.Scan(onerow...)
	if err != nil {
		panic(err)
	}
	result = append(result, onerow)
	return &result
}

//fetchAll
func FetchAll(query string, cond ...interface{}) *[]interface{} {
	checkDB()
	rows, err := db.Query(query, cond...)
	checkErr(err)
	defer rows.Close()
	cols, err := rows.Columns()
	checkErr(err)
	leng := len(cols)
	result := make([]interface{}, 0)      //结果集，数组
	scanArgs := make([]interface{}, leng) //扫描专用指针
	onerow := make([]interface{}, leng)   //数据行，无字段名
	for i := range onerow {
		scanArgs[i] = &onerow[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			continue
		}
		data := make(map[string]interface{}) //数据行，含字段名
		for i, _ := range onerow {
			data[cols[i]] = conv(onerow[i])
		}
		result = append(result, data)
	}
	return &result
}

type treeNode struct {
	Id       interface{}   `json:"id"`
	Text     interface{}   `json:"text"`
	Expanded bool          `json:"expanded"`
	Leaf     bool          `json:"leaf"`
	Children []interface{} `json:"children"`
}

func (t *treeNode) appendChild(c interface{}) {
	t.Children = append(t.Children, c)
}

//fetchTree
func FetchMenuTree(query string, cond ...interface{}) *treeNode {
	checkDB()
	rows, err := db.Query(query, cond...)
	checkErr(err)
	defer rows.Close()
	cols, err := rows.Columns()
	checkErr(err)
	leng := len(cols)
	scanArgs := make([]interface{}, leng) //扫描专用指针
	onerow := make([]interface{}, leng)   //数据行，无字段名
	for i := range onerow {
		scanArgs[i] = &onerow[i]
	}
	treeMap := make(map[string]*treeNode, 0)
	tree := &treeNode{1, "root", true, false, make([]interface{}, 0)}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			continue
		}
		data := make(map[string]interface{}) //数据行，含字段名
		for i, _ := range onerow {
			data[cols[i]] = conv(onerow[i])
		}
		menuName := fmt.Sprintf("%s", data["menu"])
		if _, ok := treeMap[menuName]; !ok {
			treeMap[menuName] = &treeNode{data["ordPath"], data["menu"], true, false, make([]interface{}, 0)}
			tree.appendChild(treeMap[menuName])
		}
		treeMap[menuName].appendChild(treeNode{data["id"], data["name"], false, true, nil})
		tree.appendChild(treeMap[menuName])

	}
	return tree
}

//fetch
func Fetch(query string, cond ...interface{}) *interface{} {
	checkDB()
	var rst interface{}
	db.QueryRow(query, cond...).Scan(&rst)
	return &rst
}

//NumRows
func NumRows(query string, cond ...interface{}) int {
	checkDB()
	rows, err := db.Query(query, cond...)
	checkErr(err)
	defer rows.Close()
	result := 0
	for rows.Next() {
		result++
	}
	return result
}

//Exec
func Exec(query string, cond ...interface{}) error {
	checkDB()
	stmt, err := db.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(cond...)
	if err != nil {
		return err
	}
	return nil
}

//检查DB是否连接，无则进行连接
func checkDB() {
	if db == nil {
		if dc == nil {
			bytes, err := ioutil.ReadFile("conf/db.json")
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(bytes, &dc); err != nil {
				panic(err)
			}
		}
		linkDb()
	}
}

//更改DB
func ChangeDb(s ...string) {
	if len(s) == 4 {
		dc = &DbConf{s[0], s[1], s[2], s[3]}
		linkDb()
	}
}

//连接DB
func linkDb() {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", dc.Server, dc.UserId, dc.Pwd, dc.DbName)
	db1, err := sql.Open("mssql", dsn)
	checkErr(err)
	db = db1
}

//通用查询
func FetchAllRowsPtr(query string, struc interface{}, cond ...interface{}) *[]interface{} {
	checkDB()
	rows, err := db.Query(query, cond...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := make([]interface{}, 0)
	s := reflect.ValueOf(struc).Elem()
	leng := s.NumField()
	onerow := make([]interface{}, leng)
	for i := 0; i < leng; i++ {
		onerow[i] = s.Field(i).Addr().Interface()
	}
	for rows.Next() {
		err = rows.Scan(onerow...)
		if err != nil {
			panic(err)
		}
		result = append(result, s.Interface())
	}
	return &result
}

//通用查询单条
func FetchOnePtr(query string, struc interface{}, cond ...interface{}) *interface{} {
	checkDB()
	rows, err := db.Query(query, cond...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	s := reflect.ValueOf(struc).Elem()
	leng := s.NumField()
	onerow := make([]interface{}, leng)
	for i := 0; i < leng; i++ {
		onerow[i] = s.Field(i).Addr().Interface()
	}
	if rows.Next() {
		err = rows.Scan(onerow...)
		if err != nil {
			panic(err)
		}
	}
	result := s.Interface()
	return &result
}
