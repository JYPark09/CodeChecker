package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func runHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	var lang *Language
	typeStr := request.PostForm.Get("lang")

	for i := range langMgr.Languages {
		if langMgr.Languages[i].TypeStr == typeStr {
			lang = &langMgr.Languages[i]
			break
		}
	}

	if lang == nil {
		log.Println("[Run Handler] Invalid lang:", lang)
		return
	}

	id := uuid.New().String()

	log.Println("[Run Handler] New analyze start id:", id)

	directory := "results/" + id
	os.MkdirAll(directory, 0755)

	input := request.PostForm.Get("stdin")

	srcFile, err := os.Create(directory + "/" + "source" + lang.Extension)
	if err != nil {
		log.Println("[Run Hander] Cannot create source file")
		return
	}
	srcFile.WriteString(request.PostForm.Get("source"))
	srcFile.Close()

	stdinFile, err := os.Create(directory + "/stdin")
	if err != nil {
		log.Println("[Run Handler] Cannot create source file")
		return
	}
	stdinFile.WriteString(input)
	stdinFile.Close()

	compileSuccess := true
	runSuccess := true
	runOutput := ""

	for i, cmdLine := range lang.CmdLines {
		compiler := exec.Command(cmdLine[0], cmdLine[1:]...)
		compiler.Dir = directory

		compileResult, err := os.Create(directory + "/compile_" + strconv.Itoa(i))
		if err != nil {
			log.Println("[Run Handler] Cannot create compileResult")
			continue
		}
		defer compileResult.Close()
		compiler.Stderr = compileResult

		if err = compiler.Start(); err != nil {
			log.Println("[Run Handler] Cannot start compiler")
			continue
		}
		compiler.Wait()

		compileSuccess = compileSuccess && compiler.ProcessState.Success()

		cmd := exec.Command("./out")
		cmd.Dir = directory

		var res bytes.Buffer
		cmd.Stdin = strings.NewReader(input)
		cmd.Stdout = &res

		if err = cmd.Run(); err != nil {
			log.Println("[Run Handler] Cannot start program")
			continue
		}
		cmd.Wait()

		runResult, err := os.Create(directory + "/run_" + strconv.Itoa(i))
		if err != nil {
			log.Println("[Run Handler] Cannot create runResult")
			continue
		}
		defer runResult.Close()
		runResult.WriteString(res.String())

		if i == 0 {
			runOutput = res.String()
		}

		runSuccess = runSuccess && (runOutput == res.String())
	}

	os.Remove(directory + "/out")

	var cmds []string
	for i := 0; i < len(lang.CmdLines); i++ {
		cmds = append(cmds, strings.Join(lang.CmdLines[i], " "))
	}

	var resStr string
	if !compileSuccess {
		resStr = "컴파일 오류"
	} else if !runSuccess {
		resStr = "출력 결과 불일치"
	} else {
		resStr = "정상"
	}

	saveResult(id, RunResult{ID: id, Result: resStr, Runs: cmds, Ext: lang.Extension})

	http.Redirect(writer, request, "/result?id="+id, 301)
}
