package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// Command 普通用法（一次性获取所有输出）
func Command(cmd string) error {
	// bash -c "cmd string"
	c := exec.Command("bash", "-c", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}

// ExecShell 阻塞方式(需要执行结果):阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func ExecShell(commandName string, params []string, logPath string) (string, error) {
	s := commandName
	for _, p := range params {
		s += " " + p
	}
	fmt.Println(s)
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s, ">>", logPath)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	return out.String(), err
}

// ExecCommand 阻塞方式(需要执行结果):需要对shell标准输出的逐行实时进行处理的
func ExecCommand(commandName string, params []string) bool {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	fmt.Println(cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()
	return true
}

// 非阻塞方式(不需要执行结果):不需要执行命令的结果与成功与否，执行命令马上就返回
func ExecShellNoResult(command string) {
	////处理启动参数，通过空格分离 如：setsid /home/luojing/gotest/src/test_main/iwatch/test/while_little &
	//command_name_and_args := strings.FieldsFunc(command, splite_command)
	////开始执行c包含的命令，但并不会等待该命令完成即返回
	//cmd.Start()
	//if err != nil {
	//	fmt.Printf("%v: exec command:%v error:%v\n", get_time(), command, err)
	//}
	//fmt.Printf("Waiting for command:%v to finish...\n", command)
	////阻塞等待fork出的子进程执行的结果，和cmd.Start()配合使用[不等待回收资源，会导致fork出执行shell命令的子进程变为僵尸进程]
	//err = cmd.Wait()
	//if err != nil {
	//	fmt.Printf("%v: Command finished with error: %v\n", get_time(), err)
	//}
	//return
}
