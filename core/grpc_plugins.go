package core

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

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
	initLanguage()
	root := strings.ReplaceAll(utils.ExecPath+"/plugins", "\\", "/")
	plugins := []string{root}
	os.Mkdir(root, 0755)
	// fmt.Println("root", root)

	files, _ := ioutil.ReadDir(root)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		name := file.Name()
		path := root + "/" + name
		plugins = append(plugins, path)
		index := path + "/main.js"
		if info, err := os.Stat(index); err == nil && !info.IsDir() {
			AddNodePlugin(index, name)
		}
	}

	// filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
	// 	fmt.Println("path", path)
	// 	path = strings.ReplaceAll(path, "\\", "/")
	// 	if !strings.HasPrefix(path, root+"/") {
	// 		return nil
	// 	}
	// 	files := strings.Split(strings.Replace(path, root+"/", "", 1), "/")
	// 	fmt.Println("files", files)
	// 	// var plugin_dir = false
	// 	// var plugin_index = false
	// 	switch len(files) {
	// 	case 1:
	// 		// plugin_dir = true
	// 		if info.IsDir() {
	// 			plugins = append(plugins, path)
	// 		}
	// 	case 2:
	// 		if (files[1] == "main.js") && !info.IsDir() { //files[1] == "main.ts" ||
	// 			AddNodePlugin(path, files[0])
	// 		}
	// 	}
	// 	return nil
	// })
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
			event.Name = strings.ReplaceAll(event.Name, "\\", "/")
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
			if plugin_name == "." {
				continue
			}
			switch event.Op.String() {
			case "CREATE":
				if plugin_dir {
					info, err := os.Stat(event.Name)
					// fmt.Println(err)
					if err == nil && info.IsDir() {
						event_name := event.Name + "/main.js"
						if info, err := os.Stat(event_name); err == nil && !info.IsDir() {
							AddNodePlugin(event_name, plugin_name)
						} else {
							time.Sleep(time.Millisecond * 100)
							if info, err := os.Stat(event_name); err == nil && !info.IsDir() {
								AddNodePlugin(event_name, plugin_name)
							}
						}
						watcher.Add(event.Name)
						// fmt.Println("增加插件目录", event.Name)
					} else {
						// fmt.Println("非插件目录", event.Name)
					}
					tf := event.Name + "/node_modules/sillygirl.d.ts"
					ti := event.Name + "/main.js"
					if _, err := os.Stat(tf); err != nil {
						os.Mkdir(event.Name+"/node_modules", 0700)
						os.WriteFile(tf, []byte(typeat), 0700)
					}
					go func() {
						time.Sleep(time.Second)
						if _, err := os.Stat(ti); err != nil {
							os.Mkdir(event.Name+"/node_modules", 0700)
							os.WriteFile(ti, []byte(defaultScript(plugin_name)), 0700)
						}
					}()
				} else if plugin_index {
					// fmt.Println("增加插件", event.Name)
					// RemNodePlugin(plugin_name)
					AddNodePlugin(event.Name, plugin_name)
				}
			case "REMOVE", "RENAME", "REMOVE|RENAME", "REMOVE|WRITE":
				if plugin_dir {
					watcher.Remove(event.Name)
					// fmt.Println("移除插件目录", event.Name)
					// fmt.Println("移除插件", plugin_name)
					RemNodePlugin(plugin_name)

				} else if plugin_index {
					// fmt.Println("移除插件", plugin_name)
					RemNodePlugin(plugin_name)
				}
			case "WRITE": //, "CHMOD"
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

func RemNodePlugin(name string) bool {
	if name == "" {
		return false
	}
	pluginLock.Lock()
	defer pluginLock.Unlock()
	key := nameUuid(name)
	for i := range Functions {
		if Functions[i].UUID == key {
			f := Functions[i]
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
			console.Log("已移除 %s%s", f.Title, f.Suffix)
			return true
		}
	}
	return false
}

func nameUuid(name string) string {
	hash := sha1.Sum([]byte(name))
	return strings.ReplaceAll(uuid.NewSHA1(uuid.Nil, hash[:]).String(), "-", "_")
}

func isNameUuid(uuid string) bool {
	return strings.Contains(uuid, "_")
}

// var plugins_id sync.Map

func AddNodePlugin(path, name string) error {
	if name == "" {
		return nil
	}
	uuid := nameUuid(name)
	plugins.Set(uuid, "")
	pluginLock.Lock()
	defer pluginLock.Unlock()
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	script := string(data)
	if script == "" {
		return nil
	}
	// plugins_id.Store(uuid, path)
	// fmt.Println("add,", uuid, name)
	f, cbs := pluginParse(script, uuid)
	f.Reload = func() { //重载
		RemNodePlugin(path)
		AddNodePlugin(path, name)
	}
	f.Suffix = ".js"
	f.Type = "node"
	f.Path = path
	f.Handle = func(s common.Sender, f func(vm *goja.Runtime)) interface{} {
		s.SetPluginID(uuid)
		plt := s.GetImType()
		// , "/home/user/.nvm/versions/node/v18.16.1/lib/node_modules/ts-node/dist/bin.js",
		cmd := exec.Command("./node", path)
		cmd.Dir = utils.ExecPath + "/language/node"
		// cmd := exec.Command(utils.ExecPath+"/language/node/node", path)
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
				fmt.Println(data)
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
				fmt.Println(data)
				// if _, err := file.WriteString(data + "\n"); err != nil {
				// 	fmt.Printf("写入文件失败：%v\n", err)
				// }
			}
		}()
		processes.Store(cmd, s)
		if (plt) != "*" {
			id := s.SetID()
			cmd.Env = append(cmd.Env, "SENDER_ID="+id)
			err = cmd.Start()
			if err != nil {

			}
			senders.Store(id, s)
			defer senders.Delete(id)
			defer processes.Delete(cmd)
			err = cmd.Wait()
			if err != nil {
				fmt.Println("命令执行失败：", err)
				return nil
			}
		} else {
			err = cmd.Start()
			if err != nil {

			}
			processes.Range(func(key, value any) bool {
				p := key.(*exec.Cmd)
				if p == cmd {
					return true
				}
				s := value.(common.Sender)
				if s.GetPluginID() == uuid {
					func() {
						defer func() {
							recover()
						}()
						if p.Process.Kill() == nil {
							processes.Delete(key)
						}
					}()
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
	if !f.OnStart {
		console.Log("已加载 %s%s", f.Title, f.Suffix)
	}
	AddCommand([]*common.Function{f})
	return nil
}

var typeat = `declare class Sender {
	uuid: string;
	private destoried;
	constructor(uuid: string);
	destructor(): void;
	getUserId(): Promise<string | undefined>;
	getUserName(): Promise<string | undefined>;
	getChatId(): Promise<string | undefined>;
	getChatName(): Promise<string | undefined>;
	getMessageId(): Promise<string | undefined>;
	getPlatform(): Promise<string | undefined>;
	getBotId(): Promise<string | undefined>;
	getContent(): Promise<string | undefined>;
	param(key: number | string): Promise<string>;
	setContent(content: string): Promise<undefined>;
	continue(): Promise<undefined>;
	getAdapter(): Promise<Adapter>;
	listen(options?: {
			rules?: string[];
			timeout?: number;
			handle?: (s: Sender) => Promise<string | void> | string | void;
			listen_private?: boolean;
			listen_group?: boolean;
			allow_platforms?: string[];
			prohibit_platforms?: string[];
			allow_groups?: string[];
			prohibit_groups?: string[];
			allow_users?: string[];
			prohibit_users?: string[];
			persistent?: boolean;
	}): Promise<Sender | undefined>;
	holdOn(str: string): string;
	reply(content: string): Promise<string | undefined>;
	action(options: any): Promise<any | undefined>;
	event(): Promise<any | undefined>;
}
declare class Bucket {
	name: string;
	constructor(name: string);
	transform(v: string | undefined): string | number | boolean | undefined;
	reverseTransform(value: any): string;
	get(key: string, defaultValue?: any): Promise<any>;
	set(key: string, value: any): Promise<{
			message?: string;
			changed?: boolean;
	}>;
	getAll(): Promise<any>;
	delete(): Promise<undefined>;
	keys(): Promise<string[] | undefined>;
	len(): Promise<number | undefined>;
	buckets(): Promise<string[] | undefined>;
	watch(key: string, handle: (old: any, now: any, key: string) => StorageFinal | void | any): void;
	_name(): Promise<string>;
}
interface StorageFinal {
	echo?: string;
	now?: any;
	message?: string;
	error?: string;
}
interface Message {
	message_id?: string;
	user_id: string;
	chat_id?: string;
	content: string;
	user_name?: string;
	chat_name?: string;
}
declare class Adapter {
	platform: string | undefined;
	bot_id: string | undefined;
	call: any;
	constructor(options: {
			platform?: string;
			bot_id?: string;
			replyHandler?: (message: Message) => string | undefined | Promise<string | undefined>;
			actionHandler?: (message: Message) => string | undefined | Promise<string | undefined>;
	});
	setActionHandler(func: (action: {}) => any): void;
	receive(message: Message): Promise<Sender>;
	push(message: Message): Promise<string>;
	destroy(): Promise<void>;
	sender(options: any): Promise<Sender>;
}
declare let sender: Sender;
declare function sleep(ms: number | undefined): Promise<unknown>;
interface CQItem {
	type: string;
	params: {};
}
declare let utils: {
	parseCQText: (text: string, prefix?: string) => (string | CQItem)[];
};
declare let console: {
	log(...args: any[]): void;
	info(...args: any[]): void;
	error(...args: any[]): void;
	debug(...args: any[]): void;
};
export { Adapter, Bucket, sender, sleep, utils, console };

`

func defaultScript(title string) string {
	create_at := time.Now().Format("2006-01-02 15:04:05")
	return `/**
	* @title ` + title + `
	* @create_at ` + create_at + `
	* @description 🐒这个人很懒什么都没有留下
	* @author ` + sillyGirl.GetString("author", "佚名") + `
	* @version v1.0.0
	*/

	const { sender: s, Bucket, Adapter, sleep } = require("sillygirl");`
}
