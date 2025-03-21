package config

import (
	"os"
	"path/filepath"
)

// TestConfig 测试配置
type TestConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Reports  ReportConfig   `yaml:"reports"`
}

// ReportConfig 测试报告配置
type ReportConfig struct {
	Dir      string `yaml:"dir"`
	MockDir  string `yaml:"mock_dir"`
	Format   string `yaml:"format"`
	Coverage bool   `yaml:"coverage"`
}

// GetTestConfig 获取测试配置
func GetTestConfig() (*TestConfig, error) {
	// 测试环境配置
	testConfig := &TestConfig{
		Database: DatabaseConfig{
			Driver:          "sqlite",
			Database:        ":memory:",
			ShowSQL:         true,
			LogLevel:        "debug",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 3600,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: 8089,
		},
		Reports: ReportConfig{
			Dir:      filepath.Join(GetProjectRoot(), "..", "temp", "test-reports"),
			MockDir:  filepath.Join(GetProjectRoot(), "tests", "fixtures"),
			Format:   "html",
			Coverage: true,
		},
	}

	return testConfig, nil
}

// GetProjectRoot 获取项目根目录
func GetProjectRoot() string {
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			// 到达根目录仍未找到 go.mod，默认返回当前目录
			return filepath.Dir(os.Args[0])
		}
		wd = parent
	}
}
