package handlerimplementation

import (
	"encoding/json"
	"main/app/utils"
	"px.dev/pxapi/proto/vizierpb"
	"px.dev/pxapi/types"
)

func ConvertPixieDataToItemStore[itemType ItemType](r *types.Record) itemType {
	var itemStore itemType
	mapObject := map[string]interface{}{}
	for i := 0; i < len(r.Data); i++ {
		tag := r.TableMetadata.ColInfo[i].Name
		datatypeName := vizierpb.DataType_name[int32(r.TableMetadata.ColInfo[i].Type)]
		value := utils.GetData(tag, datatypeName, r)
		mapObject[tag] = value
	}
	jsonStr, err := json.Marshal(mapObject)
	if err != nil {
		println(err)
	}
	err = json.Unmarshal(jsonStr, &itemStore)
	if err != nil {
		println(err)
	}

	return itemStore
}

func GetLatencies(key string, r *types.Record) (Latencies, error) {
	v, _ := utils.GetStringFromRecord(key, r)
	if *v != "" {
		data := Latencies{}
		err := json.Unmarshal([]byte(*v), &data)
		if err != nil {
			return Latencies{}, err
		}
		return data, nil
	}
	return Latencies{}, nil
}

func GetLatenciesPtr(key string, r *types.Record) *Latencies {
	v, err := GetLatencies(key, r)
	if err == nil {
		return &v
	}
	return &Latencies{}
}
