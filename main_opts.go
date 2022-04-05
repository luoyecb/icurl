package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/luoyecb/icurl/lualib"

	"github.com/yuin/gopher-lua"
)

type CommandOptions struct {
	Filename string            `flag:"f,,run this file once"`
	Method   string            `flag:"m,,http method"`
	Url      string            `flag:"url,,request url"`
	Data     string            `flag:"d,,request data"`
	Query    map[string]string `flag:"q,,request data"`
	Header   map[string]string `flag:"h,,http headers"`
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
