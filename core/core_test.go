package core
/*
import (
	"testing"
	"fmt"
	"net/http"
)

var test_private_key = ""


func TestUpdate(t *testing.T) {
	db := Database{
		PrivateKey: &test_private_key,
		TrackerKey: &test_private_key,
		Source:     make(map[string]interface{}),
	}

	if err := db.Add("testkey", "hello"); err != nil {
		t.Fatalf("Some error %s", err)
	}
}

func TestConnect(t *testing.T) {
	db, _ := HardConnect(test_private_key)
	if db != nil {
		t.Log(db)
	} else if db == nil{
		t.Fatalf("Could not get database")
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("http://localhost/add?key=%d&value=%d", i, i)
		http.Get(url)
	}
}
*/
