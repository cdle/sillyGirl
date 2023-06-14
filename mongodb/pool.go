package mongodb

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB连接池
type MongoDBPool struct {
	clients map[string]*mongo.Client // 客户端连接map
	mutex   sync.Mutex               // 互斥锁
}

var p *MongoDBPool

// 创建一个新的MongoDB连接池
func NewMongoDBPool() *MongoDBPool {
	return &MongoDBPool{
		clients: make(map[string]*mongo.Client),
	}
}

// 获取一个MongoDB客户端连接
func GetClient(uri string) (*mongo.Client, error) {
	if p == nil {
		p = NewMongoDBPool()
	}
	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 检查连接是否存在
	if client, ok := p.clients[uri]; ok {
		return client, nil
	}

	// 创建新的连接
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// 保存连接
	p.clients[uri] = client

	return client, nil
}

// 关闭一个MongoDB客户端连接
func CloseClient(uri string) error {
	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 检查连接是否存在
	if client, ok := p.clients[uri]; ok {
		// 关闭连接
		err := client.Disconnect(context.Background())
		if err != nil {
			return err
		}

		// 删除连接
		delete(p.clients, uri)

		return nil
	}

	return errors.New("client not found")
}
