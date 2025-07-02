package verify

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"checker-middleware/pkg/logger"
	"strings"
	"time"

	_ "gitee.com/chunanyong/dm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// 转换为Go的数据库连接字符串
func convertToRdbConnection(cfg RDBConfig, timeout int) RDBConnection {
	// 生成链接字符串
	var dsn string
	driver := strings.ToLower(cfg.Driver)
	switch driver {
	case "mysql", "goldendb", "mariadb", "tdsql", "oceanbase":
		// root:pass@tcp(127.0.0.1:3306)/dbname
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	case "dm", "dameng":
		// dm://user:pwd@127.0.0.1:5236
		driver = "dm"
		dsn = fmt.Sprintf("dm://%s:%s@%s:%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	case "postgresql", "postgres", "pg", "pgsql":
		// host=127.0.0.1 port=5432 user=postgres password=123456 dbname=postgres
		driver = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	default:
		dsn = "unsupported"
	}
	logger.DebugLog("convertToRdbConnection: driver=%s, dsn=%s, timeout=%d", driver, dsn, timeout)
	return RDBConnection{
		Driver:   driver,
		DSN:      dsn,
		TimeoutS: timeout,
	}
}

// 根据 driver 返回建表、插入、删表 SQL
func getTestSQL(driver string) (createSQL, insertSQL, dropSQL string) {
	switch strings.ToLower(driver) {
	case "mysql", "goldendb", "mariadb", "tdsql", "oceanbase":
		return `CREATE TABLE IF NOT EXISTS precheck_test(id INT PRIMARY KEY, val VARCHAR(32))`,
			`INSERT INTO precheck_test(id, val) VALUES (1, 'ok')`,
			`DROP TABLE IF EXISTS precheck_test`
	case "dm", "dameng":
		return `CREATE TABLE precheck_test(id INT PRIMARY KEY, val VARCHAR(32))`,
			`INSERT INTO precheck_test(id, val) VALUES (1, 'ok')`,
			`DROP TABLE precheck_test`
	case "postgresql", "postgres", "pg", "pgsql":
		return `CREATE TABLE IF NOT EXISTS precheck_test(id INT PRIMARY KEY, val VARCHAR(32))`,
			`INSERT INTO precheck_test(id, val) VALUES (1, 'ok')`,
			`DROP TABLE IF EXISTS precheck_test`
	default:
		return "", "", ""
	}
}

func openDB(cfg RDBConfig, timeout int) (*sql.DB, string, error) {
	conn := convertToRdbConnection(cfg, timeout)
	driver := strings.ToLower(conn.Driver)
	db, err := sql.Open(driver, conn.DSN)
	if err != nil {
		return nil, driver, err
	}
	db.SetConnMaxLifetime(time.Duration(timeout) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(timeout) * time.Second)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, driver, nil
}

// 通用数据库连接测试
func RdbConnect(config RDBConfig) map[string]string {
	result := map[string]string{"success": "false"}
	db, driver, err := openDB(config, 10)
	if err != nil {
		result["error"] = fmt.Sprintf("open error: %v", err)
		return result
	}
	defer db.Close()
	done := make(chan error, 1)
	go func() {
		logger.DebugLog("RdbConnect: Start db.Ping()")
		done <- db.Ping()
	}()
	select {
	case err := <-done:
		if err != nil {
			result["error"] = fmt.Sprintf("ping error: %v", err)
		} else {
			result["success"] = "true"
		}
	case <-time.After(10 * time.Second):
		result["error"] = "ping timeout"
	}
	logger.DebugLog("RdbConnect: driver=%s, result=%v", driver, result)
	return result
}

// 通用数据库写入测试
func RdbWrite(config RDBConfig, createSQL, insertSQL string) map[string]string {
	result := map[string]string{"success": "false"}
	db, _, err := openDB(config, 10)
	if err != nil {
		result["error"] = fmt.Sprintf("open error: %v", err)
		return result
	}
	defer db.Close()
	logger.DebugLog("RdbWrite: createSQL=%s, insertSQL=%s", createSQL, insertSQL)
	if _, err := db.Exec(createSQL); err != nil {
		result["error"] = fmt.Sprintf("create error: %v", err)
		return result
	}
	logger.DebugLog("RdbWrite: createSQL precheck_test table success!")
	if _, err := db.Exec(insertSQL); err != nil {
		result["error"] = fmt.Sprintf("insert error: %v", err)
		return result
	}
	logger.DebugLog("RdbWrite: insertSQL Success!")
	result["success"] = "true"
	return result
}

// 通用数据库删除测试
func RdbDelete(config RDBConfig, dropSQL string) map[string]string {
	result := map[string]string{"success": "false"}
	db, _, err := openDB(config, 10)
	if err != nil {
		result["error"] = fmt.Sprintf("open error: %v", err)
		return result
	}
	defer db.Close()
	if _, err := db.Exec(dropSQL); err != nil {
		result["error"] = fmt.Sprintf("drop error: %v", err)
		return result
	}
	logger.DebugLog("RdbDelete: dropSQL precheck_test table success!")
	result["success"] = "true"
	return result
}

// 一键检测并返回RDBResult
func VerifyRDB(config RDBConfig) RDBResult {
	conn := convertToRdbConnection(config, 10)
	createSQL, insertSQL, dropSQL := getTestSQL(conn.Driver)
	res := RDBResult{
		Connect: RdbConnect(config),
		Write:   map[string]string{"success": "skip"},
		Delete:  map[string]string{"success": "skip"},
	}
	if res.Connect["success"] == "true" && createSQL != "" {
		res.Write = RdbWrite(config, createSQL, insertSQL)
		res.Delete = RdbDelete(config, dropSQL)
	}
	return res
}

// 返回json
func VerifyRDBJson(config RDBConfig) []byte {
	res := VerifyRDB(config)
	b, _ := json.Marshal(res)
	return b
}
