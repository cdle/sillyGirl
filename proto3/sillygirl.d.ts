declare class Sender {
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
declare let utils: {
    parseCQText: (text: string, prefix?: string) => (string | {
        type: string;
        params: any;
    })[];
};
declare let console: {
    log(...args: any[]): void;
    info(...args: any[]): void;
    error(...args: any[]): void;
    debug(...args: any[]): void;
};
export { Adapter, Bucket, sender, sleep, utils, console };
