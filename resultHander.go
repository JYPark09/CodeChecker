package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

func resultHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	directory := "results/" + id

	type runResult struct {
		Command string
		Compile string
		Run     string
	}

	type resultArg struct {
		ID     string
		Source string
		Stdin  string

		Result  string
		Results []runResult
	}

	result := getResult(id)

	arg := resultArg{ID: id}

	srcFile, err := ioutil.ReadFile(directory + "/source" + result.Ext)
	if err != nil {
		fmt.Fprintln(writer, "파일을 찾을 수 없습니다.")
		return
	}
	stdinFile, err := ioutil.ReadFile(directory + "/stdin")
	if err != nil {
		fmt.Fprintln(writer, "파일을 찾을 수 없습니다.")
		return
	}

	arg.Source = string(srcFile)
	arg.Stdin = string(stdinFile)
	arg.Result = result.Result

	for i := 0; i < len(result.Runs); i++ {
		compileFile, err := ioutil.ReadFile(directory + "/compile_" + strconv.Itoa(i))
		if err != nil {
			fmt.Fprintln(writer, "파일을 찾을 수 없습니다.")
			return
		}

		runFile, err := ioutil.ReadFile(directory + "/run_" + strconv.Itoa(i))
		var runFileStr string
		if err == nil {
			runFileStr = string(runFile)
		}

		res := runResult{Command: result.Runs[i], Compile: string(compileFile), Run: runFileStr}

		arg.Results = append(arg.Results, res)
	}

	t, _ := template.ParseFiles("template/result.html")
	t.Execute(writer, &arg)
}
