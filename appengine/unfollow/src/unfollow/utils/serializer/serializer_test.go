package serializer

import (
    "reflect"
    "testing"
)

var objects = []interface{}{
    "Hello, 世界!",
    123,
    1.2,
    'f',
}

func TestSerializer(t *testing.T) {
    for _, object := range objects {
        serialized, err := Serialize(&object)
        if err != nil {
            t.Fatal(object, err)
        }

        deserialized := object
        if err := Deserialize(serialized, &deserialized); err != nil {
            t.Fatal(err)
        }

        if !reflect.DeepEqual(object, deserialized) {
            t.Fatal("serializer: unequal objects", object, deserialized)
        }
    }
}
