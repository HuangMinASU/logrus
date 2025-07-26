import (
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "time"
)

// Fire is a hook to handle writing to local log files
func (h DiHook) Fire(entry *logrus.Entry) error {
    if entry.Level > h.level {
        return nil
    }

    originData := entry.Data

    // Convert originData to bytes
    dataBytes, err := originDataToBytes(originData)
    if err != nil {
        return fmt.Errorf("error converting originData to bytes: %v", err)
    }

    // Desensitize the byte data
    desensitizedBytes := desensitizeBytes(dataBytes)

    // Convert back to originData
    newOriginData, err := bytesToOriginData(desensitizedBytes)
    if err != nil {
        return fmt.Errorf("error converting bytes back to originData: %v", err)
    }

    // Replace entry.Data with the modified version
    entry.Data = make(logrus.Fields, len(h.extFields)+2) // +2 to include Seq and Data

    for k, v := range h.extFields {
        entry.Data[k] = v
    }
    entry.Data["Seq"] = h.seq
    h.seq++

    if len(newOriginData) != 0 {
        entry.Data["Data"] = newOriginData
    }

    entry.Time = entry.Time.In(h.local)
    return nil
}

// Convert originData to byte slice
func originDataToBytes(originData map[string]interface{}) ([]byte, error) {
    byteData, err := json.Marshal(originData)
    if err != nil {
        return nil, err
    }
    return byteData, nil
}

// Convert bytes back to originData (map[string]interface{})
func bytesToOriginData(data []byte) (map[string]interface{}, error) {
    var originData map[string]interface{}
    err := json.Unmarshal(data, &originData)
    if err != nil {
        return nil, err
    }
    return originData, nil
}

// Example: Simple desensitization function for bytes
func desensitizeBytes(data []byte) []byte {
    for i := range data {
        if ('A' <= data[i] && data[i] <= 'Z') || ('a' <= data[i] && data[i] <= 'z') {
            data[i] = '*'
        }
    }
    return data
}
