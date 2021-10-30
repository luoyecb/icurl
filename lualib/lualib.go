package lualib

import (
	"os"

	"github.com/yuin/gopher-lua"
)

const (
	INIT_ENV_NAME  = "ICURL_PATH"
	INIT_HOME_PATH = "~/.icurl"
	INIT_LUA_FILE  = "init.lua"
	INIT_CODE      = `
context = {
	scheme = "http",
	host   = "localhost",
	port   = 80,
	path   = "",
	method = "GET",
	url    = "",
	data   = "",
	query  = {},
	header = {},
}
`
	LUA_CODE = `
-- 禁止设置全局变量，除了 context
setmetatable(_G, {
	__newindex = function(t, n, v)
		if n ~= "context" or type(v) ~= "function" then
			error("can not set global variable, except context")
		end
		if type(v) ~= "table" then
			error("context must be table")
		end
	end
})
`
)

func Init(vm *lua.LState) error {
	if err := InitContext(vm); err != nil {
		return err
	}
	RegisterFuncs(vm)
	return RunLuaCode(vm, LUA_CODE)
}

func InitContext(vm *lua.LState) error {
	err := RunLuaCode(vm, INIT_CODE)
	if err != nil {
		return err
	}
	return RunLuaFile(vm, GetRealPath(GetBasePath()+"/"+INIT_LUA_FILE))
}

func GetBasePath() string {
	envpath, ok := os.LookupEnv(INIT_ENV_NAME)
	if !ok && envpath == "" {
		envpath = INIT_HOME_PATH
	}
	return envpath
}

func GetContext(vm *lua.LState) (*lua.LTable, bool) {
	ctx, ok := vm.GetGlobal("context").(*lua.LTable)
	return ctx, ok
}

func CheckGetContext(vm *lua.LState) (*lua.LTable, bool) {
	ctx, ok := GetContext(vm)
	if !ok {
		vm.RaiseError("context must be table")
		return nil, false
	}
	return ctx, true
}

func GetLTableString(table *lua.LTable, field string, defval ...string) string {
	v, ok := table.RawGetString(field).(lua.LString)
	if ok {
		return string(v)
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return ""
}

func GetLTableInt(table *lua.LTable, field string, defval ...int) int {
	v, ok := table.RawGetString(field).(lua.LNumber)
	if ok {
		return int(v)
	}
	if len(defval) > 0 {
		return defval[0]
	}
	return 0
}

func GetLTableTable(table *lua.LTable, field string, deftab ...*lua.LTable) *lua.LTable {
	tab, ok := table.RawGetString(field).(*lua.LTable)
	if ok {
		return tab
	}
	if len(deftab) > 0 {
		return deftab[0]
	}
	return nil
}

func SetLTableString(table *lua.LTable, field string, val string) {
	table.RawSetString(field, lua.LString(val))
}

func SetLTable(table *lua.LTable, field string, val lua.LValue) {
	table.RawSetString(field, val)
}

func RunLuaFile(vm *lua.LState, fpath string) error {
	if !FileExists(fpath) {
		return nil
	}
	return vm.DoFile(fpath)
}

func RunLuaCode(vm *lua.LState, code string) error {
	if code == "" {
		return nil
	}
	return vm.DoString(code)
}

func CallLuaFunc(vm *lua.LState, fn string, nret int, args ...lua.LValue) ([]lua.LValue, error) {
	err := vm.CallByParam(lua.P{
		Fn:      vm.GetGlobal(fn),
		NRet:    nret,
		Protect: true,
	}, args...)
	if err != nil {
		return nil, err
	}

	ret := make([]lua.LValue, 0, nret)
	for i := nret; i >= 1; i-- {
		ret = append(ret, vm.Get(-i))
	}

	vm.Pop(nret)
	return ret, nil
}
