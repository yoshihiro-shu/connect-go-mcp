package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestParseService_FullName(t *testing.T) {
	// Create a mock file descriptor with a service
	fileDescProto := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Package: proto.String("backend.v1"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/gen/backendv1"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("TestService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("GetItem"),
						InputType:  proto.String(".backend.v1.GetItemRequest"),
						OutputType: proto.String(".backend.v1.GetItemResponse"),
					},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("GetItemRequest"),
			},
			{
				Name: proto.String("GetItemResponse"),
			},
		},
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"test.proto"},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{fileDescProto},
	}

	gen, err := protogen.Options{}.New(req)
	assert.NoError(t, err)
	assert.Len(t, gen.Files, 1)

	file := gen.Files[0]
	assert.Len(t, file.Services, 1)

	service := ParseService(file.Services[0])

	// Test that FullName includes package name
	assert.Equal(t, "TestService", service.Name)
	assert.Equal(t, "backend.v1.TestService", service.FullName)

	// Test that methods are parsed correctly
	assert.Len(t, service.Methods, 1)
	assert.Equal(t, "GetItem", service.Methods[0].Name)
}

func TestParseService_FullName_NestedPackage(t *testing.T) {
	// Test with deeply nested package
	fileDescProto := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("nested.proto"),
		Package: proto.String("com.example.api.v1"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/gen/comexampleapiv1"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("NestedService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("DoSomething"),
						InputType:  proto.String(".com.example.api.v1.DoSomethingRequest"),
						OutputType: proto.String(".com.example.api.v1.DoSomethingResponse"),
					},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("DoSomethingRequest"),
			},
			{
				Name: proto.String("DoSomethingResponse"),
			},
		},
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"nested.proto"},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{fileDescProto},
	}

	gen, err := protogen.Options{}.New(req)
	assert.NoError(t, err)

	file := gen.Files[0]
	service := ParseService(file.Services[0])

	assert.Equal(t, "NestedService", service.Name)
	assert.Equal(t, "com.example.api.v1.NestedService", service.FullName)
}

func TestParseService_FullName_NoPackage(t *testing.T) {
	// Test with no package (edge case)
	fileDescProto := &descriptorpb.FileDescriptorProto{
		Name: proto.String("nopackage.proto"),
		// No Package field
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/gen/nopackage"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("SimpleService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("Call"),
						InputType:  proto.String(".CallRequest"),
						OutputType: proto.String(".CallResponse"),
					},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("CallRequest"),
			},
			{
				Name: proto.String("CallResponse"),
			},
		},
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"nopackage.proto"},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{fileDescProto},
	}

	gen, err := protogen.Options{}.New(req)
	assert.NoError(t, err)

	file := gen.Files[0]
	service := ParseService(file.Services[0])

	assert.Equal(t, "SimpleService", service.Name)
	// When no package, FullName should just be the service name
	assert.Equal(t, "SimpleService", service.FullName)
}

func TestParseField_TrailingCommentFallback(t *testing.T) {
	// Field number of DescriptorProto.field within a message, and of
	// FileDescriptorProto.message_type within a file. Used to build the
	// SourceCodeInfo paths that attach comments to specific elements.
	const (
		fileMessageTypeField = 4 // FileDescriptorProto.message_type
		messageFieldField    = 2 // DescriptorProto.field
	)

	fileDescProto := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("field.proto"),
		Package: proto.String("backend.v1"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/gen/backendv1"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("FieldMessage"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("name"),
						Number: proto.Int32(1),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:   proto.String("id"),
						Number: proto.Int32(2),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:   proto.String("title"),
						Number: proto.Int32(3),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
			},
		},
		SourceCodeInfo: &descriptorpb.SourceCodeInfo{
			Location: []*descriptorpb.SourceCodeInfo_Location{
				{
					// name: trailing-only comment, should be used as fallback
					Path:             []int32{fileMessageTypeField, 0, messageFieldField, 0},
					Span:             []int32{0, 0, 10},
					TrailingComments: proto.String(" Name\n"),
				},
				{
					// id: trailing comment that marks it required
					Path:             []int32{fileMessageTypeField, 0, messageFieldField, 1},
					Span:             []int32{1, 0, 10},
					TrailingComments: proto.String(" ID (required)\n"),
				},
				{
					// title: both present, leading should win
					Path:             []int32{fileMessageTypeField, 0, messageFieldField, 2},
					Span:             []int32{2, 0, 10},
					LeadingComments:  proto.String(" Leading title\n"),
					TrailingComments: proto.String(" Trailing title\n"),
				},
			},
		},
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"field.proto"},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{fileDescProto},
	}

	gen, err := protogen.Options{}.New(req)
	assert.NoError(t, err)
	assert.Len(t, gen.Files, 1)

	file := gen.Files[0]
	assert.Len(t, file.Messages, 1)
	fields := file.Messages[0].Fields
	assert.Len(t, fields, 3)

	// name: leading empty -> falls back to trailing comment
	name := ParseField(fields[0])
	assert.Equal(t, "name", name.Name)
	assert.Equal(t, "Name", name.Description)
	assert.False(t, name.IsRequired)

	// id: required keyword in trailing comment marks the field required
	id := ParseField(fields[1])
	assert.Equal(t, "ID (required)", id.Description)
	assert.True(t, id.IsRequired)

	// title: leading comment takes precedence over trailing
	title := ParseField(fields[2])
	assert.Equal(t, "Leading title", title.Description)
}
