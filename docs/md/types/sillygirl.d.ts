declare class Sender {
    private uuid;
    private destoried;
    constructor(uuid: string);
    destroy(): void;
    getUserId(): Promise<string>;
    getUserName(): Promise<string>;
    getChatId(): Promise<string>;
    getChatName(): Promise<string>;
    getMessageId(): Promise<string>;
    getPlatform(): Promise<string>;
    getBotId(): Promise<string>;
    getContent(): Promise<string>;
    isAdmin(): Promise<boolean>;
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
    }): Promise<Sender | undefined>;
    holdOn(str: string): string;
    reply(content: string): Promise<string>;
    doAction(options: Record<string, any>): Promise<any>;
    getEvent(): Promise<Record<string, any>>;
}
declare class Bucket {
    private name;
    constructor(name: string);
    transform(v: string | undefined): string | number | boolean | undefined;
    reverseTransform(value: any): string;
    get(key: string, defaultValue?: any): Promise<any>;
    set(key: string, value: any): Promise<{
        message?: string;
        changed?: boolean;
    }>;
    getAll(): Promise<Record<string, any>>;
    delete(key: string): Promise<{
        message?: string;
        changed?: boolean;
    }>;
    deleteAll(): Promise<undefined>;
    keys(): Promise<string[]>;
    len(): Promise<number>;
    buckets(): Promise<string[]>;
    watch(key: string, handle: (old: any, now: any, key: string) => StorageModifier | void): void;
    getName(): Promise<string>;
}
interface StorageModifier {
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
    platform: string;
    bot_id: string;
    call: any;
    constructor(options: {
        platform: string;
        bot_id: string;
        replyHandler?: (message: Message) => Promise<string | undefined>;
        actionHandler?: (message: Message) => Promise<string | undefined>;
    });
    receive(message: Message): Promise<undefined>;
    push(message: Message): Promise<string>;
    destroy(): Promise<void>;
    sender(options: any): Promise<Sender>;
}
declare let sender: Sender;
declare function sleep(ms?: number): Promise<unknown>;
interface CQItem {
    type: string;
    params: Record<string, string>;
}
interface CQParams {
    [key: string]: string | number | boolean;
}
declare let utils: {
    buildCQTag: (type: string, params: CQParams, prefix?: string) => string;
    parseCQText: (text: string, prefix?: string) => (string | CQItem)[];
};
declare let console: {
    log(...args: any[]): void;
    info(...args: any[]): void;
    error(...args: any[]): void;
    debug(...args: any[]): void;
};

export { Adapter, Bucket, sender, sleep, utils, console };
