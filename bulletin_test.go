package bulletin

import (
	"os"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test001BulletinPersistanceToDisk(t *testing.T) {

	cv.Convey("read/write to persistent disk should work", t, func() {
		testPath := "./test-db"
		os.RemoveAll(testPath)
		os.MkdirAll(testPath, 0777)
		//defer os.Remove(testPath)

		cfg := &Config{
			MossPath: testPath,
		}
		b := NewBoard(cfg)
		err := b.openMoss()
		panicOn(err)
		k := []byte("a key")
		v := []byte("a value")

		err = b.writeMossKV(&KV{Key: k, Val: v})
		panicOn(err)
		err = b.closeMoss()
		panicOn(err)

		// verify
		err = b.openMoss()
		panicOn(err)

		v2, err := b.readMoss(k)
		panicOn(err)

		cv.So(v2, cv.ShouldResemble, v)
	})
}
