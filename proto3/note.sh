protoc-gen-grpc --js_out=import_style=commonjs:. --grpc_out=. srpc.proto
protoc-gen-grpc --ts_out=service=grpc-node:. --grpc_out=. srpc.proto

protoc --go_out=. -I. --go-grpc_out=.  srpc.proto

protoc --plugin=protoc-gen-ts=$(which protoc-gen-ts) --js_out=import_style=commonjs,binary:./ --ts_out=./ srpc.proto

protoc --js_out=import_style=commonjs,binary:. --grpc_out=.  srpc.proto


#ok
protoc --go_out=. -I. --go-grpc_out=.  srpc.proto
protoc-gen-grpc --ts_out=service=grpc-node:.  srpc.proto

#protoc-gen-grpc --python_out=.  srpc.proto
#protoc --python_out=.  srpc.proto
# pip install "grpcio-tools==1.43.0"
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. srpc.proto


#打包
npx webpack --config webpack.config.js

#linux编译：
scp /home/user/Code/sillyGirl/proto3/dist/sillygirl.js root@example.com:/root/node/node-18.16.1/lib
ssh root@example.com
cd /root/node/node-18.16.1 && ninja -C out/Release && scp -P 20211 out/Release/node a1-6@example.com:/home/user/Code/nodes/node_linux_amd64
#
macos编译：
cp /home/user/Code/sillyGirl/proto3/dist/sillygirl.js /home/user/Code/node/node-v18.16.1/lib/sillygirl.js && cd /home/user/Code/node/node-v18.16.1 && ninja -C out/Release && cp out/Release/node /home/user/Code/nodes/node_darwin_arm64

#压缩
cd /home/user/Code/nodes/node_darwin_arm64
cd /home/user/Code/nodes/node_linux_amd64

##
git add . && git commit -m "x" && git push