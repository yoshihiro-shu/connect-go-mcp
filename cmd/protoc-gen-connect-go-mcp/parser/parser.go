package parser

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Service はサービス定義情報を保持します
type Service struct {
	Name     string
	FullName string // 完全修飾名 (e.g., "backend.v1.GulliverMCPService")
	Comment  string
	Methods  []Method
}

// Method はRPCメソッド情報を保持します
type Method struct {
	Name               string
	Comment            string
	RequestType        string
	RequestComment     string
	RequestConnectType string
	RequestFields      []Field
	ResponseType       string
}

// Field はフィールド情報を保持します
type Field struct {
	Name        string
	Type        string
	Comment     string
	IsRequired  bool
	Description string
}

// ParseService はサービス定義をパースします
func ParseService(s *protogen.Service) Service {
	service := Service{
		Name:     string(s.Desc.Name()),
		FullName: string(s.Desc.FullName()),
		Comment:  extractComment(s.Comments.Leading),
		Methods:  make([]Method, 0, len(s.Methods)),
	}

	for _, m := range s.Methods {
		method := ParseMethod(m)
		service.Methods = append(service.Methods, method)
	}

	return service
}

// ParseMethod はRPCメソッドをパースします
func ParseMethod(m *protogen.Method) Method {
	method := Method{
		Name:               string(m.Desc.Name()),
		Comment:            extractComment(m.Comments.Leading),
		RequestType:        m.Input.GoIdent.GoName,
		RequestComment:     extractRequestComment(m.Input),
		RequestConnectType: m.Parent.GoName + "ServiceClient",
		ResponseType:       m.Output.GoIdent.GoName,
		RequestFields:      make([]Field, 0),
	}

	// リクエストフィールドをパース
	for _, field := range m.Input.Fields {
		f := ParseField(field)
		method.RequestFields = append(method.RequestFields, f)
	}

	return method
}

// ParseField はフィールド情報をパースします
func ParseField(field *protogen.Field) Field {
	comment := extractComment(field.Comments.Leading)

	return Field{
		Name:        string(field.Desc.Name()), //　FIX: リクエストの構造体と同じ名前になっていない
		Type:        getFieldType(field.Desc),
		Comment:     comment,
		IsRequired:  isRequired(comment),
		Description: extractDescription(comment),
	}
}

// extractComment はコメントテキストを抽出します
func extractComment(comment protogen.Comments) string {
	text := strings.TrimSpace(string(comment))
	if text == "" {
		return ""
	}

	// 先頭の// を除去
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(strings.TrimPrefix(line, "//"))
	}

	return strings.Join(lines, "\n")
}

// extractRequestComment はリクエストメッセージのコメントを抽出します
func extractRequestComment(msg *protogen.Message) string {
	return extractComment(msg.Comments.Leading)
}

// getFieldType はフィールドのGo型を取得します
func getFieldType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind:
		return "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind:
		return "int64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.DoubleKind:
		return "float64"
	default:
		return "interface{}"
	}
}

// isRequired はコメントからフィールドが必須かどうかを判断します
func isRequired(comment string) bool {
	return strings.Contains(strings.ToLower(comment), "必須") ||
		strings.Contains(strings.ToLower(comment), "required")
}

// extractDescription はコメントから説明を抽出します
func extractDescription(comment string) string {
	return comment
}
