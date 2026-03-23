package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/leanote/leanote/tools/migration/core"
	"github.com/leanote/leanote/tools/migration/mongo2pg"
	"github.com/leanote/leanote/tools/migration/pg2mongo"
)

func main() {
	// 解析命令行参数
	configFile := flag.String("config", "config/migration.yaml", "配置文件路径")
	direction := flag.String("direction", "", "迁移方向: mongo_to_pg 或 pg_to_mongo")
	dryRun := flag.Bool("dry-run", false, "试运行模式")
	flag.Parse()

	if *direction == "" {
		log.Fatal("请指定迁移方向: -direction mongo_to_pg 或 -direction pg_to_mongo")
	}

	// 加载配置
	config, err := core.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	config.Direction = *direction
	config.DryRun = *dryRun

	// 创建迁移器
	var migrator core.Migrator

	switch config.Direction {
	case "mongo_to_pg":
		migrator = mongo2pg.NewMigrator(config)
	case "pg_to_mongo":
		migrator = pg2mongo.NewMigrator(config)
	default:
		log.Fatalf("未知的迁移方向: %s", config.Direction)
	}

	// 执行迁移
	fmt.Printf("开始迁移: %s\n", config.Direction)
	fmt.Printf("源数据库: %s\n", config.Source.Database)
	fmt.Printf("目标数据库: %s\n", config.Target.Database)
	fmt.Printf("试运行: %v\n", config.DryRun)
	fmt.Println("----------------------------------------")

	startTime := time.Now()

	stats, err := migrator.Migrate()

	duration := time.Since(startTime)
	if err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	// 输出统计信息
	fmt.Println("----------------------------------------")
	fmt.Printf("迁移完成! 耗时: %v\n", duration)
	fmt.Printf("总表数: %d\n", stats.TotalTables)
	fmt.Printf("成功表数: %d\n", stats.SuccessTables)
	fmt.Printf("失败表数: %d\n", stats.FailedTables)
	fmt.Printf("跳过表数: %d\n", stats.SkippedTables)
	fmt.Printf("总记录数: %d\n", stats.TotalRecords)
	fmt.Printf("成功记录数: %d\n", stats.SuccessRecords)
	fmt.Printf("失败记录数: %d\n", stats.FailedRecords)
}
