package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	callTemplateText = `{{.head}}


package {{.filePackage}}

import (
	"context"

	{{.package}}

	"github.com/tal-tech/go-zero/zrpc"
)

type (
	{{.alias}}

	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli zrpc.Client
	}
)

func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return &default{{.serviceName}}{
		cli: cli,
	}
}

{{.functions}}
`

	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context{{if .hasReq}},in *{{.pbRequest}}{{end}}) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error)`

	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.serviceName}}) {{.method}}(ctx context.Context{{if .hasReq}},in *{{.pbRequest}}{{end}}) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error) {
	client := {{.package}}.New{{.rpcServiceName}}Client(m.cli.Conn())
	return client.{{.method}}(ctx,{{if .hasReq}} in{{end}})
}
`
)

// GenCall generates the rpc client code, which is the entry point for the rpc service call.
// It is a layer of encapsulation for the rpc client and shields the details in the pb.
func (g *DefaultGenerator) GenCall(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetCall()
	service := proto.Service
	head := util.GetHead(proto.Name)

	callFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name)
	if err != nil {
		return err
	}

	filename := filepath.Join(dir.Filename, fmt.Sprintf("%s.go", callFilename))
	functions, err := g.genFunction(proto.PbPackage, service)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, service)
	if err != nil {
		return err
	}

	text, err := util.LoadTemplate(category, callTemplateFile, callTemplateText)
	if err != nil {
		return err
	}

	alias := collection.NewSet()
	for _, item := range proto.Message {
		alias.AddStr(fmt.Sprintf("%s = %s", parser.CamelCase(item.Name), fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(item.Name))))
	}

	err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"name":        callFilename,
		"alias":       strings.Join(alias.KeysStr(), util.NL),
		"head":        head,
		"filePackage": dir.Base,
		"package":     fmt.Sprintf(`"%s"`, ctx.GetPb().Package),
		"serviceName": stringx.From(service.Name).ToCamel(),
		"functions":   strings.Join(functions, util.NL),
		"interface":   strings.Join(iFunctions, util.NL),
	}, filename, true)
	return err
}

func (g *DefaultGenerator) genFunction(goPackage string, service parser.Service) ([]string, error) {
	functions := make([]string, 0)
	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, callFunctionTemplateFile, callFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		var streamServer string
		if rpc.StreamsRequest && rpc.StreamsReturns {
			streamServer = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(service.Name)+"_StreamClient")
		} else if rpc.StreamsRequest {
			streamServer = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(service.Name)+"_ClientStreamClient")
		} else {
			streamServer = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(service.Name)+"_ServerStreamClient")
		}
		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]interface{}{
			"serviceName":    stringx.From(service.Name).ToCamel(),
			"rpcServiceName": parser.CamelCase(service.Name),
			"method":         parser.CamelCase(rpc.Name),
			"package":        goPackage,
			"pbRequest":      parser.CamelCase(rpc.RequestType),
			"pbResponse":     parser.CamelCase(rpc.ReturnsType),
			"hasComment":     len(comment) > 0,
			"comment":        comment,
			"hasReq":         !rpc.StreamsRequest,
			"notStream":      !rpc.StreamsRequest && !rpc.StreamsReturns,
			"streamBody":     streamServer,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}
	return functions, nil
}

func (g *DefaultGenerator) getInterfaceFuncs(goPackage string, service parser.Service) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, callInterfaceFunctionTemplateFile, callInterfaceFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		var streamServer string
		if rpc.StreamsRequest && rpc.StreamsReturns {
			streamServer = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(service.Name)+"_StreamClient")
		} else if rpc.StreamsRequest {
			streamServer = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(service.Name)+"_ClientStreamClient")
		} else {
			streamServer = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(service.Name)+"_ServerStreamClient")
		}

		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]interface{}{
				"hasComment": len(comment) > 0,
				"comment":    comment,
				"method":     parser.CamelCase(rpc.Name),
				"hasReq":     !rpc.StreamsRequest,
				"pbRequest":  parser.CamelCase(rpc.RequestType),
				"notStream":  !rpc.StreamsRequest && !rpc.StreamsReturns,
				"pbResponse": parser.CamelCase(rpc.ReturnsType),
				"streamBody": streamServer,
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}
