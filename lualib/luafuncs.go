package lualib

import (
	"fmt"
	"io/ioutil"

	"github.com/yuin/gopher-lua"
)

var (
	FuncsMap = map[string]lua.LGFunction{
		"reset":       reset,
		"loadf":       loadf,
		"loadh":       loadh,
		"listh":       listh,
		"saveh":       saveh,
		"debugc":      debugc,
		"send":        send,
		"send_get":    send_get,
		"send_post":   send_post,
		"send_form":   send_form,
		"set_query":   set_query,
		"set_header":  set_header,
		"json_encode": json_encode,
		"shell":       shell,
		"help":        help,
	}
)

func RegisterFuncs(vm *lua.LState) {
	for fnName, fn := range FuncsMap {
		vm.SetGlobal(fnName, vm.NewFunction(fn))
	}
}

func reset(vm *lua.LState) int {
	InitContext(vm)
	return 0
}

func loadf(vm *lua.LState) int {
	if !CheckArg(vm, 1, "too few args, need absolute filepath") {
		return 1
	}

	fpath := GetRealPath(vm.ToString(1))
	if !FileExists(fpath) {
		return 0
	}
	err := RunLuaFile(vm, fpath)
	if err != nil {
		vm.RaiseError("call lua file error: %v", err)
		return 1
	}
	return 0
}

func loadh(vm *lua.LState) int {
	if !CheckArg(vm, 1, "too few args, need relative filepath with base path") {
		return 1
	}

	fpath := GetRealPath(GetBasePath() + "/" + vm.ToString(1))
	if !FileExists(fpath) {
		return 0
	}
	err := RunLuaFile(vm, fpath)
	if err != nil {
		vm.RaiseError("call lua file error: %v", err)
		return 1
	}
	return 0
}

func listh(vm *lua.LState) int {
	fpath := GetRealPath(GetBasePath())
	dirs := ListDir(fpath)
	for _, dir := range dirs {
		fmt.Println(dir)
	}
	return 0
}

func saveh(vm *lua.LState) int {
	if !CheckArg(vm, 1, "too few args, need relative filepath with base path") {
		return 1
	}

	argpath := vm.ToString(1)
	fpath := GetRealPath(GetBasePath() + "/" + argpath)

	// overwrite existing file?
	overwrite := false
	if vm.GetTop() > 1 {
		overwrite = vm.CheckBool(2)
	}

	if !overwrite {
		if FileExists(fpath) {
			vm.RaiseError("%s exists", argpath)
			return 1
		}
	}

	ctx, ok := CheckGetContext(vm)
	if !ok {
		return 1
	}

	strcode, err := LTableToLuaCode(ctx)
	if err != nil {
		vm.RaiseError("saveh error: %v", err)
		return 1
	}

	if err := ioutil.WriteFile(fpath, []byte("context = "+strcode), 0644); err != nil {
		vm.RaiseError("saveh write file error: %v", err)
		return 1
	}
	return 0
}

func debugc(vm *lua.LState) int {
	ctx, ok := CheckGetContext(vm)
	if !ok {
		return 1
	}

	str, err := LTableToJsonString(ctx, true)
	if err != nil {
		vm.RaiseError("debugc error: %v", err)
		return 1
	}
	fmt.Println(str)
	return 0
}

func send0(vm *lua.LState, method string, header map[string]string) (nres int) {
	defer func() {
		if err := recover(); err != nil {
			vm.RaiseError("call send error: %v", err)
			nres = 1
		}
	}()

	var formatJson = false
	if vm.GetTop() > 0 {
		formatJson = vm.CheckBool(1)
	}

	ctx, ok := CheckGetContext(vm)
	if !ok {
		return 1
	}

	httpCtx := NewHttpContext()
	httpCtx.Scheme = GetLTableString(ctx, "scheme", "http")
	httpCtx.Host = GetLTableString(ctx, "host", "localhost")
	httpCtx.Port = GetLTableInt(ctx, "port", 80)
	httpCtx.Path = GetLTableString(ctx, "path", "")
	httpCtx.Url = GetLTableString(ctx, "url", "")
	httpCtx.Data = GetLTableString(ctx, "data", "")
	httpCtx.Query = LTableToMapString(GetLTableTable(ctx, "query"))

	httpCtx.Header = LTableToMapString(GetLTableTable(ctx, "header"))
	if len(header) > 0 {
		for k, v := range header {
			httpCtx.Header[k] = v
		}
	}

	if method == "" {
		httpCtx.Method = GetLTableString(ctx, "method", "GET")
	} else {
		httpCtx.Method = method
	}

	bodyStr, err := httpCtx.Send()
	if err != nil {
		panic(err)
	}

	if formatJson {
		fmt.Println(JsonPrettyFormat(bodyStr))
	} else {
		fmt.Println(bodyStr)
	}
	return 0
}

func send(vm *lua.LState) int {
	return send0(vm, "", nil)
}

func send_get(vm *lua.LState) int {
	return send0(vm, "GET", nil)
}

func send_post(vm *lua.LState) int {
	return send0(vm, "POST", nil)
}

func send_form(vm *lua.LState) int {
	return send0(vm, "POST", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
}

func set_query(vm *lua.LState) int {
	if !CheckArg(vm, 2, "too few args, need query pair (key,value)") {
		return 1
	}

	qk := vm.CheckString(1)
	vm.CheckTypes(2, lua.LTString, lua.LTNumber)
	qv := vm.Get(2)

	ctx, ok := CheckGetContext(vm)
	if !ok {
		return 1
	}

	query := GetLTableTable(ctx, "query", nil)
	if query == nil {
		query = vm.NewTable()
		SetLTable(ctx, "query", query)
	}
	SetLTable(query, qk, qv)
	return 0
}

func set_header(vm *lua.LState) int {
	if !CheckArg(vm, 2, "too few args, need header pair (key,value)") {
		return 1
	}

	hk := vm.CheckString(1)
	hv := vm.CheckString(2)

	ctx, ok := CheckGetContext(vm)
	if !ok {
		return 1
	}

	header := GetLTableTable(ctx, "header", nil)
	if header == nil {
		header = vm.NewTable()
		SetLTable(ctx, "header", header)
	}
	SetLTableString(header, hk, hv)
	return 0
}

func json_encode(vm *lua.LState) int {
	if !CheckArg(vm, 1, "too few args") {
		return 1
	}

	var formatJson bool
	if vm.GetTop() > 1 {
		formatJson = vm.CheckBool(2)
	}

	tab := vm.CheckTable(1)

	s, err := LTableToJsonString(tab, formatJson)
	if err != nil {
		vm.RaiseError("json_encode error: %v", err)
		return 1
	}
	vm.Push(lua.LString(s))
	return 1
}

func shell(vm *lua.LState) int {
	if !CheckArg(vm, 1, "too few args") {
		return 1
	}

	cmdStr := vm.CheckString(1)
	out, err := ShellExec(cmdStr)
	if err != nil {
		vm.RaiseError("shell error: %v", err)
		return 1
	}
	fmt.Println(out)
	return 0
}

func help(vm *lua.LState) int {
	fmt.Println(`=== context
context = {
	scheme = "http", # http|https
	host   = "localhost", # must string
	port   = 80,     # must number
	path   = "",     # must string
	method = "GET",  # GET|PUT|POST|DELETE
	url    = "",     # must string, if set, other fields are ignored
	data   = "",     # must string, if method=GET, data are ignored, otherwise use data instead of query
	query  = {},     # must table
	header = {},     # must table
}

=== functions
exit|quit                 : exit
reset()                   : reset context
loadf(string)             : load lua file, absolute path
loadh(string)             : load lua file, default in dir ~/.icurl/
listh(string)             : list lua file, default in dir ~/.icurl/
saveh(string, [bool])     : save lua file, default in dir ~/.icurl/, bool arg means whether overwrite existing file or not
debugc()                  : print context information
send([bool])              : send http requeset, http method is context.method, bool arg means json pretty formatting
send_get([bool])          : send http get requeset, bool arg means json pretty formatting
send_post([bool])         : send http post requeset, bool arg means json pretty formatting
send_form([bool])         : send http post requeset, with header "Content-Type:application/x-www-form-urlencoded", bool arg means json pretty formatting
set_query(string, string) : set context.query
set_header(string, string): set context.header
json_encode(table, [bool]): json encode, bool arg means json pretty formatting
shell(string)             : exec shell command
!string                   : exec shell command
help()                    : show this help information

Everything follows Lua grammar.
Good luck.
`)
	return 0
}
