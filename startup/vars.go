package startup

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"wechat-robot-client/vars"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupVars() error {
	if err := InitMySQLClient(); err != nil {
		return fmt.Errorf("MySQL连接失败: %v", err)
	}
	log.Println("MySQL连接成功")
	if err := InitRedisClient(); err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}
	log.Println("Redis连接成功")
	return nil
}

func InitMySQLClient() (err error) {
	// 获取本机IP地址
	localIP := getLocalIP()
	log.Printf("[DEBUG] 本机IP地址: %s", localIP)

	// 检查网络连接
	log.Printf("[DEBUG] 检查网络连接到 %s:%s", vars.MysqlSettings.Host, vars.MysqlSettings.Port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", vars.MysqlSettings.Host, vars.MysqlSettings.Port))
	if err != nil {
		log.Printf("[ERROR] 无法连接到MySQL服务器 %s:%s: %v", vars.MysqlSettings.Host, vars.MysqlSettings.Port, err)
		return fmt.Errorf("网络连接失败: %v", err)
	}
	conn.Close()
	log.Printf("[DEBUG] 网络连接测试成功")

	// 创建机器人实例连接对象
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		vars.MysqlSettings.User, vars.MysqlSettings.Password, vars.MysqlSettings.Host, vars.MysqlSettings.Port, vars.MysqlSettings.Db)

	// 打印DSN信息（隐藏密码）
	maskedDsn := fmt.Sprintf("%s:***@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		vars.MysqlSettings.User, vars.MysqlSettings.Host, vars.MysqlSettings.Port, vars.MysqlSettings.Db)
	log.Printf("[DEBUG] 尝试连接MySQL数据库: %s", maskedDsn)
	mysqlConfig := mysql.Config{
		DSN:                     dsn,
		DontSupportRenameIndex:  true, // 重命名索引时采用删除并新建的方式
		DontSupportRenameColumn: true, // 用 `change` 重命名列
	}
	// gorm 配置
	gormConfig := gorm.Config{}
	// 是否开启调试模式
	if flag, _ := strconv.ParseBool(os.Getenv("GORM_DEBUG")); flag {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	vars.DB, err = gorm.Open(mysql.New(mysqlConfig), &gormConfig)
	if err != nil {
		return err
	}

	// 创建机器人管理后台连接对象
	adminDsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		vars.MysqlSettings.User, vars.MysqlSettings.Password, vars.MysqlSettings.Host, vars.MysqlSettings.Port, vars.MysqlSettings.AdminDb)

	// 打印管理后台DSN信息（隐藏密码）
	adminMaskedDsn := fmt.Sprintf("%s:***@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		vars.MysqlSettings.User, vars.MysqlSettings.Host, vars.MysqlSettings.Port, vars.MysqlSettings.AdminDb)
	log.Printf("[DEBUG] 尝试连接MySQL管理后台数据库: %s", adminMaskedDsn)
	adminYysqlConfig := mysql.Config{
		DSN:                     adminDsn,
		DontSupportRenameIndex:  true, // 重命名索引时采用删除并新建的方式
		DontSupportRenameColumn: true, // 用 `change` 重命名列
	}
	// gorm 配置
	adminGormConfig := gorm.Config{}
	// 是否开启调试模式
	if flag, _ := strconv.ParseBool(os.Getenv("GORM_DEBUG")); flag {
		adminGormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	vars.AdminDB, err = gorm.Open(mysql.New(adminYysqlConfig), &adminGormConfig)
	return err
}

func InitRedisClient() (err error) {
	vars.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", vars.RedisSettings.Host, vars.RedisSettings.Port),
		Password: vars.RedisSettings.Password,
		DB:       vars.RedisSettings.Db,
	})
	_, err = vars.RedisClient.Ping(context.Background()).Result()
	return err
}

// getLocalIP 获取本机IP地址
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
