package common

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

// IDGenerator ID生成器接口
type IDGenerator interface {
	Generate() string
	IsValid(id string) bool
}

// UUIDGenerator UUID生成器（PostgreSQL）
type UUIDGenerator struct{}

func (g *UUIDGenerator) Generate() string {
	return uuid.New().String()
}

func (g *UUIDGenerator) IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// ObjectIdGenerator ObjectId对象ID生成器（MongoDB）
type ObjectIdGenerator struct{}

func (g *ObjectIdGenerator) Generate() string {
	return bson.NewObjectId().Hex()
}

func (g *ObjectIdGenerator) IsValid(id string) bool {
	return bson.IsObjectIdHex(id)
}

// IDConverter ID转换器，处理跨数据库的ID转换
type IDConverter struct {
	ObjectIdToUUIDMap map[string]string
	UUIDToObjectIdMap map[string]string
}

func NewIDConverter() *IDConverter {
	return &IDConverter{
		ObjectIdToUUIDMap: make(map[string]string),
		UUIDToObjectIdMap: make(map[string]string),
	}
}

func (c *IDConverter) AddMapping(objectId, uuid string) {
	c.ObjectIdToUUIDMap[objectId] = uuid
	c.UUIDToObjectIdMap[uuid] = objectId
}

func (c *IDConverter) ObjectIdToToUUID(objectId string) (string, bool) {
	uuid, ok := c.ObjectIdToUUIDMap[objectId]
	if !ok {
		uuid = c.simpleObjectIdToUUID(objectId)
		c.ObjectIdToUUIDMap[objectId] = uuid
		c.UUIDToObjectIdMap[uuid] = objectId
		return uuid, true
	}
	return uuid, true
}

func (c *IDConverter) UUIDToObjectId(uuid string) (string, bool) {
	objectId, ok := c.UUIDToObjectIdMap[uuid]
	if !ok {
		return "", false
	}
	return objectId, true
}

func (c *IDConverter) simpleObjectIdToUUID(objectId string) string {
	if len(objectId) == 24 {
		return fmt.Sprintf("%s-%s-%s-%s-%s",
			objectId[0:8], objectId[8:12], objectId[12:16],
			objectId[16:20], objectId[20:24])
	}
	return objectId
}
