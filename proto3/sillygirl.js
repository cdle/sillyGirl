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
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var _a, _b;
Object.defineProperty(exports, "__esModule", { value: true });
exports.sleep = exports.sender = exports.Bucket = exports.Adapter = void 0;
const srpc_1 = require("./srpc");
const grpc_1 = __importStar(require("@grpc/grpc-js"));
let client = new srpc_1.srpc.SillyGirlServiceClient("localhost:50051", grpc_1.credentials.createInsecure());
class Sender {
    constructor(uuid) {
        this.uuid = uuid;
    }
    getUserId() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetUserId(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getUserName() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetUserName(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getChatId() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetChatId(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getChatName() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetChatName(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getMessageId() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetMessageId(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getPlatform() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetPlatform(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getBotId() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetBotId(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    getContent() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderGetContent(new srpc_1.srpc.SenderRequest({
                    uuid: this.uuid,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    setContent(content) {
        return __awaiter(this, void 0, void 0, function* () {
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
        });
    }
    continue() {
        return __awaiter(this, void 0, void 0, function* () {
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
        });
    }
    getAdapter() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Adapter({
                bot_id: yield this.getBotId(),
                platform: yield this.getPlatform(),
            });
        });
    }
    listen(options) {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                let params = {
                    uuid: this.uuid,
                    rules: options === null || options === void 0 ? void 0 : options.rules,
                    timeout: options === null || options === void 0 ? void 0 : options.timeout,
                    listen_private: options === null || options === void 0 ? void 0 : options.listen_private,
                    listen_group: options === null || options === void 0 ? void 0 : options.listen_group,
                    allow_platforms: options === null || options === void 0 ? void 0 : options.allow_platforms,
                    prohibit_platforms: options === null || options === void 0 ? void 0 : options.prohibit_platforms,
                    allow_groups: options === null || options === void 0 ? void 0 : options.allow_groups,
                    prohibit_groups: options === null || options === void 0 ? void 0 : options.prohibit_groups,
                    allow_users: options === null || options === void 0 ? void 0 : options.allow_users,
                    prohibit_users: options === null || options === void 0 ? void 0 : options.prohibit_users,
                };
                client.SenderListen(new srpc_1.srpc.SenderListenRequest(params), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        if (resp === null || resp === void 0 ? void 0 : resp.value) {
                            resolve(new Sender(resp.value));
                        }
                        else {
                            reject(new Error("timeout"));
                        }
                    }
                });
            });
        });
    }
    reply(content) {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.SenderReply(new srpc_1.srpc.ReplyRequest({
                    uuid: this.uuid,
                    content,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.value);
                    }
                });
            });
        });
    }
    action(options) {
        return __awaiter(this, void 0, void 0, function* () {
            return;
        });
    }
    event() {
        return __awaiter(this, void 0, void 0, function* () {
            return;
        });
    }
}
class Bucket {
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
    get(key, defaultValue = undefined) {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.BucketGet(new srpc_1.srpc.BucketKeyRequest({ name: this.name, key }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(this.transform(resp === null || resp === void 0 ? void 0 : resp.value) || defaultValue);
                    }
                });
            });
        });
    }
    set(key, value) {
        return __awaiter(this, void 0, void 0, function* () {
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
                            message: resp === null || resp === void 0 ? void 0 : resp.message,
                            changed: resp === null || resp === void 0 ? void 0 : resp.changed,
                        });
                    }
                });
            });
        });
    }
    getAll() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.BucketGetAll(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        let values = {};
                        if (resp === null || resp === void 0 ? void 0 : resp.value) {
                            values = JSON.parse(resp === null || resp === void 0 ? void 0 : resp.value);
                            for (let key in values) {
                                values[key] = this.transform(values[key]);
                            }
                        }
                        resolve(values);
                    }
                });
            });
        });
    }
    delete() {
        return __awaiter(this, void 0, void 0, function* () {
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
        });
    }
    keys() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.BucketKeys(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.keys);
                    }
                });
            });
        });
    }
    len() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.BucketLen(new srpc_1.srpc.BucketRequest({ name: this.name }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.length);
                    }
                });
            });
        });
    }
    buckets() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.BucketBuckets(new srpc_1.srpc.Empty(), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        resolve(resp === null || resp === void 0 ? void 0 : resp.buckets);
                    }
                });
            });
        });
    }
    _name() {
        return __awaiter(this, void 0, void 0, function* () {
            return this.name;
        });
    }
}
exports.Bucket = Bucket;
class Adapter {
    constructor(options) {
        this.platform = options.platform;
        this.bot_id = options.bot_id;
        const call = client.AdapterRegist();
        if (options.replyHandler) {
            let callback = options.replyHandler;
            call.on("data", (response) => {
                let message = JSON.parse(response.value);
                const { echo, __type__ } = message;
                delete (message.__type__);
                delete (message.echo);
                if (__type__ == "reply") {
                    call.write(new srpc_1.srpc.AdapterRegistRequest({
                        bot_id: echo,
                        platform: callback(message),
                    }));
                }
            });
            call.on("error", (err) => {
                // console.error(err);
            });
            call.write(new srpc_1.srpc.AdapterRegistRequest({
                bot_id: options.bot_id,
                platform: options.platform,
            }));
            this.call = call;
        }
    }
    setActionHandler(func) {
        // 将从服务端不断接收action消息，并处理
        // 事件处理巨饼
    }
    receive(message) {
        return __awaiter(this, void 0, void 0, function* () {
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
                    else if (resp === null || resp === void 0 ? void 0 : resp.value) {
                        resolve(new Sender(resp.value));
                    }
                });
            });
        });
    }
    push(message) {
        return __awaiter(this, void 0, void 0, function* () {
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
                    else if (resp === null || resp === void 0 ? void 0 : resp.value) {
                        resolve(resp.value);
                    }
                });
            });
        });
    }
    destroy() {
        return __awaiter(this, void 0, void 0, function* () {
            this.call.cancel();
        });
    }
    sender(options) {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                client.AdapterSender(new srpc_1.srpc.AdapterRequest({
                    platform: this.platform,
                    bot_id: this.bot_id,
                }), (err, resp) => {
                    if (err) {
                        reject(err);
                    }
                    else if (resp === null || resp === void 0 ? void 0 : resp.value) {
                        resolve(new Sender(resp.value));
                    }
                });
            });
        });
    }
}
exports.Adapter = Adapter;
let sender = new Sender((_b = (_a = process.env) === null || _a === void 0 ? void 0 : _a.SENDER_ID) !== null && _b !== void 0 ? _b : "4d6371a8-2778-11ee-a3c2-821680fbbf6b");
exports.sender = sender;
function sleep(ms) {
    return __awaiter(this, void 0, void 0, function* () {
        return new Promise((resolve) => setTimeout(resolve, ms));
    });
}
exports.sleep = sleep;
