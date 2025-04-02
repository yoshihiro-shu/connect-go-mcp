package generator

import (
	"fmt"
	"path/filepath"

	"github.com/yoshihiro-shu/protoc-gen-connect-go-mpcserver/internal/parser"
	"github.com/yoshihiro-shu/protoc-gen-connect-go-mpcserver/internal/template"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generate はProtocol Bufferファイルからコードを生成します
func Generate(gen *protogen.Plugin) error {
	// 各ファイルを処理
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}

		// サービス定義がなければスキップ
		if len(f.Services) == 0 {
			continue
		}

		// 出力ファイル名を決定
		outputName := f.GeneratedFilenamePrefix + "_mpcserver.go"

		// 出力ファイルを生成
		g := gen.NewGeneratedFile(outputName, f.GoImportPath)

		// パッケージ名を取得
		pkgName := string(f.GoPackageName)
		if pkgName == "" {
			pkgName = filepath.Base(f.GoImportPath.String())
		}

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
