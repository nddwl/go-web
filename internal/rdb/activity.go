package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-web/internal/model"
	"go-web/utils/ecode"
)

type Activity struct {
	*Rdb
}

func NewActivity(rdb *Rdb) *Activity {
	return &Activity{rdb}
}

func (t *Activity) List(activityUUID string, prize []model.Prize) (result []string, err error) {
	ctx := context.Background()
	err = t.db.Watch(ctx, func(tx *redis.Tx) error {
		cmder, err1 := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			members := make([]redis.Z, len(prize))
			values := make(map[string]interface{}, len(prize))
			for k, v := range prize {
				members[k] = redis.Z{
					Score:  float64(v.Score),
					Member: v.UUID,
				}
				values[v.UUID] = v.Stock
			}
			pipe.Del(ctx, "activity_stock:"+activityUUID)
			pipe.HSet(ctx, "activity_stock:"+activityUUID, values)
			pipe.Del(ctx, "activity:"+activityUUID)
			pipe.ZAddNX(ctx, "activity:"+activityUUID, members...)
			return nil
		})
		if err1 != nil {
			return err1
		}
		result = make([]string, len(cmder))
		for i := 0; i < len(result); i++ {
			result[i] = cmder[i].String()
		}
		return nil
	}, "activity_stock:"+activityUUID, "activity:"+activityUUID)
	return
}

func (t *Activity) GetPrizeStock(activityUUID string) (map[string]string, error) {
	ctx := context.Background()
	const luaScript = `return redis.call("HGETALL", KEYS[1])`
	result, err := t.db.Eval(ctx, luaScript, []string{"activity_stock:" + activityUUID}).Result()
	if err != nil {
		return nil, err
	}
	if resultArray, ok := result.([]interface{}); ok {
		prizeStock := make(map[string]string)
		for i := 0; i < len(resultArray); i += 2 {
			key := resultArray[i].(string)
			value := resultArray[i+1].(string)
			prizeStock[key] = value
		}
		return prizeStock, nil
	}
	return nil, fmt.Errorf("unexpected result type")
}

func (t *Activity) UnList(activityUUID string) error {
	const luaScript = `return redis.call("del", KEYS[1], KEYS[2])`
	return t.db.Eval(context.Background(), luaScript, []string{
		"activity:" + activityUUID,
		"activity_stock:" + activityUUID,
	}).Err()
}

func (t *Activity) Lottery(activityUUID string) (prizeUUID string, err error) {

	const luaScript = `
    local prizeUUIDs = redis.call("ZRANDMEMBER", "activity:" .. ARGV[1], 1)
    if #prizeUUIDs == 0 then
        return ""
    end

    local prizeUUID = prizeUUIDs[1]
    local stockKey = "activity_stock:" .. ARGV[1]
    local stock = redis.call("HINCRBY", stockKey, prizeUUID, -1)

    if stock <= 0 then
        redis.call("ZREM", "activity:" .. ARGV[1], prizeUUID)
    end

    return prizeUUID
    `

	ctx := context.Background()
	result, err := t.db.Eval(ctx, luaScript, []string{}, activityUUID).Result()
	if err != nil {
		return
	}
	prizeUUID = result.(string)
	if prizeUUID == "" {
		err = ecode.ActivityIsOver
	}
	return
}
