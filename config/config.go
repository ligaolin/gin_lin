package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// LoadConfig 从指定路径加载配置文件，如果文件不存在则返回默认配置
func LoadConfig(path string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultConfig := DefaultConfig()
		return &defaultConfig, nil
	}

	// 根据文件扩展名选择解析器
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return loadJSON(path)
	case ".toml":
		return loadTOML(path)
	case ".yaml", ".yml":
		return loadYAML(path)
	default:
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}
}

// loadJSON 加载 JSON 配置文件
func loadJSON(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 JSON 文件失败: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("解析 JSON 文件失败: %w", err)
	}

	return &cfg, nil
}

// loadTOML 加载 TOML 配置文件
func loadTOML(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 TOML 文件失败: %w", err)
	}

	var cfg Config
	if _, err := toml.Decode(string(file), &cfg); err != nil {
		return nil, fmt.Errorf("解析 TOML 文件失败: %w", err)
	}

	return &cfg, nil
}

// loadYAML 加载 YAML 配置文件
func loadYAML(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 YAML 文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("解析 YAML 文件失败: %w", err)
	}

	return &cfg, nil
}
