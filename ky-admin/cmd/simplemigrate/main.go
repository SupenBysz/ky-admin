package main

import (
	"github.com/SupenBysz/ky-admin/cmd/simplemigrate/utils"
)

func main() {
	// 调用utils包中的RunMigration函数执行数据库迁移
	utils.RunMigration()
}
