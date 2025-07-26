package main

import "encoding/json"
import "fmt"

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

// Print bytes as a readable string
func printBytesAsString(data []byte) {
    fmt.Println("Data as String:", string(data))
}

func main() {
    // Example originData
    originData := map[string]interface{}{
        "username": "JohnDoe",
        "password": "secret123",
    }

    // Convert to bytes
    dataBytes, err := originDataToBytes(originData)
    if err != nil {
        fmt.Println("Error converting to bytes:", err)
        return
    }
    
    // Print the byte data as string for debugging
    printBytesAsString(dataBytes)

    // Desensitize
    desensitizedBytes := desensitizeBytes(dataBytes)
    
    // Print the desensitized byte data as string for debugging
    printBytesAsString(desensitizedBytes)

    // Convert back to originData
    newOriginData, err := bytesToOriginData(desensitizedBytes)
    if err != nil {
        fmt.Println("Error converting back to originData:", err)
        return
    }

    fmt.Println("Original Data:", originData)
    fmt.Println("Converted Back Data:", newOriginData)
}
