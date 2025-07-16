package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "bu",
		Usage:       "跨平台模板渲染工具",
		Description: `使用 Go template 语法渲染指定模板文件，支持从 JSON 文件、标准输入或命令行参数提供数据。`,
		Version:     "1.0.0",
		Authors: []*cli.Author{
			{
				Name: "Sokach7c",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "template",
				Aliases:  []string{"t"},
				Usage:    "模板文件路径",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Usage:   "JSON 数据文件路径",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "输出文件路径 (默认输出到标准输出)",
			},
			&cli.StringSliceFlag{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "设置变量 (格式: key=value)",
			},
			&cli.StringFlag{
				Name:    "json",
				Aliases: []string{"i"},
				Usage:   "直接提供 JSON 数据字符串",
			},
		},
		Action: renderTemplate,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func renderTemplate(c *cli.Context) error {
	templatePath := c.String("template")
	dataPath := c.String("data")
	outputPath := c.String("output")
	setVars := c.StringSlice("set")
	jsonDataStr := c.String("json")

	// 读取模板文件
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 解析模板
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 准备数据
	data := make(map[string]interface{})

	// 从 JSON 文件读取数据
	if dataPath != "" {
		jsonData, err := os.ReadFile(dataPath)
		if err != nil {
			return fmt.Errorf("读取数据文件失败: %w", err)
		}
		if err := json.Unmarshal(jsonData, &data); err != nil {
			return fmt.Errorf("解析 JSON 数据失败: %w", err)
		}
	}

	// 从命令行直接读取 JSON 数据
	if jsonDataStr != "" {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(jsonDataStr), &jsonMap); err != nil {
			return fmt.Errorf("解析 JSON 数据字符串失败: %w", err)
		}
		// 合并数据
		for k, v := range jsonMap {
			data[k] = v
		}
	}

	// 处理命令行设置的变量
	for _, setVar := range setVars {
		parts := strings.SplitN(setVar, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("无效的变量格式: %s (应该是 key=value)", setVar)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 尝试解析为 JSON 值（支持数字、布尔值等）
		var jsonValue interface{}
		if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {
			data[key] = jsonValue
		} else {
			// 如果不是有效的 JSON，则作为字符串处理
			data[key] = value
		}
	}

	// 准备输出
	var output io.Writer
	if outputPath != "" {
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("创建输出文件失败: %w", err)
		}
		defer file.Close()
		output = file
	} else {
		output = os.Stdout
	}

	// 渲染模板
	if err := tmpl.Execute(output, data); err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	return nil
}
