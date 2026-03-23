package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/db/common"
)

// MappingManager ID映射管理器
type MappingManager struct {
	ObjectIdToUUIDMap map[string]string
	UUIDToObjectIdMap map[string]string
	mappingFile       string
}

func NewMappingManager(mappingFile string) *MappingManager {
	return &MappingManager{
		ObjectIdToUUIDMap: make(map[string]string),
		UUIDToObjectIdMap: make(map[string]string),
		mappingFile:       mappingFile,
	}
}

func (mm *MappingManager) Load() error {
	if mm.mappingFile == "" {
		return nil
	}

	data, err := ioutil.ReadFile(mm.mappingFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在是正常的
		}
		return fmt.Errorf("failed to read mapping file: %w", err)
	}

	var mappings []IDMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		return fmt.Errorf("failed to parse mapping file: %w", err)
	}

	mm.ObjectIdToUUIDMap = make(map[string]string)
	mm.UUIDToObjectIdMap = make(map[string]string)

	for _, mapping := range mappings {
		mm.ObjectIdToUUIDMap[mapping.ObjectId] = mapping.UUID
		mm.UUIDToObjectIdMap[mapping.UUID] = mapping.ObjectId
	}

	return nil
}

func (mm *MappingManager) Save() error {
	if mm.mappingFile == "" {
		return nil
	}

	var mappings []IDMapping
	now := time.Now().Format(time.RFC3339)

	// 收集所有映射（只保存从ObjectId到UUID的映射）
	for objectId, uuid := range mm.ObjectIdToUUIDMap {
		mappings = append(mappings, IDMapping{
			ObjectId:   objectId,
			UUID:       uuid,
			TableName:  "", // 可以为空
			MigratedAt: now,
			Direction:  "mongo_to_pg",
		})
	}

	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mappings: %w", err)
	}

	if err := ioutil.WriteFile(mm.mappingFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write mapping file: %w", err)
	}

	return nil
}

func (mm *MappingManager) Add(objectId, uuid, table string) {
	mm.ObjectIdToUUIDMap[objectId] = uuid
	mm.UUIDToObjectIdMap[uuid] = objectId
}

func (mm *MappingManager) ObjectIdToUUID(objectId string) (string, bool) {
	uuid, ok := mm.ObjectIdToUUIDMap[objectId]
	if !ok {
		// 如果没有映射，使用简单转换
		uuid = mm.simpleObjectIdToUUID(objectId)
		mm.ObjectIdToUUIDMap[objectId] = uuid
		mm.UUIDToObjectIdMap[uuid] = objectId
		return uuid, true
	}
	return uuid, true
}

func (mm *MappingManager) UUIDToObjectId(uuid string) (string, bool) {
	objectId, ok := mm.UUIDToObjectIdMap[uuid]
	return objectId, ok
}

func (mm *MappingManager) simpleObjectIdToUUID(objectId string) string {
	if len(objectId) == 24 {
		return fmt.Sprintf("%s-%s-%s-%s-%s",
			objectId[0:8], objectId[8:12], objectId[12:16],
			objectId[16:20], objectId[20:24])
	}
	return objectId
}
