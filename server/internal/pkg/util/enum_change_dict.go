package util

import (
	"admin/internal/types"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EnumChangeDict
// 扫描目录constant/enum 下的文件，文件名称为 map 的key, 文件内容 以结构体 types.Options 为值， 其中 value 为常量的值，label 为常量的注解 生成 map[string][]types.Options
func EnumChangeDict(enumDir string) map[string][]*types.Options {
	result := make(map[string][]*types.Options)

	files, err := os.ReadDir(enumDir)
	if err != nil {
		log.Printf("Error reading directory %s: %v\n", enumDir, err)
		return result
	}

	fset := token.NewFileSet()
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(enumDir, file.Name())
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing file %s: %v\n", filePath, err)
			continue
		}

		fileNameKey := strings.TrimSuffix(file.Name(), ".go")
		var options []*types.Options

		ast.Inspect(node, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			// 检查是否是常量声明
			if !ok || decl.Tok != token.CONST {
				return true // 继续遍历子节点
			}

			// 实际的 iota 值需要根据 spec 在 decl.Specs 中的索引来确定
			// currentIota := 0 // 移除简单的 iota 计数器

			// 用于跟踪当前的 iota 值和基础偏移量
			currentIotaValue := 0
			iotaOffset := 0
			hasIotaExpression := false

			for i, spec := range decl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				var val interface{} // 改为 interface{} 以支持不同类型
				isIntValue := false
				isStringValue := false

				// 检查是否有显式赋值
				if len(valueSpec.Values) > 0 {
					// 只处理第一个值
					valueExpr := valueSpec.Values[0]
					if basicLit, ok := valueExpr.(*ast.BasicLit); ok {
						if basicLit.Kind == token.INT {
							// 尝试解析整数值
							parsedVal, err := strconv.Atoi(basicLit.Value)
							if err == nil {
								val = parsedVal
								isIntValue = true
								// 重置 iota 跟踪
								currentIotaValue = parsedVal
								iotaOffset = 0
								hasIotaExpression = false
							} else {
								log.Printf("Warn: Could not parse integer literal '%s' in %s: %v\n", basicLit.Value, filePath, err)
							}
						} else if basicLit.Kind == token.STRING {
							// 解析字符串值 (去除引号)
							parsedVal, err := strconv.Unquote(basicLit.Value)
							if err == nil {
								val = parsedVal
								isStringValue = true
							} else {
								log.Printf("Warn: Could not parse string literal '%s' in %s: %v\n", basicLit.Value, filePath, err)
							}
						}
						// 其他类型 (float, complex, char) 暂不处理
					} else if binExpr, ok := valueExpr.(*ast.BinaryExpr); ok {
						// 处理 iota + 1 这样的表达式
						if ident, ok := binExpr.X.(*ast.Ident); ok && ident.Name == "iota" {
							if basicLit, ok := binExpr.Y.(*ast.BasicLit); ok && basicLit.Kind == token.INT {
								if offset, err := strconv.Atoi(basicLit.Value); err == nil {
									iotaOffset = offset
									val = i + iotaOffset
									isIntValue = true
									hasIotaExpression = true
									currentIotaValue = val.(int)
								}
							}
						} else {
							// 其他二元表达式，使用 iota 近似
							val = i
							isIntValue = true
							currentIotaValue = i
						}
					} else {
						// 如果不是基本字面量 (可能是标识符、调用等)，尝试使用 iota
						val = i
						isIntValue = true
						currentIotaValue = i
					}
				} else {
					// 没有显式赋值，使用前一个值的递增
					if hasIotaExpression {
						// 如果前面有 iota 表达式，继续递增
						currentIotaValue++
						val = currentIotaValue
					} else {
						// 否则使用索引
						val = i
						currentIotaValue = i
					}
					isIntValue = true
				}

				// 只有整数或字符串类型的常量才添加到 options
				if isIntValue || isStringValue {
					// 尝试获取注释作为 label
					var label string
					if valueSpec.Comment != nil && len(valueSpec.Comment.List) > 0 {
						label = strings.TrimSpace(strings.TrimPrefix(valueSpec.Comment.List[0].Text, "//"))
					} else if valueSpec.Doc != nil && len(valueSpec.Doc.List) > 0 {
						label = strings.TrimSpace(strings.TrimPrefix(valueSpec.Doc.List[0].Text, "//"))
					}

					// 添加到 options 列表 (只添加有名字的常量)
					if len(valueSpec.Names) > 0 && valueSpec.Names[0].Name != "_" {
						options = append(options, &types.Options{
							Value: val, // 使用解析出的或 iota 的值
							Label: label,
						})
					}
				}
				// iota 的递增由 Go 编译器处理，我们只需在需要时使用其值 (这里用索引近似)
				// currentIota++ // 移除简单的 iota 计数器
			}
			// 退出 const 块的处理
			return false // 不需要深入遍历 const 块内部的具体表达式等
		})

		// 将 Value 为 整数 0 的选项移到末尾
		if len(options) > 1 { // 只有一个元素时无需移动
			zeroIndex := -1
			for i, opt := range options {
				// 检查 Value 是否为整数 0
				if intVal, ok := opt.Value.(int); ok && intVal == 0 {
					zeroIndex = i
					break
				}
			}
			if zeroIndex != -1 && zeroIndex != len(options)-1 {
				zeroOption := options[zeroIndex]
				// 从原位置删除
				options = append(options[:zeroIndex], options[zeroIndex+1:]...)
				// 添加到末尾
				options = append(options, zeroOption)
			}
		}

		if len(options) > 0 {
			result[fileNameKey] = options
		}
	}

	return result
}
