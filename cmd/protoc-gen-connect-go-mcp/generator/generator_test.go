package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	// Create a simple test file descriptor
	greetFileDesc := &descriptorpb.FileDescriptorProto{
		Name:    stringPtr("greet.proto"),
		Package: stringPtr("greet"),
		Options: &descriptorpb.FileOptions{
			GoPackage: stringPtr("github.com/example/greet/greetv1;greetv1"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: stringPtr("GreetService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       stringPtr("Greet"),
						InputType:  stringPtr(".greet.GreetRequest"),
						OutputType: stringPtr(".greet.GreetResponse"),
					},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: stringPtr("GreetRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   stringPtr("name"),
						Number: int32Ptr(1),
						Type:   typePtr(descriptorpb.FieldDescriptorProto_TYPE_STRING),
					},
				},
			},
			{
				Name: stringPtr("GreetResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   stringPtr("message"),
						Number: int32Ptr(1),
						Type:   typePtr(descriptorpb.FieldDescriptorProto_TYPE_STRING),
					},
				},
			},
		},
	}

	compilerVersion := &pluginpb.Version{
		Major:  int32Ptr(0),
		Minor:  int32Ptr(0),
		Patch:  int32Ptr(1),
		Suffix: stringPtr("test"),
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate:        []string{greetFileDesc.GetName()},
		Parameter:             nil,
		ProtoFile:             []*descriptorpb.FileDescriptorProto{greetFileDesc},
		SourceFileDescriptors: []*descriptorpb.FileDescriptorProto{greetFileDesc},
		CompilerVersion:       compilerVersion,
	}

	gen, err := protogen.Options{}.New(req)
	assert.Nil(t, err)

	config := Config{} // デフォルト設定
	err = Generate(gen, config)
	assert.Nil(t, err)
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func typePtr(t descriptorpb.FieldDescriptorProto_Type) *descriptorpb.FieldDescriptorProto_Type {
	return &t
}