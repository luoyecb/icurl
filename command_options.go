package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/luoyecb/icurl/lualib"

	"github.com/yuin/gopher-lua"
)

type CommandOptions struct {
	Filename string `eflag:"f,,run this file once"`

	Scheme string            `eflag:"scheme"`
	Host   string            `eflag:"host"`
	Port   int               `eflag:"port"`
	Path   string            `eflag:"path"`
	Method string            `eflag:"method"`
	Url    string            `eflag:"url"`
	Data   string            `eflag:"data"`
	Query  map[string]string `eflag:"query"`
	Header map[string]string `eflag:"header"`
}

func RunWithCommandOptions(vm *lua.LState, cmdOpts *CommandOptions) {
	if cmdOpts.Filename != "" {
		if !lualib.FileExists(cmdOpts.Filename) {
			fmt.Fprintf(os.Stderr, "file %s not exists.", cmdOpts.Filename)
			os.Exit(1)
		}
		lualib.RunLuaFile(vm, cmdOpts.Filename)
		os.Exit(0)
	} else {
		codes := make([]string, 0)
		if cmdOpts.Scheme != "" {
			codes = append(codes, fmt.Sprintf(`context.scheme = "%s"`, cmdOpts.Scheme))
		}
		if cmdOpts.Host != "" {
			codes = append(codes, fmt.Sprintf(`context.host = "%s"`, cmdOpts.Host))
		}
		if cmdOpts.Port != 0 {
			codes = append(codes, fmt.Sprintf(`context.port = %d`, cmdOpts.Port))
		}
		if cmdOpts.Path != "" {
			codes = append(codes, fmt.Sprintf(`context.path = "%s"`, cmdOpts.Path))
		}
		if cmdOpts.Method != "" {
			codes = append(codes, fmt.Sprintf(`context.method = "%s"`, cmdOpts.Method))
		}
		if cmdOpts.Url != "" {
			codes = append(codes, fmt.Sprintf(`context.url = "%s"`, cmdOpts.Url))
		}
		if cmdOpts.Data != "" {
			codes = append(codes, fmt.Sprintf(`context.data = "%s"`, cmdOpts.Data))
		}
		if len(cmdOpts.Header) > 0 {
			for key, val := range cmdOpts.Header {
				codes = append(codes, fmt.Sprintf(`set_header("%s", "%s")`, key, val))
			}
		}
		if len(cmdOpts.Query) > 0 {
			for key, val := range cmdOpts.Query {
				codes = append(codes, fmt.Sprintf(`set_query("%s", "%s")`, key, val))
			}
		}

		lualib.RunLuaCode(vm, strings.Join(codes, "\n"))
	}
}
