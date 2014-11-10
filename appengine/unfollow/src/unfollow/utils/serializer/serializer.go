package serializer

import (
    "bytes"
    "encoding/gob"
)

func Serialize(source interface{}) ([]byte, error) {
    buffer := new(bytes.Buffer)
    encoder := gob.NewEncoder(buffer)
    if err := encoder.Encode(source); err != nil {
        return nil, err
    }
    return buffer.Bytes(), nil
}

func Deserialize(source []byte, destination interface{}) error {
    decoder := gob.NewDecoder(bytes.NewBuffer(source))
    return decoder.Decode(destination)
}
