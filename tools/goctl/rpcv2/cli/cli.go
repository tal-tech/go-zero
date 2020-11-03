package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/generator"
	"github.com/urfave/cli"
)

// Rpc is to generate rpc service code from a proto file by specifying a proto file using flag src,
// you can specify a target folder for code generation,when the proto file has import, you can specify
// the import search directory through the proto_path command. For specific usage, please refer to protoc -h
func Rpc(c *cli.Context) error {
	src := c.String("src")
	out := c.String("dir")
	protoImportPath := c.StringSlice("proto_path")
	if len(src) == 0 {
		return errors.New("the proto source can not be nil")
	}
	if len(out) == 0 {
		return errors.New("the target directory can not be nil")
	}
	g := generator.NewDefaultRpcGenerator()
	return g.Generate(src, out, protoImportPath)
}

// RpcNew is to generate rpc greet service,this greet service can speed
// up your understanding of the zrpc service structure
func RpcNew(c *cli.Context) error {
	name := c.Args().First()
	ext := filepath.Ext(name)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}

	protoName := name + ".proto"
	filename := filepath.Join(".", name, protoName)
	src, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	err = generator.ProtoTmpl(src)
	if err != nil {
		return err
	}

	workDir := filepath.Dir(src)
	_, err = execx.Run("go mod init "+name, workDir)
	if err != nil {
		return err
	}

	g := generator.NewDefaultRpcGenerator()
	return g.Generate(src, filepath.Dir(src), nil)
}

func RpcTemplate(c *cli.Context) error {
	name := c.Args().First()
	if len(name) == 0 {
		name = "greet.proto"
	}
	return generator.ProtoTmpl(name)
}
