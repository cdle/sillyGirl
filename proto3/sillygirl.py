import asyncio
import json, pickle, base64, os, grpc

import srpc_pb2, srpc_pb2_grpc

plugin_id = os.environ.get("PLUGIN_ID", "")
runtime_id = os.environ.get("RUNTIME_ID", "")

metadata = [("runtime_id", runtime_id)]
channel = None
stub = None


def get_stub():  # 异步
    global stub, channel
    if stub == None:
        channel = grpc.insecure_channel("localhost:50051")
        stub = srpc_pb2_grpc.SillyGirlServiceStub(channel)
    return stub


channel2 = grpc.aio.insecure_channel("localhost:50051")
stub2 = srpc_pb2_grpc.SillyGirlServiceStub(channel2)


def transform(value):
    # 将从远程存储桶中获取的值进行转换，将其转换为相应的 Python 类型
    if value.startswith("f:"):
        return float(value[2:])
    elif value.startswith("i:"):
        return int(value[2:])
    elif value.startswith("b:"):
        return value[2:] == "true"
    elif value.startswith("o:"):
        return json.loads(value[2:])
    elif value.startswith("p:"):
        return pickle.loads(base64.b64decode(value[2:]))
    elif value == "":
        return None
    return value


def reverseTransform(value):
    # 将 Python 值转换为存储桶中的字符串表示形式
    try:
        if isinstance(value, float):
            return f"f:{value}"
        elif isinstance(value, bool):
            res = "true" if value else "false"
            return f"b:{res}"
        elif isinstance(value, int):
            return f"i:{value}"
        elif isinstance(value, str):
            return value
        elif value == None:
            return ""
        else:
            return f"o:{json.dumps(value)}"
    except:
        return "p:%s" % base64.b64encode(pickle.dumps(value)).decode("utf-8")


class Bucket:
    def __init__(self, name):
        self.__name = name

    def __getitem__(self, key):
        return self.__get(key)

    def __getattr__(self, name):
        return self.__get(name)

    def __setitem__(self, key, value):
        self.__set(key, value)
        # asyncio.ensure_future(self.set(key, value))

    def __setattr__(self, name, value):
        if name in ["data", "_Bucket__name"]:
            # 使用object类的方法设置属性，避免无限递归
            object.__setattr__(self, name, value)
        else:
            self.__set(name, value)
            # asyncio.ensure_future(self.set(name, value))

    async def get(self, key, defaultValue=None):
        request = srpc_pb2.BucketKeyRequest(name=self.__name, key=key)
        response = await stub2.BucketGet(request)
        if response.value != "":
            return transform(response.value)
        else:
            return defaultValue

    def __get(self, key, defaultValue=None):
        request = srpc_pb2.BucketKeyRequest(name=self.__name, key=key)
        response = get_stub().BucketGet(request)
        if response.value != "":
            return transform(response.value)
        else:
            return defaultValue

    def __set(self, key, value):
        try:
            request = srpc_pb2.BucketSetRequest(
                name=self.__name, key=key, value=reverseTransform(value)
            )
            response = get_stub().BucketSet(request)
            return {"message": response.message, "changed": response.changed}
        except Exception as e:
            return {"error": "fail to set value"}

    async def set(self, key, value):
        try:
            request = srpc_pb2.BucketSetRequest(
                name=self.__name, key=key, value=reverseTransform(value)
            )
            response = await stub2.BucketSet(request)
            return {"message": response.message, "changed": response.changed}
        except Exception as e:
            return {"error": "fail to set value"}

    async def delete(self, key):
        return await self.set(key, None)

    async def deleteAll(self):
        request = srpc_pb2.BucketRequest(name=self.__name)
        await stub2.BucketDelete(request)

    async def keys(self):
        print("name", self.__name)
        request = srpc_pb2.BucketKeyRequest(name=self.__name)
        response = await stub2.BucketKeys(request)
        return list(response.keys)

    async def len(self):
        request = srpc_pb2.BucketRequest(name=self.__name)
        response = await stub2.BucketLen(request)
        return response.length

    async def buckets(self):
        request = srpc_pb2.Empty()
        response = await stub2.BucketBuckets(request)
        return list(response.buckets)

    def watch(self, key, handle):
        async def __watch(self, key, handle):
            queue = asyncio.Queue()

            async def entry_request_iterator():
                yield srpc_pb2.BucketWatchRequest(
                    name=self.__name, key=key, plugin_id=plugin_id
                )
                while True:
                    yield await queue.get()

            generator = entry_request_iterator()
            response_iterator = stub2.BucketWatch(generator)
            async for response in response_iterator:
                old = transform(response.old)
                now = transform(response.now)
                try:
                    result = await handle(old, now, response.key)
                except Exception as e:
                    print(e)
                    continue
                storage_modifier = {
                    "echo": response.echo,
                }
                if not result:
                    storage_modifier["error"] = "VOID"
                else:
                    if "now" in result:
                        storage_modifier["now"] = reverseTransform(result["now"])
                    if "message" in result:
                        storage_modifier["message"] = result["message"]
                    if "error" in result:
                        storage_modifier["error"] = result["error"]
                await queue.put(srpc_pb2.BucketWatchRequest(**storage_modifier))

        asyncio.ensure_future(__watch(self, key, handle))


# test = Bucket("test")
# test.name = "777"
# print(test.name)

# async def handle(old_value, new_value, key):
#     print("Bucket value changed!")
#     print("Key:", key)
#     print("Old value:", old_value)
#     print("New value:", new_value)
#     return {
#         "now": "shit",
#         "message": "hehe",
#         "error": "xxx",
#     }


# async def main():
#     asyncio.ensure_future(test.watch(key="*", handle=handle))
#     await asyncio.sleep(0.1)
#     print("test.set", await test.set("name", random.randint(1, 10)))
#     print("test.name", await test.get("name"))
#     print("test.keys", await test.keys())
#     print("test.buckets", await test.buckets())
#     # print("test.delete", await test.delete("name"))
#     # print("test.deleteAll", await test.deleteAll())
#     print("test.name", await test.get("name"))

# loop = asyncio.get_event_loop()
# loop.run_until_complete(main())


class Sender:
    def __init__(self, uuid):
        self.__uuid = uuid
        self.destoried = False

    async def destroy(self):
        if self.destoried:
            return
        self.destoried = True
        await stub2.SenderDestroy(
            srpc_pb2.ReplyRequest(uuid=self.__uuid), metadata=metadata
        )

    async def getUserId(self):
        response = await stub2.SenderGetUserId(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getUserName(self):
        response = await stub2.SenderGetUserName(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getChatId(self):
        response = await stub2.SenderGetChatId(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getChatName(self):
        response = await stub2.SenderGetChatName(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getMessageId(self):
        response = await stub2.SenderGetMessageId(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getPlatform(self):
        response = await stub2.SenderGetPlatform(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getBotId(self):
        response = await stub2.SenderGetBotId(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def getContent(self):
        response = await stub2.SenderGetContent(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def isAdmin(self):
        response = await stub2.SenderIsAdmin(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )
        return response.value

    async def param(self, key):
        response = await stub2.SenderParam(
            request=srpc_pb2.ReplyRequest(uuid=self.__uuid, content=str(key)),
            metadata=metadata,
        )
        return response.value

    async def setContent(self, content):
        await stub2.SenderSetContent(
            request=srpc_pb2.SenderContentRequest(uuid=self.__uuid, content=content),
            metadata=metadata,
        )

    async def next(self):
        await stub2.SenderContinue(
            request=srpc_pb2.SenderRequest(uuid=self.__uuid), metadata=metadata
        )

    async def getAdapter(self):
        bot_id = await self.getBotId()
        platform = await self.getPlatform()
        return Adapter(bot_id, platform)

    async def listen(
        self,
        timeout=0,
        rules=[],
        handle=None,
        listen_private=False,
        listen_group=False,
        allow_platforms=[],
        prohibit_platforms=[],
        allow_groups=[],
        prohibit_groups=[],
        allow_users=[],
        prohibit_users=[],
    ):
        queue = asyncio.Queue()

        async def entry_request_iterator():
            yield srpc_pb2.SenderListenRequest(
                uuid=self.__uuid,
                timeout=timeout,
                rules=rules,
                listen_private=listen_private,
                listen_group=listen_group,
                allow_platforms=allow_platforms,
                prohibit_platforms=prohibit_platforms,
                allow_groups=allow_groups,
                prohibit_groups=prohibit_groups,
                allow_users=allow_users,
                prohibit_users=prohibit_users,
                persistent=self.__uuid == "",
                plugin_id=plugin_id,
            )
            while True:
                next_request = await queue.get()
                if next_request is None:
                    return
                yield next_request

        generator = entry_request_iterator()
        response_iterator = stub2.SenderListen(generator, metadata=metadata)
        s = None
        async for response in response_iterator:
            if response.echo == "END":
                break
            s = Sender(response.uuid) if response.uuid != "" else None
            if handle is not None and s is not None:
                try:
                    value = await handle(s)
                    await queue.put(
                        srpc_pb2.SenderListenRequest(
                            uuid=response.echo,
                            value=value,
                        )
                    )
                except Exception as e:
                    print(e)
            else:
                await queue.put(
                    srpc_pb2.SenderListenRequest(
                        uuid=response.echo,
                        value="",
                    )
                )

        await queue.put(None)
        return s

    def holdOn(self, string):
        return "go_again_" + string

    async def reply(self, content):
        return (
            await stub2.SenderReply(
                request=srpc_pb2.ReplyRequest(uuid=self.__uuid, content=content),
                metadata=metadata,
            )
        ).value

    async def doAction(self, properties):
        # Perform action based on the provided properties
        # ...
        pass


sender = Sender(os.environ.get("SENDER_ID", ""))


class Adapter:
    def __init__(self, platform, bot_id, replyHandler=None, actionHandler=None):
        self.platform = platform
        self.bot_id = bot_id

        if replyHandler is not None:
            asyncio.create_task(self.__run(replyHandler, actionHandler))

    async def __run(self, replyHandler, actionHandler):
        self.queue = asyncio.Queue()

        async def entry_request_iterator():
            yield srpc_pb2.AdapterRegistRequest(
                bot_id=self.bot_id, platform=self.platform
            )
            while True:
                next_request = await self.queue.get()
                if next_request is None:
                    return
                yield next_request

        generator = entry_request_iterator()

        response_iterator = stub2.AdapterRegist(generator, metadata=metadata)
        try:
            async for response in response_iterator:
                message = json.loads(response.value)
                echo = message["echo"]
                __type__ = message["__type__"]
                del message["__type__"]
                del message["echo"]
                if __type__ == "reply":
                    try:
                        v = replyHandler(message) or ""
                        await self.queue.put(
                            srpc_pb2.AdapterRegistRequest(bot_id=echo, platform=v)
                        )
                    except Exception as e:
                        print(e)

                if __type__ == "action" and actionHandler is not None:
                    try:
                        v = actionHandler(message) or ""
                        await self.queue.put(
                            srpc_pb2.AdapterRegistRequest(bot_id=echo, platform=v)
                        )
                    except Exception as e:
                        print(e)
        except:
            pass

    async def receive(self, message):
        # 投递消息
        return await stub2.AdapterReceive(
            srpc_pb2.AdapterRequest(
                platform=self.platform, bot_id=self.bot_id, value=json.dumps(message)
            ),
            metadata=metadata,
        )

    async def push(self, message):
        # 推送消息
        response = await stub2.AdapterPush(
            srpc_pb2.AdapterRequest(
                platform=self.platform, bot_id=self.bot_id, value=json.dumps(message)
            ),
            metadata=metadata,
        )
        return response.value or ""

    async def destroy(self):
        if not hasattr(self, "self.__destroyed"):
            await self.queue.put(None)
            self.__destroyed = True

    async def sender(self, options={}):
        response = await stub2.AdapterSender(
            srpc_pb2.AdapterRequest(
                platform=self.platform, bot_id=self.bot_id, value=json.dumps(options)
            ),
            metadata=metadata,
        )
        if response.value:
            return Sender(response.value)


import re


class Utils:
    def buildCQTag(self, type, params, prefix="CQ"):
        param_strings = []
        for key, value in params.items():
            param_string = f"{key}={value}"
            param_strings.append(param_string)
        param_string = ",".join(param_strings)
        cq_string = f"[{prefix}:{type}{',' + param_string if param_string else ''}]"
        return cq_string

    def parseCQText(self, text, prefix="CQ"):
        cq_regex = re.compile(rf"\[{prefix}:(\w+)(.*?)\]", re.DOTALL)
        cq_matches = cq_regex.finditer(text)
        result = []

        last_index = 0
        for match in cq_matches:
            # 添加 CQ 码前的文本
            match_index = text.index(match.group(0), last_index)
            if match_index > last_index:
                result.append(text[last_index:match_index])

            # 解析 CQ 码
            params = {}
            param_regex = re.compile(r"(\w+)=([^,]+)")
            param_matches = param_regex.findall(match.group(2))
            for param_match in param_matches:
                params[param_match[0]] = param_match[1].strip()
            result.append(
                {
                    "type": match.group(1),
                    "params": params,
                }
            )

            last_index = match_index + len(match.group(0))

        if last_index < len(text):
            result.append(text[last_index:])

        return result

    def image(self, url):
        return self.buildCQTag("image", {"url": url})

    def video(self, url):
        return self.buildCQTag("video", {"url": url})


utils = Utils()

import srpc_pb2
import srpc_pb2_grpc


class Console:
    def __init__(self, plugin_id):
        self.plugin_id = plugin_id

    def log(self, *args):
        self.send_console_request("log", *args)

    def info(self, *args):
        self.send_console_request("info", *args)

    def error(self, *args):
        self.send_console_request("error", *args)

    def debug(self, *args):
        self.send_console_request("debug", *args)

    def send_console_request(self, console_type, *args):
        content = " ".join(map(str, args))
        request = srpc_pb2.ConsoleRequest(
            type=console_type, content=content, plugin_id=self.plugin_id
        )
        stub2.Console(request)


console = Console(plugin_id)
