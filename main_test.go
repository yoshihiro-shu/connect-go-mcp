package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	greetv1 "github.com/yoshihiro-shu/protoc-gen-connect-go-mpcserver/testdata/greet/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	stdout, stderr, exitCode := testRunCommand(t, nil, "--version")
	t.Logf("stdout: %s", stdout.String())
	t.Logf("stderr: %s", stderr.String())
	t.Logf("exitCode: %d", exitCode)
	assert.Equal(t, stdout.String(), version+"\n")
	assert.Equal(t, stderr.String(), "")
	assert.Equal(t, exitCode, 0)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	greetFileDesc := protodesc.ToFileDescriptorProto(greetv1.File_greet_proto)
	compilerVersion := &pluginpb.Version{
		Major:  ptr(int32(0)),
		Minor:  ptr(int32(0)),
		Patch:  ptr(int32(1)),
		Suffix: ptr("test"),
	}
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate:        []string{greetFileDesc.GetName()},
		Parameter:             nil,
		ProtoFile:             []*descriptorpb.FileDescriptorProto{greetFileDesc},
		SourceFileDescriptors: []*descriptorpb.FileDescriptorProto{greetFileDesc},
		CompilerVersion:       compilerVersion,
	}
	rsp := testGenerate(t, req)
	assert.Nil(t, rsp.Error)

	t.Logf("rsp: %+v", rsp)

	assert.Equal(t, len(rsp.File), 1)
	file := rsp.File[0]
	t.Logf("file: %+v", file.GetName())
	t.Logf("file: %+v", file.GetContent())
	// assert.Equal(t, file.GetName(), "./testdata/greet/gen/greetv1connect/greet.connect.go")
	// assert.NotZero(t, file.GetContent())
}

func testRunCommand(t *testing.T, stdin io.Reader, args ...string) (stdout, stderr *bytes.Buffer, exitCode int) {
	t.Helper()

	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	args = append([]string{"run", "main.go"}, args...)

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	assert.Nil(t, cmd.Run(), fmt.Sprintf("Run go %v", args))
	exitCode = cmd.ProcessState.ExitCode()
	return stdout, stderr, exitCode
}

func testGenerate(t *testing.T, req *pluginpb.CodeGeneratorRequest) *pluginpb.CodeGeneratorResponse {
	t.Helper()

	inputBytes, err := proto.Marshal(req)
	assert.Nil(t, err)

	stdout, stderr, exitCode := testRunCommand(t, bytes.NewReader(inputBytes))
	assert.Equal(t, exitCode, 0)
	assert.Equal(t, stderr.String(), "")
	assert.True(t, len(stdout.Bytes()) > 0)

	var output pluginpb.CodeGeneratorResponse
	assert.Nil(t, proto.Unmarshal(stdout.Bytes(), &output))
	return &output
}

func ptr[T any](v T) *T {
	return &v
}
