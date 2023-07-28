"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.console = exports.utils = exports.sleep = exports.sender = exports.Bucket = exports.Adapter = void 0;
const srpc_1 = require("./srpc");
const grpc_1 = __importStar(require("@grpc/grpc-js"));
let client = new srpc_1.srpc.SillyGirlServiceClient("localhost:50051", grpc_1.credentials.createInsecure());
let senders = [];
let plugin_id = process.env?.PLUGIN_ID ?? "";
process.on("beforeExit", () => {
    for (let sender of senders) {
        sender.destructor();
    }
});
class Sender {
    uuid;
    destoried = false;
    constructor(uuid) {
        this.uuid = uuid;
        senders.push(this);
    }
    destructor() {
        if (this.destoried)
            return;
        this.destoried = true;
        client.SenderDestroy(new srpc_1.srpc.ReplyRequest({ uuid: sender.uuid }), (err, resp) => { });
    }
    async getUserId() {
        return new Promise((resolve, reject) => {
            client.SenderGetUserId(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getUserName() {
        return new Promise((resolve, reject) => {
            client.SenderGetUserName(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getChatId() {
        return new Promise((resolve, reject) => {
            client.SenderGetChatId(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getChatName() {
        return new Promise((resolve, reject) => {
            client.SenderGetChatName(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getMessageId() {
        return new Promise((resolve, reject) => {
            client.SenderGetMessageId(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getPlatform() {
        return new Promise((resolve, reject) => {
            client.SenderGetPlatform(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getBotId() {
        return new Promise((resolve, reject) => {
            client.SenderGetBotId(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async getContent() {
        return new Promise((resolve, reject) => {
            client.SenderGetContent(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async param(key) {
        return new Promise((resolve, reject) => {
            client.SenderParam(new srpc_1.srpc.ReplyRequest({
                uuid: this.uuid,
                content: `${key}`,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value ?? "");
                }
            });
        });
    }
    async setContent(content) {
        return new Promise((resolve, reject) => {
            client.SenderSetContent(new srpc_1.srpc.SenderContentRequest({
                uuid: this.uuid,
                content,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(undefined);
                }
            });
        });
    }
    async continue() {
        return new Promise((resolve, reject) => {
            client.SenderContinue(new srpc_1.srpc.SenderRequest({
                uuid: this.uuid,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(undefined);
                }
            });
        });
    }
    async getAdapter() {
        return new Adapter({
            bot_id: await this.getBotId(),
            platform: await this.getPlatform(),
        });
    }
    async listen(options) {
        return new Promise(async (resolve, reject) => {
            let params = {
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
                        call.write(new srpc_1.srpc.SenderListenRequest({
                            uuid: response.echo,
                            value: obj,
                        }));
                    }
                    else if (obj) {
                        obj
                            .then((v) => {
                            call.write(new srpc_1.srpc.SenderListenRequest({
                                uuid: response.echo,
                                value: v ?? "",
                            }));
                        })
                            .catch((e) => {
                            call.write(new srpc_1.srpc.SenderListenRequest({
                                uuid: response.echo,
                                value: "",
                            }));
                        });
                    }
                    else {
                        call.write(new srpc_1.srpc.SenderListenRequest({
                            uuid: response.echo,
                            value: "",
                        }));
                    }
                    // console.log(`options?.handle`, options.persistent)
                }
                else {
                    // console.log(`call.cancel()`, options?.persistent)
                    call.write(new srpc_1.srpc.SenderListenRequest({
                        uuid: response.echo,
                        value: "",
                    }));
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
            call.write(new srpc_1.srpc.SenderListenRequest(params));
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
    holdOn(str) {
        return "go_again_" + str;
    }
    async reply(content) {
        return new Promise((resolve, reject) => {
            client.SenderReply(new srpc_1.srpc.ReplyRequest({
                uuid: this.uuid,
                content,
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async action(options) {
        return new Promise((resolve, reject) => {
            client.SenderAction(new srpc_1.srpc.ReplyRequest({
                uuid: this.uuid,
                content: JSON.stringify(options),
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value);
                }
            });
        });
    }
    async event() {
        return;
    }
}
class Bucket {
    name;
    constructor(name) {
        this.name = name;
    }
    transform(v) {
        if (!v) {
            return undefined;
        }
        let result;
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
    reverseTransform(value) {
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
    async get(key, defaultValue = undefined) {
        return new Promise((resolve, reject) => {
            client.BucketGet(new srpc_1.srpc.BucketKeyRequest({ name: this.name, key }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(this.transform(resp?.value) || defaultValue);
                }
            });
        });
    }
    async set(key, value) {
        return new Promise((resolve, reject) => {
            client.BucketSet(new srpc_1.srpc.BucketSetRequest({
                name: this.name,
                key,
                value: this.reverseTransform(value),
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve({
                        message: resp?.message,
                        changed: resp?.changed,
                    });
                }
            });
        });
    }
    async getAll() {
        return new Promise((resolve, reject) => {
            client.BucketGetAll(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    let values = {};
                    if (resp?.value) {
                        values = JSON.parse(resp?.value);
                        for (let key in values) {
                            values[key] = this.transform(values[key]);
                        }
                    }
                    resolve(values);
                }
            });
        });
    }
    async delete(key) {
        return this.set(key, "");
    }
    async deleteAll() {
        return new Promise((resolve, reject) => {
            client.BucketDelete(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(undefined);
                }
            });
        });
    }
    async keys() {
        return new Promise((resolve, reject) => {
            client.BucketKeys(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.keys);
                }
            });
        });
    }
    async len() {
        return new Promise((resolve, reject) => {
            client.BucketLen(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.length);
                }
            });
        });
    }
    async buckets() {
        return new Promise((resolve, reject) => {
            client.BucketBuckets(new srpc_1.srpc.Empty(), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.buckets);
                }
            });
        });
    }
    watch(key, handle) {
        const call = client.BucketWatch();
        call.on("data", async (response) => {
            let fin = handle(this.transform(response.old), this.transform(response.now), response.key);
            fin = await fin;
            let result = {
                echo: response.echo,
            };
            if (!fin) {
                result.error = "VOID";
            }
            else {
                result.now = this.reverseTransform(fin.now);
                result.message = fin.message;
                result.error = fin.error;
            }
            call.write(new srpc_1.srpc.BucketWatchRequest(result));
        });
        call.on("error", (err) => {
            // console.error(err);
        });
        call.write(new srpc_1.srpc.BucketWatchRequest({
            name: this.name,
            key: key,
            plugin_id,
        }));
    }
    async _name() {
        return this.name;
    }
}
exports.Bucket = Bucket;
class Adapter {
    platform;
    bot_id;
    call;
    constructor(options) {
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
                    call.write(new srpc_1.srpc.AdapterRegistRequest({
                        bot_id: echo,
                        platform: v,
                    }));
                }
                if (__type__ == "action" && options.actionHandler) {
                    let v = await options.actionHandler(message);
                    call.write(new srpc_1.srpc.AdapterRegistRequest({
                        bot_id: echo,
                        platform: v,
                    }));
                }
                // console.log("end on data")
            });
            call.on("error", (err) => {
                console.error("adapter disc", err);
            });
            // console.log("before write")
            call.write(new srpc_1.srpc.AdapterRegistRequest({
                bot_id: options.bot_id,
                platform: options.platform,
            }));
            // console.log("after write write")
            this.call = call;
        }
    }
    setActionHandler(func) {
        // 将从服务端不断接收action消息，并处理
        // 事件处理巨饼
    }
    async receive(message) {
        //投递消息
        return new Promise((resolve, reject) => {
            client.AdapterReceive(new srpc_1.srpc.AdapterRequest({
                platform: this.platform,
                bot_id: this.bot_id,
                value: JSON.stringify(message),
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else if (resp?.value) {
                    resolve(new Sender(resp.value));
                }
            });
        });
    }
    async push(message) {
        //推送消息
        return new Promise((resolve, reject) => {
            client.AdapterPush(new srpc_1.srpc.AdapterRequest({
                platform: this.platform,
                bot_id: this.bot_id,
                value: JSON.stringify(message),
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resp?.value ?? "");
                }
            });
        });
    }
    async destroy() {
        this.call.cancel();
    }
    async sender(options) {
        return new Promise((resolve, reject) => {
            client.AdapterSender(new srpc_1.srpc.AdapterRequest({
                platform: this.platform,
                bot_id: this.bot_id,
                value: JSON.stringify(options),
            }), (err, resp) => {
                if (err) {
                    reject(err);
                }
                else if (resp?.value) {
                    resolve(new Sender(resp.value));
                }
            });
        });
    }
}
exports.Adapter = Adapter;
let sender = new Sender(process.env?.SENDER_ID ?? "");
exports.sender = sender;
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
exports.sleep = sleep;
class Console {
    error = (message, ...optionalParams) => { };
    info = (message, ...optionalParams) => { };
    log = (message, ...optionalParams) => { };
    debug = (message, ...optionalParams) => { };
}
let utils = {
    parseCQText: (text, prefix = "CQ") => {
        const cqRegex = new RegExp(`\\[${prefix}:(\\w+)(.*?)\\]`, "g");
        const cqMatches = text.matchAll(cqRegex);
        const result = [];
        let lastIndex = 0;
        for (const match of cqMatches) {
            // 添加 CQ 码前的文本
            const matchIndex = text.indexOf(match[0], lastIndex);
            if (matchIndex > lastIndex) {
                result.push(text.slice(lastIndex, matchIndex));
            }
            // 解析 CQ 码
            const params = {};
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
exports.utils = utils;
let slog = (type, ...args) => { };
let console = {
    log(...args) {
        const content = args.reduce((acc, arg) => acc + " " + arg, "");
        client.Console(new srpc_1.srpc.ConsoleRequest({ type: "log", content, plugin_id }), (err, resp) => { });
    },
    info(...args) {
        const content = args.reduce((acc, arg) => acc + " " + arg, "");
        client.Console(new srpc_1.srpc.ConsoleRequest({ type: "info", content, plugin_id }), (err, resp) => { });
    },
    error(...args) {
        const content = args.reduce((acc, arg) => acc + " " + arg, "");
        client.Console(new srpc_1.srpc.ConsoleRequest({ type: "error", content, plugin_id }), (err, resp) => { });
    },
    debug(...args) {
        const content = args.reduce((acc, arg) => acc + " " + arg, "");
        client.Console(new srpc_1.srpc.ConsoleRequest({ type: "debug", content, plugin_id }), (err, resp) => { });
    },
};
exports.console = console;
