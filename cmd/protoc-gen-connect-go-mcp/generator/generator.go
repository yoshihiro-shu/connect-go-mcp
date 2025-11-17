package generator

import (
	"fmt"
	"path/filepath"

	"github.com/yoshihiro-shu/connect-go-mcp/cmd/protoc-gen-connect-go-mcp/parser"
	"github.com/yoshihiro-shu/connect-go-mcp/cmd/protoc-gen-connect-go-mcp/template"
	"google.golang.org/protobuf/compiler/protogen"
)

// Config はジェネレーターの設定を保持します
type Config struct {
	PackageSuffix    string
	PackageSuffixSet bool // package_suffixパラメータが設定されたかどうか
}

// Generate はProtocol Bufferファイルからコードを生成します
func Generate(gen *protogen.Plugin, config Config) error {
	// 各ファイルを処理
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}

		// サービス定義がなければスキップ
		if len(f.Services) == 0 {
			continue
		}

		// パッケージ名を取得
		pkgName := string(f.GoPackageName)
		if pkgName == "" {
			pkgName = filepath.Base(f.GoImportPath.String())
		}

		// 出力ファイル名とパッケージ名を決定
		var outputName string
		var importPath protogen.GoImportPath

		// package_suffixの設定に応じて出力先を決定
		if config.PackageSuffixSet && config.PackageSuffix == "" {
			// package_suffix=（空文字列）の場合：現在のディレクトリに直接配置、パッケージ名はそのまま
			outputName = f.GeneratedFilenamePrefix + ".mcpserver.go"
			importPath = f.GoImportPath
			// パッケージ名はそのまま（greetv1）
		} else {
			// package_suffixが未設定またはpackage_suffix=somevalueの場合
			suffix := "mcp" // デフォルト
			if config.PackageSuffixSet && config.PackageSuffix != "" {
				suffix = config.PackageSuffix // 指定された値を使用
			}

			// パッケージ名とディレクトリ名を決定
			packageNameWithSuffix := pkgName + suffix
			connectDirName := packageNameWithSuffix

			// ファイル名のベース部分のみを取得（connect-goと同じ仕様）
			baseName := filepath.Base(f.GeneratedFilenamePrefix)

			// 出力ファイル名とインポートパスを決定
			outputName = filepath.Join(connectDirName, baseName+".mcpserver.go")
			importPath = f.GoImportPath + "/" + protogen.GoImportPath(connectDirName)

			// パッケージ名を更新
			pkgName = packageNameWithSuffix
		}

		// 出力ファイルを生成
		g := gen.NewGeneratedFile(outputName, importPath)

		// サービス情報をパース
		services := make([]parser.Service, 0, len(f.Services))
		for _, s := range f.Services {
			service := parser.ParseService(s)
			services = append(services, service)
		}

		// コード生成
		if err := template.GenerateCode(g, pkgName, services); err != nil {
			return fmt.Errorf("failed to generate code for %s: %w", f.Desc.Path(), err)
		}
	}

	return nil
}
