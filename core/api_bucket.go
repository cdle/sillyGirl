package core

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func init() {
	var sillyGirl = MakeBucket("sillyGirl")
	GinApi(GET, "/api/storage/list", RequireAuth, func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
		perPage, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
		keys := ctx.Query("keys")
		data := []map[string]string{}
		arr := strings.Split(keys, ",")
		if keys == "" {
			ctx.JSON(200, map[string]interface{}{
				"success": true,
				"data":    data,
				"page":    page,
				"total":   len(data),
			})
			return
		}
		for _, bk := range arr {
			ar := strings.Split(bk, ".")
			if len(ar) == 2 {
				if ar[0] == "plugins" && false { //todo
					// data[bk] = halfDeEct(MakeBucket(ar[0]).GetString(ar[1]))
				} else {
					// data[bk] = MakeBucket(ar[0]).GetString(ar[1])
					data = append(data, map[string]string{
						"bucket": ar[0],
						"key":    ar[1],
						"value":  MakeBucket(ar[0]).GetString(ar[1]),
					})
				}
			}
			if len(ar) == 1 {
				MakeBucket(ar[0]).Foreach(func(b1, b2 []byte) error {
					data = append(data, map[string]string{
						"bucket": bk,
						"key":    string(b1),
						"value":  string(b2),
					})
					return nil
				})
			}
		}
		start := (page - 1) * perPage
		end := start + perPage
		if end > len(data) {
			end = len(data)
		}
		res := data[start:end]
		index := start + 1
		for i := range res {
			res[i]["index"] = fmt.Sprint(index)
			index++
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data":    res,
			"page":    page,
			"total":   len(data),
		})
	})
	GinApi(GET, "/api/storage", RequireAuth, func(ctx *gin.Context) {
		keys := ctx.Query("keys")
		if keys == "" {
			buckets := sillyGirl.Buckets()
			search := ctx.Query("search")
			res := []map[string]interface{}{}
			if search == "" {
				for _, bucket := range buckets {
					if bucket == "plugins" {
						continue
					}
					res = append(res, map[string]interface{}{
						"value": bucket,
						"text":  "[桶] " + bucket,
					})
				}
				ctx.JSON(200, map[string]interface{}{
					"success": true,
					"data":    res,
				})
				return
			}
			for _, bucket := range buckets {
				if bucket == "plugins" {
					continue
				}
				if strings.Contains(bucket, search) {
					res = append(res, map[string]interface{}{
						"value": bucket,
						"text":  "[桶] " + bucket,
					})
				}
				b := MakeBucket(bucket)
				b.Foreach(func(b1, b2 []byte) error {
					key := string(b1)
					value := string(b2)
					if strings.Contains(key, search) {
						res = append(res, map[string]interface{}{
							"value": bucket + "." + key,
							"text":  "[键] " + key,
						})
					}
					if strings.Contains(value, search) {
						res = append(res, map[string]interface{}{
							"value": bucket + "." + key,
							"text":  "[值] " + value,
						})
					}
					return nil
				})
			}

			ctx.JSON(200, map[string]interface{}{
				"success": true,
				"data":    res,
			})
			return
		}
		data := map[string]interface{}{}
		arr := strings.Split(keys, ",")
		for _, bk := range arr {
			ar := strings.Split(bk, ".")
			if len(ar) == 2 {
				if ar[0] == "plugins" && false { //todo
					data[bk] = halfDeEct(MakeBucket(ar[0]).GetString(ar[1]))
				} else {
					data[bk] = TransformBucketKeyValue(MakeBucket(ar[0]).GetString(ar[1]))
				}
			}
			if len(ar) == 1 {
				MakeBucket(ar[0]).Foreach(func(b1, b2 []byte) error {
					data[bk+"."+string(b1)] = TransformBucketKeyValue(string(b2))
					return nil
				})
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data":    data,
		})
	})
	GinApi(PUT, "/api/storage", RequireAuth, func(ctx *gin.Context) {
		data, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
				"showType":     2,
			})
			return
		}
		updates := map[string]interface{}{}
		err = json.Unmarshal(data, &updates)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
				"showType":     2,
			})
			return
		}
		messages := map[string]interface{}{}
		errors := map[string]interface{}{}
		for bk, v := range updates {
			ar := strings.Split(bk, ".")
			if len(ar) == 2 {
				msg, err := SetBucketKeyValue(MakeBucket(ar[0]), ar[1], v)
				if msg != "" {
					messages[bk] = msg
				}
				if err != nil {
					errors[bk] = err.Error()
				}
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success":  true,
			"messages": messages,
			"errors":   errors,
		})
	})
}
