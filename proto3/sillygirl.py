import asyncio
import json, pickle, base64, os, grpc
import time
import random
import threading
from queue import Queue

import srpc_pb2, srpc_pb2_grpc

plugin_id = os.environ.get("PLUGIN_ID", "")
runtime_id = os.environ.get("RUNTIME_ID", "")

metadata = (("RUNTIME_ID", runtime_id),)
metadata = grpc.metadata_call_credentials(metadata)

channel = grpc.insecure_channel("localhost:50051")
stub = srpc_pb2_grpc.SillyGirlServiceStub(channel)


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
        print("name=", name, "self.__name=", self.__name)

    def __getitem__(self, key):
        return self.get(key)

    def __getattr__(self, name):
        return self.get(name)

    def __setitem__(self, key, value):
        self.set(key, value)

    def __setattr__(self, name, value):
        if name in ["data", "_Bucket__name"]:
            # 使用object类的方法设置属性，避免无限递归
            object.__setattr__(self, name, value)
        else:
            self.set(name, value)

    def get(self, key, defaultValue=None):
        request = srpc_pb2.BucketKeyRequest(name=self.__name, key=key)
        response = stub.BucketGet(request)
        if response.value != "":
            return transform(response.value)
        else:
            return defaultValue

    def set(self, key, value):
        request = srpc_pb2.BucketSetRequest(
            name=self.__name, key=key, value=reverseTransform(value)
        )
        response = stub.BucketSet(request)
        return {"message": response.message, "changed": response.changed}

    def delete(self, key):
        return self.set(key, None)

    def deleteAll(self):
        request = srpc_pb2.BucketRequest(name=self.__name)
        stub.BucketDelete(request)

    def keys(self):
        print("name", self.__name)
        request = srpc_pb2.BucketKeyRequest(name=self.__name)
        response = stub.BucketKeys(request)
        return list(response.keys)

    def len(self):
        request = srpc_pb2.BucketRequest(name=self.__name)
        response = stub.BucketLen(request)
        return response.length

    def buckets(self):
        request = srpc_pb2.Empty()
        response = stub.BucketBuckets(request)
        return list(response.buckets)

    def watch(self, key, handle):
        queue = Queue()

        def entry_request_iterator():
            yield srpc_pb2.BucketWatchRequest(
                name=self.__name, key=key, plugin_id=plugin_id
            )
            while True:
                yield queue.get()
                print("queue.get()")
            # 阻塞A区 等待B区发送数据

        generator = entry_request_iterator()
        for response in stub.BucketWatch(generator):
            print("Server response:", response)
            old = transform(response.old)
            now = transform(response.now)
            result = handle(old, now, response.key)
            try:
                fin = result
            except Exception as e:
                print(e)
                continue
            storage_modifier = {
                "echo": response.echo,
            }
            if not fin:
                storage_modifier["error"] = "VOID"
            else:
                if "now" in fin:
                    storage_modifier["now"] = reverseTransform(fin["now"])
                if "message" in fin:
                    storage_modifier["message"] = fin["message"]
                if "error" in fin:
                    storage_modifier["error"] = fin["error"]
            queue.put(srpc_pb2.BucketWatchRequest(**storage_modifier))


app = Bucket("test")


def handle(old_value, new_value, key):
    print("Bucket value changed!")
    print("Key:", key)
    print("Old value:", old_value)
    print("New value:", new_value)
    return {
        "now": "shit",
    }


app.watch(key="*", handle=handle)
# app.name = 1
# print(app.name)
# app["name"] = (1,3, {3,4,5})
# print(app["name"])
# app.set("name", 1.23)
# print(app.get("name"))
# print("keys", app.keys())
# print("len", app.len())
# print("buckets", app.buckets())
# # app.delete("name")
# app.name = 0
# print(app.name)
# print("==")
# app.name = random.randrange(1,100,1)
# print("==")
# app.deleteAll()
print(app.name)


class Sender:
    pass


class Adapter:
    pass
