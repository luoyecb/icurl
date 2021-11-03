package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/luoyecb/icurl/lualib"

	"github.com/luoyecb/eflag"
	"github.com/peterh/liner"
	"github.com/yuin/gopher-lua"
)

const (
	DEFAULT_HISTORY_FILE = "~/.icurl_history"
	DEFAULT_PROMPT       = "icurl> "
	LOGO_PROMPT          = `
  _____     ____   __    __   ______     _____
 (_   _)   / ___)  ) )  ( (  (   __ \   (_   _)
   | |    / /     ( (    ) )  ) (__) )    | |
   | |   ( (       ) )  ( (  (    __/     | |
   | |   ( (      ( (    ) )  ) \ \  _    | |   __
  _| |__  \ \___   ) \__/ (  ( ( \ \_)) __| |___) )
 /_____(   \____)  \______/   )_) \__/  \________/

You can get help information through the help() function.`
)

func OpenHistoryFile() *os.File {
	f, err := os.OpenFile(lualib.GetRealPath(DEFAULT_HISTORY_FILE), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		ErrExit(err)
	}
	return f
}

func ErrExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "system error: %v", err)
		os.Exit(1)
	}
}

func CheckExit(line string) bool {
	if line == "exit" || line == "quit" {
		fmt.Printf("Bye.\n")
		return true
	}
	return false
}

// <Tab>键补全命令
func CommandCompleter(line string) []string {
	candidates := make([]string, 0)
	for name, _ := range lualib.FuncsMap {
		if strings.HasPrefix(name, line) {
			candidates = append(candidates, name)
		}
	}
	return candidates
}

func main() {
	commandOptions := &CommandOptions{}
	eflag.Parse(commandOptions)

	// Lua VM
	vm := lua.NewState()
	defer vm.Close()
	ErrExit(lualib.Init(vm))
	RunWithCommandOptions(vm, commandOptions)

	historyFile := OpenHistoryFile()
	defer historyFile.Close()

	// Readline
	lineState := liner.NewLiner()
	defer lineState.Close()

	lineState.SetCtrlCAborts(true)
	lineState.SetMultiLineMode(true)
	lineState.SetTabCompletionStyle(liner.TabPrints)
	lineState.SetCompleter(CommandCompleter)
	_, err := lineState.ReadHistory(historyFile)
	if err != nil {
		ErrExit(err)
	}

	// Main loop
	fmt.Println(LOGO_PROMPT)
	for {
		line, err := lineState.Prompt(DEFAULT_PROMPT)
		if err == liner.ErrPromptAborted {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		} else if CheckExit(line) {
			break
		}

		lineState.AppendHistory(line)

		// Run shell command
		if line[0] == '!' {
			str, err := lualib.ShellExec(line[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "run shell command error: %v\n", err)
			} else {
				fmt.Println(str)
			}
			continue
		}

		err = lualib.RunLuaCode(vm, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "run lua error: %v\n", err)
		}
	}

	// Truncate historyFile
	historyFile.Truncate(0)
	historyFile.Seek(0, os.SEEK_SET)

	lineState.WriteHistory(historyFile)
}
