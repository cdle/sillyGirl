package core

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	cron "github.com/robfig/cron/v3"
)

func init() {
	go initNodePlugins()
}

var processes sync.Map

func initNodePlugins() {
	root := utils.ExecPath + "/plugins"
	plugins := []string{root}
	os.Mkdir(root, 0755)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !strings.HasPrefix(path, utils.ExecPath+"/plugins/") {
			return nil
		}
		files := strings.Split(strings.Replace(path, root+"/", "", 1), "/")
		// var plugin_dir = false
		// var plugin_index = false
		switch len(files) {
		case 1:
			// plugin_dir = true
			if info.IsDir() {
				plugins = append(plugins, path)
			}
		case 2:
			if (files[1] == "main.js") && !info.IsDir() { //files[1] == "main.ts" ||
				AddNodePlugin(path, files[0])
			}
		}
		return nil
	})
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("创建监视器失败：", err)
		return
	}
	defer watcher.Close()
	// 要监控的文件夹路径
	for _, dir := range plugins {
		err = watcher.Add(dir)
		if err != nil {
			fmt.Println("添加监视目录失败：", err)
			return
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// fmt.Println(event.Name, "op", event.Op.String())
			files := strings.Split(strings.Replace(event.Name, root+"/", "", 1), "/")
			var plugin_dir = false
			var plugin_index = false
			var plugin_name = ""
			switch len(files) {
			case 1:
				plugin_dir = true
				// fmt.Println("目录事件")
				plugin_name = files[0]
			case 2:
				if files[1] == "main.ts" || files[1] == "main.js" {
					if files[1] == "main.js" {
						plugin_index = true
					}
					// fmt.Println("入口文件事件")
				}
				plugin_name = files[0]
			}
			switch event.Op.String() {
			case "CREATE":
				if plugin_dir {
					info, err := os.Stat(event.Name)
					// fmt.Println(err)
					if err == nil && info.IsDir() {
						watcher.Add(event.Name)
						// fmt.Println("增加插件目录", event.Name)
					} else {
						// fmt.Println("非插件目录", event.Name)
					}
				} else if plugin_index {
					// fmt.Println("增加插件", event.Name)
					AddNodePlugin(event.Name, plugin_name)
				}
			case "REMOVE", "RENAME", "REMOVE|RENAME":
				if plugin_dir {
					watcher.Remove(event.Name)
					// fmt.Println("移除插件目录", event.Name)
					// fmt.Println("移除插件", plugin_name)
					RemNodePlugin(plugin_name)
				} else if plugin_index {
					// fmt.Println("移除插件", plugin_name)
					RemNodePlugin(plugin_name)
				}
			case "WRITE":
				if plugin_index {
					RemNodePlugin(plugin_name)
					AddNodePlugin(event.Name, plugin_name)
					// fmt.Println("变更插件", event.Name, plugin_name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("错误：", err)
		}
	}
}

func RemNodePlugin(name string) error {
	pluginLock.Lock()
	defer pluginLock.Unlock()
	key := nameUuid(name)
	// fmt.Println("rem", key, name)
	for i := range Functions {
		if Functions[i].UUID == key {
			// fmt.Println("pl", key)
			DestroyAdapterByUUID(key)
			Functions[i].Running = false
			if len(Functions[i].CronIds) != 0 {
				for _, id := range Functions[i].CronIds {
					CRON.Remove(cron.EntryID(id))
				}
			}
			Functions = append(Functions[:i], Functions[i+1:]...)
			CancelPluginCrons(key)
			CancelPluginWebs(key)
			CancelPluginlistening(key)
			CancelHttpListen(key)
			remStatic(key)
			storage.DisableHandle(key)
			break
		}
	}
	return nil
}

func nameUuid(name string) string {
	hash := sha1.Sum([]byte(name))
	return uuid.NewSHA1(uuid.Nil, hash[:]).String()
}

func AddNodePlugin(path, name string) error {
	pluginLock.Lock()
	defer pluginLock.Unlock()
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	uuid := nameUuid(name)
	// fmt.Println("add,", uuid, name)
	f, cbs := pluginParse(string(data), uuid)
	f.Suffix = ".ts"
	f.Type = "typescript"
	f.Handle = func(s common.Sender, f func(vm *goja.Runtime)) interface{} {
		s.SetPluginID(uuid)
		plt := s.GetImType()
		// , "/home/user/.nvm/versions/node/v18.16.1/lib/node_modules/ts-node/dist/bin.js",
		cmd := exec.Command(utils.ExecPath+"/language/node/node", path)
		// cmd := exec.Command(utils.ExecPath+"/language/node/node", path)
		id := s.SetID()
		cmd.Env = append(cmd.Env, "SENDER_ID="+id)
		cmd.Env = append(cmd.Env, "PLUGIN_ID="+uuid)

		// 获取标准输出和标准错误输出的管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			// fmt.Printf("获取标准输出管道失败：%v\n", err)
			// os.Exit(1)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			// fmt.Printf("获取标准错误输出管道失败：%v\n", err)
			// os.Exit(1)
		}
		err = cmd.Start()
		if err != nil {

		}

		// file, err := os.Create("output.log")
		// if err != nil {
		// 	fmt.Printf("创建文件失败：%v\n", err)
		// 	os.Exit(1)
		// }
		// defer file.Close()
		var wg sync.WaitGroup
		wg.Add(2)
		// 处理标准输出
		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				data := scanner.Text()
				fmt.Println("log", data)
				// if _, err := file.WriteString(data + "\n"); err != nil {
				// 	fmt.Printf("写入文件失败：%v\n", err)
				// }
			}
		}()
		// 处理标准错误输出
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				data := scanner.Text()
				fmt.Fprintln(os.Stderr, "err "+data)
				// if _, err := file.WriteString(data + "\n"); err != nil {
				// 	fmt.Printf("写入文件失败：%v\n", err)
				// }
			}
		}()
		processes.Store(cmd, s)
		if (plt) != "*" {
			senders.Store(id, s)
			defer senders.Delete(id)
			defer processes.Delete(cmd)
			err = cmd.Wait()
			if err != nil {
				fmt.Println("命令执行失败：", err)
				return nil
			}
		} else {
			processes.Range(func(key, value any) bool {
				p := key.(*exec.Cmd)
				if p == cmd {
					return true
				}
				s := value.(common.Sender)
				if s.GetPluginID() == uuid {
					p.Process.Kill()
				}
				return true
			})
			go func() {
				defer processes.Delete(cmd)
				err = cmd.Wait()
			}()
		}
		return nil
	}
	for _, cb := range cbs {
		cb()
	}
	AddCommand([]*common.Function{f})
	return nil
}
