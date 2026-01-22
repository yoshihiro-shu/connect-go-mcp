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
