import { srpc } from "./srpc";
import * as grpc_1 from "@grpc/grpc-js";

let client = new srpc.SillyGirlServiceClient(
  "localhost:50051",
  grpc_1.credentials.createInsecure()
);

class Sender {
  private uuid: string;
  constructor(uuid: string) {
    this.uuid = uuid;
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
    handle?: (s: Sender) => string | void; // 如果匹配成功，则进入消息处理逻辑。如果将 holdOn(content) 的结果作为返回值，会继续监听
    listen_private?: boolean; // 监听用户群内消息时，同时监听用户消息
    listen_group?: boolean; // 监听用户群内消息时，同时监听群员消息
    allow_platforms?: string[]; // 平台白名单
    prohibit_platforms?: string[]; // 平台黑名单
    allow_groups?: string[]; // 群聊白名单
    prohibit_groups?: string[]; // 群聊黑名单
    allow_users?: string[]; // 用户白名单
    prohibit_users?: string[]; // 群聊白名单
  }): Promise<Sender> {
    return new Promise((resolve, reject) => {
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
      };
      client.SenderListen(new srpc.SenderListenRequest(params), (err, resp) => {
        if (err) {
          reject(err);
        } else {
          if (resp?.value) {
            resolve(new Sender(resp.value));
          } else {
            reject(new Error("timeout"));
          }
        }
      });
    });
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
    return;
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

  reverseTransform(value: any) {
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
        return undefined;
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

  async _name(): Promise<string> {
    return this.name;
  }
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
  call: any
  constructor(options: {
    platform?: string;
    bot_id?: string;
    replyHandler?: (message: Message) => string | undefined;
  }) {
    this.platform = options.platform;
    this.bot_id = options.bot_id;
    const call = client.AdapterRegist();
    if (options.replyHandler) {
      let callback: any = options.replyHandler;
      call.on("data", (response) => {
        let message = JSON.parse(response.value);
        const { echo, __type__ } = message;
        delete(message.__type__)
        delete(message.echo)
        if (__type__ == "reply") {
          call.write(
            new srpc.AdapterRegistRequest({
              bot_id: echo,
              platform: callback(message),
            })
          );
        }
      });
      call.on("error", (err) => {
        // console.error(err);
      });
      call.write(
        new srpc.AdapterRegistRequest({
          bot_id: options.bot_id,
          platform: options.platform,
        })
      );
      this.call = call
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
          } else if (resp?.value) {
            resolve(resp.value);
          }
        }
      );
    });
  }
  async destroy(): Promise<void> {
    this.call.cancel()
  }
  async sender(options: any): Promise<Sender> {
    return new Promise<Sender>((resolve, reject) => {
      client.AdapterSender(
        new srpc.AdapterRequest({
          platform: this.platform,
          bot_id: this.bot_id,
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

let sender = new Sender(
  process.env?.SENDER_ID ?? "4d6371a8-2778-11ee-a3c2-821680fbbf6b"
);

async function sleep(ms: number | undefined) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export { Adapter, Bucket, sender, sleep };
