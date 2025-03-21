package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/repository"
)

// TestErrorPackage 测试错误包是否正常工作
func TestErrorPackage(t *testing.T) {
	// 创建一个普通错误
	err := errors.NewUserError(errors.ErrUserNotFound)
	if err == nil {
		t.Error("应该返回一个错误")
	}

	// 检查错误类型
	if !errors.Is(err.Err, errors.ErrUserNotFound) {
		t.Error("错误类型识别失败")
	}

	// 测试仓库错误
	repoErr := repository.ErrUserNotFound
	if repoErr == nil {
		t.Error("仓库错误应该不为空")
	}
}

// TestProjectStructure 测试项目结构是否正确
func TestProjectStructure(t *testing.T) {
	// 获取当前路径
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	// 获取项目根目录
	rootDir := filepath.Dir(workDir)

	// 定义需要检查的目录
	dirsToCheck := []string{
		"internal/common/errors",
		"internal/dto",
		"internal/model",
		"internal/repository",
		"internal/service",
		"pkg/utils",
		"tests/pure_api_test",
	}

	// 检查目录是否存在
	for _, dir := range dirsToCheck {
		path := filepath.Join(rootDir, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("目录 %s 不存在", dir)
		} else {
			fmt.Printf("目录 %s: 正常\n", dir)
		}
	}
}
