import { srpc } from "./srpc";
import * as grpc_1 from "@grpc/grpc-js";

let client = new srpc.SillyGirlServiceClient(
  "localhost:50051",
  grpc_1.credentials.createInsecure()
);

let senders: Sender[] = [];
let plugin_id = process.env?.PLUGIN_ID ?? "";

process.on("beforeExit", () => {
  for (let sender of senders) {
    sender.destructor();
  }
});

class Sender {
  public uuid: string;
  private destoried = false;
  constructor(uuid: string) {
    this.uuid = uuid;
    senders.push(this);
  }

  destructor() {
    if (this.destoried) return;
    this.destoried = true;
    client.SenderDestroy(
      new srpc.ReplyRequest({ uuid: sender.uuid }),
      (err, resp) => {}
    );
  }

  async getUserId(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetUserId(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getUserName(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetUserName(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getChatId(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetChatId(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getChatName(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetChatName(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getMessageId(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetMessageId(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getPlatform(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetPlatform(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getBotId(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetBotId(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async getContent(): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderGetContent(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async param(key: number | string): Promise<string> {
    return new Promise((resolve, reject) => {
      client.SenderParam(
        new srpc.ReplyRequest({
          uuid: this.uuid,
          content: `${key}`,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value ?? "");
          }
        }
      );
    });
  }
  async setContent(content: string): Promise<undefined> {
    return new Promise((resolve, reject) => {
      client.SenderSetContent(
        new srpc.SenderContentRequest({
          uuid: this.uuid,
          content,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(undefined);
          }
        }
      );
    });
  }
  async continue(): Promise<undefined> {
    return new Promise((resolve, reject) => {
      client.SenderContinue(
        new srpc.SenderRequest({
          uuid: this.uuid,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(undefined);
          }
        }
      );
    });
  }
  async getAdapter() {
    return new Adapter({
      bot_id: await this.getBotId(),
      platform: await this.getPlatform(),
    });
  }

  async listen(options?: {
    rules?: string[]; // 匹配规则
    timeout?: number; // 超时，单位毫秒
    handle?: (s: Sender) => Promise<string | void> | string | void; // 如果匹配成功，则进入消息处理逻辑。如果将 holdOn(content) 的结果作为返回值，会继续监听
    listen_private?: boolean; // 监听用户群内消息时，同时监听用户消息
    listen_group?: boolean; // 监听用户群内消息时，同时监听群员消息
    allow_platforms?: string[]; // 平台白名单
    prohibit_platforms?: string[]; // 平台黑名单
    allow_groups?: string[]; // 群聊白名单
    prohibit_groups?: string[]; // 群聊黑名单
    allow_users?: string[]; // 用户白名单
    prohibit_users?: string[]; // 群聊白名单
    persistent?: boolean; //持久化监听
  }): Promise<Sender | undefined> {
    return new Promise(
      async (resolve, reject) => {
        let params: any = {
          uuid: this.uuid,
          rules: options?.rules,
          timeout: options?.timeout,
          listen_private: options?.listen_private,
          listen_group: options?.listen_group,
          allow_platforms: options?.allow_platforms,
          prohibit_platforms: options?.prohibit_platforms,
          allow_groups: options?.allow_groups,
          prohibit_groups: options?.prohibit_groups,
          allow_users: options?.allow_users,
          prohibit_users: options?.prohibit_users,
          persistent: options?.persistent,
          plugin_id,
        };
        if (!this.uuid) {
          params.persistent = true;
        }
        const call = client.SenderListen();
        // let callback: any = options.replyHandler;
        // console.log("===",this.uuid)
        call.on("data", (response) => {
          if (response.echo == "END") {
            call.cancel();
            return;
          }
          let s = response.uuid ? new Sender(response.uuid) : undefined;
          if (options?.handle && s) {
            // console.log(`options?.handle`, options.persistent)
            let obj = options?.handle(s);
            if (typeof obj == "string") {
              call.write(
                new srpc.SenderListenRequest({
                  uuid: response.echo,
                  value: obj,
                })
              );
            } else if (obj) {
              obj
                .then((v: any) => {
                  call.write(
                    new srpc.SenderListenRequest({
                      uuid: response.echo,
                      value: v ?? "",
                    })
                  );
                })
                .catch((e) => {
                  call.write(
                    new srpc.SenderListenRequest({
                      uuid: response.echo,
                      value: "",
                    })
                  );
                });
            } else {
              call.write(
                new srpc.SenderListenRequest({
                  uuid: response.echo,
                  value: "",
                })
              );
            }
            // console.log(`options?.handle`, options.persistent)
          } else {
            // console.log(`call.cancel()`, options?.persistent)
            call.write(
              new srpc.SenderListenRequest({
                uuid: response.echo,
                value: "",
              })
            );
            // console.log(`call.cancel()`, options?.persistent)
          }
          // console.log("response", JSON.stringify(response));
          // call.cancel()
          resolve(s);
        });
        call.on("error", (err) => {
          reject(err);
          // console.error(err);
        });
        // console.log("params", JSON.stringify(params));
        call.write(new srpc.SenderListenRequest(params));
      }
      // client.SenderListen(new srpc.SenderListenRequest(params), (err, resp) => {
      //   if (err) {
      //     reject(err);
      //   } else {
      //     if (resp?.value) {
      //       let handle = options?.handle;
      //       let s = new Sender(resp.value);
      //       resolve(s);
      //       if (handle) {
      //         handle(s);
      //       }
      //     } else {
      //       reject(new Error("timeout"));
      //     }
      //   }
      // });
      // }
    );
  }
  holdOn(str: string) {
    return "go_again_" + str;
  }
  async reply(content: string): Promise<string | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderReply(
        new srpc.ReplyRequest({
          uuid: this.uuid,
          content,
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async action(options: any): Promise<any | undefined> {
    return new Promise((resolve, reject) => {
      client.SenderAction(
        new srpc.ReplyRequest({
          uuid: this.uuid,
          content: JSON.stringify(options),
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value);
          }
        }
      );
    });
  }
  async event(): Promise<any | undefined> {
    return;
  }
}

class Bucket {
  name: string;

  constructor(name: string) {
    this.name = name;
  }

  transform(v: string | undefined) {
    if (!v) {
      return undefined;
    }
    let result: number | boolean;
    if (v.startsWith("f:")) {
      result = parseFloat(v.replace("f:", ""));
      return result;
    }
    if (v.startsWith("d:")) {
      result = parseInt(v.replace("d:", ""));
      return result;
    }
    if (v.startsWith("b:")) {
      result = v.replace("b:", "") === "true";
      return result;
    }
    if (v.startsWith("o:")) {
      result = JSON.parse(v.replace("o:", ""));
      return result;
    }
    return v;
  }

  reverseTransform(value: any): string {
    if (typeof value === "number" || typeof value === "boolean") {
      return value.toString();
    }
    if (typeof value === "object" && value !== null) {
      return "o:" + JSON.stringify(value);
    }
    if (value === undefined || value === null) {
      return "";
    }
    if (typeof value === "string") {
      if (!value) {
        return "";
      }
      if (!isNaN(parseFloat(value))) {
        return "f:" + parseFloat(value);
      }
      if (!isNaN(parseInt(value))) {
        return "d:" + parseInt(value);
      }
      if (value === "true" || value === "false") {
        return "b:" + (value === "true");
      }
    }
    return value;
  }

  async get(key: string, defaultValue: any = undefined): Promise<any> {
    return new Promise((resolve, reject) => {
      client.BucketGet(
        new srpc.BucketKeyRequest({ name: this.name, key }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(this.transform(resp?.value) || defaultValue);
          }
        }
      );
    });
  }

  async set(
    key: string,
    value: any
  ): Promise<{ message?: string; changed?: boolean }> {
    return new Promise((resolve, reject) => {
      client.BucketSet(
        new srpc.BucketSetRequest({
          name: this.name,
          key,
          value: this.reverseTransform(value),
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve({
              message: resp?.message,
              changed: resp?.changed,
            });
          }
        }
      );
    });
  }

  async getAll(): Promise<any> {
    return new Promise((resolve, reject) => {
      client.BucketGetAll(
        new srpc.BucketRequest({ name: this.name }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            let values: any = {};
            if (resp?.value) {
              values = JSON.parse(resp?.value);
              for (let key in values) {
                values[key] = this.transform(values[key]);
              }
            }
            resolve(values);
          }
        }
      );
    });
  }

  async delete(): Promise<undefined> {
    return new Promise((resolve, reject) => {
      client.BucketDelete(
        new srpc.BucketRequest({ name: this.name }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(undefined);
          }
        }
      );
    });
  }

  async keys(): Promise<string[] | undefined> {
    return new Promise((resolve, reject) => {
      client.BucketKeys(
        new srpc.BucketRequest({ name: this.name }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.keys);
          }
        }
      );
    });
  }

  async len(): Promise<number | undefined> {
    return new Promise((resolve, reject) => {
      client.BucketLen(
        new srpc.BucketRequest({ name: this.name }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.length);
          }
        }
      );
    });
  }

  async buckets(): Promise<string[] | undefined> {
    return new Promise((resolve, reject) => {
      client.BucketBuckets(new srpc.Empty(), (err, resp) => {
        if (err) {
          reject(err);
        } else {
          resolve(resp?.buckets);
        }
      });
    });
  }

  watch(
    key: string,
    handle: (old: any, now: any, key: string) => StorageFinal | void | any
  ) {
    const call = client.BucketWatch();
    call.on("data", async (response) => {
      let fin: any = handle(
        this.transform(response.old),
        this.transform(response.now),
        response.key
      );
      fin = await fin;
      let result: StorageFinal = {
        echo: response.echo,
      };
      if (!fin) {
        result.error = "VOID";
      } else {
        result.now = this.reverseTransform(fin.now);
        result.message = fin.message;
        result.error = fin.error;
      }
      call.write(new srpc.BucketWatchRequest(result));
    });
    call.on("error", (err) => {
      // console.error(err);
    });
    call.write(
      new srpc.BucketWatchRequest({
        name: this.name,
        key: key,
        plugin_id,
      })
    );
  }

  async _name(): Promise<string> {
    return this.name;
  }
}
interface StorageFinal {
  echo?: string;
  now?: any;
  message?: string;
  error?: string;
}

interface Message {
  message_id?: string; // 消息ID
  user_id: string; // 用户ID
  chat_id?: string; // 聊天ID
  content: string; // 聊天内容
  user_name?: string; // 用户名
  chat_name?: string; // 群组名
}

class Adapter {
  platform: string | undefined;
  bot_id: string | undefined;
  call: any;
  constructor(options: {
    platform?: string;
    bot_id?: string;
    replyHandler?: (
      message: Message
    ) => string | undefined | Promise<string | undefined>;
    actionHandler?: (
      message: Message
    ) => string | undefined | Promise<string | undefined>;
  }) {
    this.platform = options.platform;
    this.bot_id = options.bot_id;
    if (options.replyHandler) {
      const call = client.AdapterRegist();
      // let callback: any = ;
      call.on("data", async (response) => {
        // console.log("start on data")
        let message = JSON.parse(response.value);
        const { echo, __type__ } = message;
        delete message.__type__;
        delete message.echo;
        if (__type__ == "reply" && options.replyHandler) {
          let v = (await options.replyHandler(message)) ?? "";
          call.write(
            new srpc.AdapterRegistRequest({
              bot_id: echo,
              platform: v,
            })
          );
        }
        if (__type__ == "action" && options.actionHandler) {
          let v = await options.actionHandler(message);
          call.write(
            new srpc.AdapterRegistRequest({
              bot_id: echo,
              platform: v,
            })
          );
        }
        // console.log("end on data")
      });
      call.on("error", (err) => {
        console.error("adapter disc", err);
      });
      // console.log("before write")
      call.write(
        new srpc.AdapterRegistRequest({
          bot_id: options.bot_id,
          platform: options.platform,
        })
      );
      // console.log("after write write")
      this.call = call;
    }
  }
  setActionHandler(func: (action: {}) => any): void {
    // 将从服务端不断接收action消息，并处理
    // 事件处理巨饼
  }
  async receive(message: Message): Promise<Sender> {
    //投递消息
    return new Promise<Sender>((resolve, reject) => {
      client.AdapterReceive(
        new srpc.AdapterRequest({
          platform: this.platform,
          bot_id: this.bot_id,
          value: JSON.stringify(message),
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else if (resp?.value) {
            resolve(new Sender(resp.value));
          }
        }
      );
    });
  }
  async push(message: Message): Promise<string> {
    //推送消息
    return new Promise<string>((resolve, reject) => {
      client.AdapterPush(
        new srpc.AdapterRequest({
          platform: this.platform,
          bot_id: this.bot_id,
          value: JSON.stringify(message),
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else {
            resolve(resp?.value ?? "");
          }
        }
      );
    });
  }
  async destroy(): Promise<void> {
    this.call.cancel();
  }
  async sender(options: any): Promise<Sender> {
    return new Promise<Sender>((resolve, reject) => {
      client.AdapterSender(
        new srpc.AdapterRequest({
          platform: this.platform,
          bot_id: this.bot_id,
          value: JSON.stringify(options),
        }),
        (err, resp) => {
          if (err) {
            reject(err);
          } else if (resp?.value) {
            resolve(new Sender(resp.value));
          }
        }
      );
    });
  }
}

let sender: Sender = new Sender(process.env?.SENDER_ID ?? "");

async function sleep(ms: number | undefined) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

class Console {
  error = (message?: any, ...optionalParams: any[]) => {};
  info = (message?: any, ...optionalParams: any[]) => {};
  log = (message?: any, ...optionalParams: any[]) => {};
  debug = (message?: any, ...optionalParams: any[]) => {};
}

interface CQItem {
  type: string;
  params: {};
}

let utils = {
  parseCQText: (text: string, prefix = "CQ") => {
    const cqRegex = new RegExp(`\\[${prefix}:(\\w+)(.*?)\\]`, "g");
    const cqMatches = text.matchAll(cqRegex);
    const result: (CQItem | string)[] = [];

    let lastIndex = 0;
    for (const match of cqMatches) {
      // 添加 CQ 码前的文本
      const matchIndex = text.indexOf(match[0], lastIndex);
      if (matchIndex > lastIndex) {
        result.push(text.slice(lastIndex, matchIndex));
      }

      // 解析 CQ 码
      const params: any = {};
      const paramRegex = /(\w+)=([^,]+)/g;
      const paramMatches = match[2].matchAll(paramRegex);
      for (const paramMatch of paramMatches) {
        params[paramMatch[1]] = paramMatch[2].trim();
      }
      result.push({
        type: match[1],
        params: params,
      });

      lastIndex = matchIndex + match[0].length;
    }
    if (lastIndex < text.length) {
      result.push(text.slice(lastIndex));
    }
    return result;
  },
};

let slog = (type: string, ...args: any[]) => {};

let console = {
  log(...args: any[]) {
    const content = args.reduce((acc, arg) => acc + " " + arg, "");
    client.Console(
      new srpc.ConsoleRequest({ type: "log", content, plugin_id }),
      (err, resp) => {}
    );
  },
  info(...args: any[]) {
    const content = args.reduce((acc, arg) => acc + " " + arg, "");
    client.Console(
      new srpc.ConsoleRequest({ type: "info", content, plugin_id }),
      (err, resp) => {}
    );
  },
  error(...args: any[]) {
    const content = args.reduce((acc, arg) => acc + " " + arg, "");
    client.Console(
      new srpc.ConsoleRequest({ type: "error", content, plugin_id }),
      (err, resp) => {}
    );
  },
  debug(...args: any[]) {
    const content = args.reduce((acc, arg) => acc + " " + arg, "");
    client.Console(
      new srpc.ConsoleRequest({ type: "debug", content, plugin_id }),
      (err, resp) => {}
    );
  },
};

export { Adapter, Bucket, sender, sleep, utils, console };
