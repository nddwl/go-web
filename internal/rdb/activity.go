package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-web/internal/model"
	"go-web/utils/ecode"
	"strconv"
)

type Activity struct {
	*Rdb
}

func NewActivity(rdb *Rdb) *Activity {
	return &Activity{rdb}
}

func (t *Activity) List(activityUUID int64, prize []model.Prize) (err error) {
	ctx := context.Background()
	uuid := strconv.FormatInt(activityUUID, 10)
	err = t.db.Watch(ctx, func(tx *redis.Tx) error {
		_, err1 := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			members := make([]redis.Z, len(prize))
			values := make(map[string]interface{}, len(prize))
			for k, v := range prize {
				members[k] = redis.Z{
					Score:  float64(v.Score),
					Member: v.UUID,
				}
				values[strconv.FormatInt(v.UUID, 10)] = v.Stock
			}
			pipe.Del(ctx, "activity_stock:"+uuid)
			pipe.HSet(ctx, "activity_stock:"+uuid, values)
			pipe.Del(ctx, "activity:"+uuid)
			pipe.ZAdd(ctx, "activity:"+uuid, members...)
			return nil
		})
		if err1 != nil {
			return err1
		}
		return nil
	}, "activity_stock:"+uuid, "activity:"+uuid)
	return
}

func (t *Activity) GetPrizeStock(activityUUID int64) (map[string]string, error) {
	ctx := context.Background()
	const luaScript = `return redis.call("HGETALL", KEYS[1])`
	result, err := t.db.Eval(ctx, luaScript, []string{"activity_stock:" + strconv.FormatInt(activityUUID, 10)}).Result()
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

func (t *Activity) UnList(activityUUID int64) error {
	uuid := strconv.FormatInt(activityUUID, 10)
	const luaScript = `return redis.call("del", KEYS[1], KEYS[2])`
	return t.db.Eval(context.Background(), luaScript, []string{
		"activity:" + uuid,
		"activity_stock:" + uuid,
	}).Err()
}

func (t *Activity) Lottery(activityUUID int64) (prizeUUID int64, err error) {

	const luaScript = `
    local prizeUUIDs = redis.call("ZRANDMEMBER", "activity:" .. ARGV[1], 1)
    if #prizeUUIDs == 0 then
        return 0
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
	prizeUUID, _ = strconv.ParseInt(result.(string), 10, 64)
	if prizeUUID == 0 {
		err = ecode.ActivityIsOver
	}
	return
}

func (t *Activity) IsList(activityUUID int64) (list bool, err error) {
	i, err := t.db.Exists(context.Background(), "activity:"+strconv.FormatInt(activityUUID, 10)).Result()
	if err != nil {
		return
	}
	if i == 1 {
		list = true
	}
	return
}
