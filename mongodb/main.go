package mongodb

import (
	"context"

	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	Uri    string
	Vm     *goja.Runtime
}

func (mc *MongoClient) Connect(uri string) {
	var err error
	if mc.Uri == "" {
		mc.Uri = uri
	}
	mc.client, err = GetClient(mc.Uri)
	if err != nil {
		panic(mc.Vm.NewGoError(err))
	}
}

func (mc *MongoClient) Db(dbname string) interface{} {
	if mc.client == nil {
		mc.Connect("")
	}
	return &MongoDB{
		Vm:       mc.Vm,
		database: mc.client.Database(dbname),
	}
}

type MongoDB struct {
	database *mongo.Database
	Vm       *goja.Runtime
}

func (db *MongoDB) Collection(clname string) interface{} {
	return &Collection{
		Vm:         db.Vm,
		collection: db.database.Collection(clname),
	}
}

type Collection struct {
	collection *mongo.Collection
	Vm         *goja.Runtime
}

func (cl *Collection) InsertOne(insert map[string]interface{}) interface{} {
	result, err := cl.collection.InsertOne(context.Background(), insert)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) FindOne(filter interface{}) interface{} {
	var result = map[string]interface{}{}
	err := cl.collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) FindOneAndDelete(filter interface{}) interface{} {
	var result = map[string]interface{}{}
	err := cl.collection.FindOneAndDelete(context.Background(), filter).Decode(result)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) FindOneAndReplace(filter interface{}, replace interface{}) interface{} {
	var result = map[string]interface{}{}
	err := cl.collection.FindOneAndReplace(context.Background(), filter, replace).Decode(result)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) FindOneAndUpdate(filter interface{}, update interface{}) interface{} {
	var result = map[string]interface{}{}
	err := cl.collection.FindOneAndUpdate(context.Background(), filter, update).Decode(result)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) Find(filter interface{}, params map[string]interface{}) interface{} {
	opts := []*options.FindOptions{}
	count := false
	if v, ok := params["sort"]; ok {
		opt := options.Find()
		opt.SetSort(v)
		opts = append(opts, opt)
	}
	if _, ok := params["count"]; ok {
		count = true
	} else {
		if v, ok := params["skip"]; ok {
			opt := options.Find()
			opt.SetSkip(utils.Int64(v))
			opts = append(opts, opt)
		}
		if v, ok := params["limit"]; ok {
			opt := options.Find()
			opt.SetLimit(utils.Int64(v))
			opts = append(opts, opt)
		}
	}
	// opts.SetSkip((pageNumber - 1) * pageSize)
	// opts.SetLimit(pageSize)
	cur, err := cl.collection.Find(context.Background(), filter, opts...)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	defer cur.Close(context.Background())
	var results = []map[string]interface{}{}
	for cur.Next(context.Background()) {
		var person = map[string]interface{}{}
		err := cur.Decode(&person)
		if err != nil {
			panic(cl.Vm.NewGoError(err))
		}
		results = append(results, person)
	}
	if err := cur.Err(); err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	if count {
		v, _ := cl.collection.CountDocuments(context.Background(), filter)
		return map[string]interface{}{
			"count":   v,
			"results": results,
		}
	}
	return results
}

func (cl *Collection) UpdateOne(filter, update interface{}, params map[string]interface{}) interface{} {
	opts := []*options.UpdateOptions{}
	if v, ok := params["upsert"]; ok {
		opts = append(opts, options.Update().SetUpsert(v.(bool)))
	}
	result, err := cl.collection.UpdateOne(context.Background(), filter, update, opts...)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) UpdateMany(filter, update interface{}) interface{} {

	result, err := cl.collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) DeleteOne(filter interface{}) interface{} {
	result, err := cl.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}

func (cl *Collection) DeleteMany(filter interface{}) interface{} {
	result, err := cl.collection.DeleteMany(context.Background(), filter)
	if err != nil {
		panic(cl.Vm.NewGoError(err))
	}
	return result
}
