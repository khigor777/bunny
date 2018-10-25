package bunny

import "testing"

func TestOpen(t *testing.T) {
	db := Open()
	in := []struct {
		key string
		val []byte
	}{
		{"2", []byte("2")},
		{"key", []byte("val")},
		{"keys", []byte("esdfkeswfjreut934-59-3u59tu4tgjoreg")},
	}
	for _, v := range in {
		db.Set(v.key, v.val)
		if string(db.Get(v.key)) != string(v.val) {
			t.Error("not equal value", v.key)
			db.Delete(v.key)
			if db.Get(v.key) != nil {
				t.Error("key did't delete")
			}
		}
	}

}
