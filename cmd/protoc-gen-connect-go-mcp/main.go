package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/yoshihiro-shu/connect-go-mcp/cmd/protoc-gen-connect-go-mcp/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("%v\n", version)
		return
	}

	var config generator.Config

	// グローバルフラグのパース
	protogen.Options{
		ParamFunc: func(name, value string) error {
			// パラメータをカンマ区切りで分割して処理
			if name == "" && value != "" {
				// 単一の値の場合（例：paths=source_relative,package_suffix=mcp）
				pairs := strings.Split(value, ",")
				for _, pair := range pairs {
					if keyValue := strings.SplitN(pair, "=", 2); len(keyValue) == 2 {
						switch keyValue[0] {
						case "package_suffix":
							config.PackageSuffix = keyValue[1]
							config.PackageSuffixSet = true
						default:
							// 未知のパラメータは無視
						}
					}
				}
			} else {
				// 名前付きパラメータの場合
				switch name {
				case "package_suffix":
					config.PackageSuffix = value
					config.PackageSuffixSet = true
				default:
					// 未知のパラメータは無視
				}
			}
			return nil
		},
	}.Run(func(gen *protogen.Plugin) error {
		// ジェネレーターを呼び出す
		err := generator.Generate(gen, config)
		if err != nil {
			log.Printf("Error during generation: %v", err)
		}
		return err
	})
}
