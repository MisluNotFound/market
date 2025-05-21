package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var client *elasticsearch.Client

func Init() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	cli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	res, err := cli.Info()
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	client = cli

	if err := InitIndex(); err != nil {
		panic(err)
	}
}

func InitIndex() error {
    mapping := map[string]interface{}{
        "settings": map[string]interface{}{
            "number_of_shards":   1,
            "number_of_replicas": 0,
            "analysis": map[string]interface{}{
                "analyzer": map[string]interface{}{
                    "ik_max_word_analyzer": map[string]interface{}{
                        "type":      "custom",
                        "tokenizer": "ik_max_word",
                    },
                },
            },
        },
        // TODO 添加商品的所有信息，查询仅依赖于es
        "mappings": map[string]interface{}{
            "properties": map[string]interface{}{
                "describe": map[string]interface{}{
                    "type":     "text",
                    "analyzer": "ik_max_word",
                    "fields": map[string]interface{}{
                        "keyword": map[string]interface{}{
                            "type": "keyword",
                        },
                    },
                },
                "id": map[string]interface{}{
                    "type": "keyword",
                },
                "category": map[string]interface{}{
                    "type": "keyword",
                },
                "created_at": map[string]interface{}{
                    "type":   "date",
                    "format": "strict_date_optional_time||epoch_millis",
                },
                "attributes": map[string]interface{}{
                    "type": "nested",
                    "properties": map[string]interface{}{
                        "key": map[string]interface{}{
                            "type": "keyword",
                        },
                        "value": map[string]interface{}{
                            "type": "keyword",
                        },
                    },
                },
            },
        },
    }

    index := "m-market"
    ctx := context.Background()

    // 检查索引是否存在
    existsReq := esapi.IndicesExistsRequest{Index: []string{index}}
    existsRes, err := existsReq.Do(ctx, client)
    if err != nil {
        return fmt.Errorf("检查索引失败: %w", err)
    }
    defer existsRes.Body.Close()

    // 如果索引不存在，直接创建
    if existsRes.StatusCode != 200 {
        body, err := json.Marshal(mapping)
        if err != nil {
            return fmt.Errorf("序列化映射失败: %v", err)
        }
        req := esapi.IndicesCreateRequest{
            Index: index,
            Body:  bytes.NewReader(body),
        }
        res, err := req.Do(ctx, client)
        if err != nil {
            return fmt.Errorf("创建索引失败: %w", err)
        }
        defer res.Body.Close()
        if res.IsError() {
            return fmt.Errorf("创建索引错误: %s", res.String())
        }
        log.Printf("✅ Index [%s] created", index)
        return nil
    }

    // 索引存在，获取当前映射
    getMappingReq := esapi.IndicesGetMappingRequest{Index: []string{index}}
    getMappingRes, err := getMappingReq.Do(ctx, client)
    if err != nil {
        return fmt.Errorf("获取映射失败: %w", err)
    }
    defer getMappingRes.Body.Close()

    // 解析当前映射
    var currentMapping map[string]interface{}
    if err := json.NewDecoder(getMappingRes.Body).Decode(&currentMapping); err != nil {
        return fmt.Errorf("解析当前映射失败: %v", err)
    }

    // 提取索引的 mappings 部分
    currentIndexMapping, ok := currentMapping[index].(map[string]interface{})
    if !ok {
        return fmt.Errorf("无法解析索引 %s 的映射", index)
    }
    currentMappings, ok := currentIndexMapping["mappings"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("无法解析 mappings 部分")
    }

    // 序列化当前映射和目标映射以比较
    currentJSON, err := json.Marshal(currentMappings)
    if err != nil {
        return fmt.Errorf("序列化当前映射失败: %v", err)
    }
    targetJSON, err := json.Marshal(mapping["mappings"])
    if err != nil {
        return fmt.Errorf("序列化目标映射失败: %v", err)
    }

    // 比较映射
    if string(currentJSON) == string(targetJSON) {
        log.Printf("🔁 Index [%s] mappings unchanged", index)
        return nil
    }

    // 映射不同，尝试更新
    log.Printf("🔄 Index [%s] mappings changed, updating...", index)
    updateBody, err := json.Marshal(mapping["mappings"])
    if err != nil {
        return fmt.Errorf("序列化更新映射失败: %v", err)
    }
    updateReq := esapi.IndicesPutMappingRequest{
        Index: []string{index},
        Body:  bytes.NewReader(updateBody),
    }
    updateRes, err := updateReq.Do(ctx, client)
    if err != nil {
        return fmt.Errorf("更新映射失败: %w", err)
    }
    defer updateRes.Body.Close()

    if updateRes.IsError() {
        // 如果更新失败（例如字段类型冲突），需考虑重建索引
        return fmt.Errorf("更新映射错误: %s", updateRes.String())
    }

    log.Printf("✅ Index [%s] mappings updated", index)
    return nil
}