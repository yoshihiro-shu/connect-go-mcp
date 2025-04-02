package main

import (
	"flag"
	"fmt"

	"github.com/yoshihiro-shu/protoc-gen-connect-go-mpcserver/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-connect-go-mpcserver %v\n", version)
		return
	}

	// グローバルフラグのパース
	protogen.Options{
		ParamFunc: func(name, value string) error {
			// 将来のオプション処理
			return nil
		},
	}.Run(func(gen *protogen.Plugin) error {
		// ジェネレーターを呼び出す
		return generator.Generate(gen)
	})
}
