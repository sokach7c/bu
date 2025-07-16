# 模板渲染工具使用说明

## 基本用法

```bash
# 使用模板文件和 JSON 数据文件
./bu -t template.tmpl -d data.json

# 使用模板文件和命令行 JSON 字符串
./bu -t template.tmpl -i '{"name":"张三","age":25}'

# 使用模板文件和命令行变量设置
./bu -t template.tmpl -s name=张三 -s age=25

# 输出到文件
./bu -t template.tmpl -d data.json -o output.txt

# 组合使用多种数据源
./bu -t template.tmpl -d data.json -i '{"city":"北京"}' -s department=技术部
```

## 参数说明

- `-t, --template`: 模板文件路径（必需）
- `-d, --data`: JSON 数据文件路径
- `-i, --json`: 直接提供 JSON 数据字符串
- `-o, --output`: 输出文件路径（默认输出到标准输出）
- `-s, --set`: 设置变量（格式: key=value）

## 数据优先级

当多个数据源同时提供时，优先级为：
1. 命令行变量设置（-s）
2. JSON 字符串（-i）
3. JSON 文件（-d）
