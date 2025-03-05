package gamedatabase

import (
	"ThreeKingdoms/src/config"
	"fmt"
	"log"
	"time"
	"xorm.io/xorm"

	// 导入 MySQL 驱动，注意这里使用了匿名导入
	_ "github.com/go-sql-driver/mysql"
)

var Engine *xorm.Engine

func TestDatabase() {
	// DSN 格式: username:password@protocol(address)/dbname?param=value
	//// 示例中使用 TCP 连接到本地 MySQL 数据库，记得替换为实际的用户名、密码、数据库名等信息
	//dsn := "root:loveyou@tcp(8.138.106.163:3306)/electric_dispatch?parseTime=true"
	mysqlConfig := config.Config.MySqlSection
	dsn := mysqlConfig.Dsn
	fmt.Println(dsn)
	engine, err := xorm.NewEngine("mysql", dsn)
	Engine = engine
	// 配置连接池参数
	engine.SetMaxOpenConns(mysqlConfig.MaxOpenConns)                     // 最大打开连接数
	engine.SetMaxIdleConns(mysqlConfig.MaxIdleConns)                     // 最大空闲连接数
	engine.SetConnMaxLifetime(mysqlConfig.MaxConnLifetime * time.Minute) // 连接最大存活时间，例如 30 分钟

	// 测试数据库连接
	err = engine.Ping()
	type XormUser struct {
		Id        int64
		Name      string
		Salary    float64
		Age       int64
		Password  string    `xorm:"varchar(255)"`
		CreatedAt time.Time `xorm:"created"`
		UpdatedAt time.Time `xorm:"updated"`
	}
	//user1 := XormUser{Age: 1, Name: "胡伟切片版本", Salary: 12}
	//user2 := XormUser{Age: 11, Name: "谢祥龙切片版本", Salary: 999999}
	//user3 := XormUser{Age: 32, Name: "李强切片版本", Salary: 333}
	//
	//sliceUser := []XormUser{user1, user2, user3}
	//userx := XormUser{
	//	Name: "赵信",
	//	Age:  999,
	//

	//engine.Iterate(&XormUser{Age: 999},
	//	func(idx int, bean interface{}) error {
	//
	//		user, ok := bean.(*XormUser)
	//		if ok {
	//			fmt.Println(user.Name)
	//		}
	//		return nil
	//	})

	rows, _ := engine.Rows(&XormUser{
		Age: 999})

	defer rows.Close()

	userBean := &XormUser{}

	for rows.Next() {
		err := rows.Scan(userBean)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(*userBean)
	}

	if err != nil {
		fmt.Println(err)
	}

}

// 创建 users 表
//func createTable(db *sql.DB) {
//	query := `CREATE TABLE IF NOT EXISTS users_go (
//		id INT AUTO_INCREMENT PRIMARY KEY,
//		name VARCHAR(100) NOT NULL,
//		age INT NOT NULL,
//		email VARCHAR(100) UNIQUE NOT NULL,
//		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//	);`
//	_, err := db.Exec(query)
//	if err != nil {
//		log.Fatalf("创建表失败: %v", err)
//	}
//	fmt.Println("表 users 创建成功（如果不存在）")
//}
//
//// 插入用户
//func insertUser(db *sql.DB, name string, age int, email string) int64 {
//	query := "INSERT INTO users (name, age, email) VALUES (?, ?, ?)"
//	result, err := db.Exec(query, name, age, email)
//	if err != nil {
//		log.Fatalf("插入用户失败: %v", err)
//	}
//	id, _ := result.LastInsertId()
//	return id
//}
//
//// 查询单个用户
//func getUser(db *sql.DB, id int64) {
//	query := "SELECT id, name, age, email, created_at FROM users WHERE id = ?"
//	var user struct {
//		ID        int64
//		Name      string
//		Age       int
//		Email     string
//		CreatedAt string
//	}
//	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Age, &user.Email, &user.CreatedAt)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			fmt.Println("用户不存在")
//		} else {
//			log.Fatalf("查询用户失败: %v", err)
//		}
//		return
//	}
//	fmt.Printf("用户信息: ID=%d, Name=%s, Age=%d, Email=%s, CreatedAt=%s\n",
//		user.ID, user.Name, user.Age, user.Email, user.CreatedAt)
//}
//
//// 查询所有用户
//func getAllUsers(db *sql.DB) {
//	query := "SELECT id, name, age, email, created_at FROM users"
//	rows, err := db.Query(query)
//	if err != nil {
//		log.Fatalf("查询所有用户失败: %v", err)
//	}
//	defer rows.Close()
//
//	fmt.Println("所有用户信息:")
//	for rows.Next() {
//		var id int64
//		var name string
//		var age int
//		var email string
//		var createdAt string
//		if err := rows.Scan(&id, &name, &age, &email, &createdAt); err != nil {
//			log.Fatalf("扫描用户数据失败: %v", err)
//		}
//		fmt.Printf("ID=%d, Name=%s, Age=%d, Email=%s, CreatedAt=%s\n",
//			id, name, age, email, createdAt)
//	}
//}
//
//// 更新用户
//func updateUser(db *sql.DB, id int64, name string, age int) {
//	query := "UPDATE users SET name = ?, age = ? WHERE id = ?"
//	_, err := db.Exec(query, name, age, id)
//	if err != nil {
//		log.Fatalf("更新用户失败: %v", err)
//	}
//	fmt.Printf("用户 ID=%d 已更新\n", id)
//}
//
//// 删除用户
//func deleteUser(db *sql.DB, id int64) {
//	query := "DELETE FROM users WHERE id = ?"
//	_, err := db.Exec(query, id)
//	if err != nil {
//		log.Fatalf("删除用户失败: %v", err)
//	}
//	fmt.Printf("用户 ID=%d 已删除\n", id)
//}
