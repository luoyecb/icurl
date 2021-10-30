package lualib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/yuin/gopher-lua"
)

func CheckArg(vm *lua.LState, narg int, errmsg string) bool {
	if vm.GetTop() < narg {
		vm.RaiseError(errmsg)
		return false
	}
	return true
}

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(fpath string) bool {
	dir, err := os.Stat(fpath)
	if err != nil {
		return false
	}
	return dir.IsDir()
}

func ListDir(fpath string) []string {
	res := make([]string, 0)
	if !IsDir(fpath) {
		return res
	}

	infos, err := ioutil.ReadDir(fpath)
	if err != nil {
		return res
	}
	for _, finfo := range infos {
		res = append(res, finfo.Name())
	}
	return res
}

func GetRealPath(path string) string {
	home := GetCurrentUserHomeDir()
	if home == "" {
		return path
	}
	if strings.HasPrefix(path, "~/") {
		return strings.Replace(path, "~", home, 1)
	}
	return path
}

func GetCurrentUserHomeDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}

func CompareStringIgnoreCase(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

func StringIsLetter(s string) bool {
	for i, j := 0, len(s); i < j; i++ {
		if !IsLetter(s[i]) {
			return false
		}
	}
	return true
}

func IsLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func LTableToJsonString(table *lua.LTable, formatJson bool) (string, error) {
	m, err := LTableToMap(table)
	if err != nil {
		return "", err
	}

	var bytes []byte

	if formatJson {
		bytes, err = json.MarshalIndent(m, "", "    ")
	} else {
		bytes, err = json.Marshal(m)
	}
	if err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

func LTableToMap(table *lua.LTable) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})

	tableMap := make(map[lua.LValue]lua.LValue)
	table.ForEach(func(tk, tv lua.LValue) {
		tableMap[tk] = tv
	})

	for k, v := range tableMap {
		if tk, ok := k.(lua.LString); ok {
			sk := string(tk)
			var sv interface{}

			switch v := v.(type) {
			case lua.LBool:
				sv = bool(v)
			case lua.LNumber:
				sv = float64(v)
			case lua.LString:
				sv = string(v)
			case *lua.LNilType:
				sv = nil
			case *lua.LTable:
				ssv, err := LTableToMap(v)
				if err != nil {
					return nil, err
				}
				sv = ssv
			default:
				return nil, errors.New("table value only supported type of bool|number|string|nil|table")
			}

			jsonMap[sk] = sv
		} else {
			return nil, errors.New("table key only supported string type")
		}
	}

	return jsonMap, nil
}

func JsonPrettyFormat(s string) string {
	var holder interface{}
	if err := json.Unmarshal([]byte(s), &holder); err != nil {
		return s
	} else {
		if bytes, err := json.MarshalIndent(&holder, "", "    "); err != nil {
			return s
		} else {
			return string(bytes)
		}
	}
}

func LTableToMapString(table *lua.LTable) map[string]string {
	res := make(map[string]string)

	if table != nil {
		table.ForEach(func(k, v lua.LValue) {
			res[k.String()] = v.String()
		})
	}
	return res
}

func LTableToLuaCode(table *lua.LTable) (string, error) {
	maps, err := LTableToMap(table)
	if err != nil {
		return "", err
	}

	s, err := MapToLuaCode(maps, "\t")
	if err != nil {
		return "", err
	}

	return s, nil
}

func MapToLuaCode(m map[string]interface{}, prefix string) (string, error) {
	if len(m) <= 0 {
		return "{}", nil
	}

	var buf bytes.Buffer

	buf.WriteString("{\n")
	for key, val := range m {
		// handle key
		buf.WriteString(prefix)
		if StringIsLetter(key) {
			buf.WriteString(key)
		} else {
			buf.WriteString(`["`)
			buf.WriteString(key)
			buf.WriteString(`"]`)
		}
		buf.WriteString(" = ")
		// handle value
		if s, ok := val.(string); ok {
			buf.WriteByte('"')
			buf.WriteString(s)
			buf.WriteByte('"')
		} else if f, ok := val.(float64); ok {
			buf.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
		} else if val == nil {
			buf.WriteString("nil")
		} else if b, ok := val.(bool); ok {
			if b {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		} else if mm, ok := val.(map[string]interface{}); ok {
			str, err := MapToLuaCode(mm, prefix+"\t")
			if err != nil {
				return "", err
			}
			buf.WriteString(str)
		} else {
			return "", errors.New("table value only supported type of bool|number|string|nil|table")
		}
		buf.WriteByte(',')
		buf.WriteByte('\n')
	}
	buf.WriteString(prefix)
	buf.WriteByte('}')
	return buf.String(), nil
}

func ShellExec(cmd string) (string, error) {
	var out bytes.Buffer

	command := exec.Command("/bin/bash", "-c", cmd)
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
