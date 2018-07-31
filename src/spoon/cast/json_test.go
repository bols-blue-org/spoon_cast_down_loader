package cast

import (
	"fmt"
	"testing"
)

func TestLoadMetaData(t *testing.T) {
	data, err := LoadMetaData("29236")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", data)
	err = data.Download()
	if err != nil {
		t.Fatal(err)
	}
}
