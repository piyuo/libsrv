package key

import (
	"strings"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("should get key path", t, func() {
		path, err := getKeyPath("gcloud.json")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/gcloud.json"), ShouldBeTrue)
	})
	Convey("should get region key path", t, func() {
		path, err := getKeyPath("region/us.json")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/region/us.json"), ShouldBeTrue)
	})
	Convey("should get key content", t, func() {
		text, err := Text("gcloud.json")
		So(err, ShouldBeNil)
		So(text, ShouldNotBeEmpty)

		bytes, err := Bytes("gcloud.json")
		So(err, ShouldBeNil)
		So(bytes, ShouldNotBeNil)

		json, err := JSON("gcloud.json")
		So(err, ShouldBeNil)
		So(json["project_id"], ShouldNotBeEmpty)

	})
	Convey("should get region key content", t, func() {
		text, err := getKeyPath("region/us.json")
		So(err, ShouldBeNil)
		So(text, ShouldNotBeEmpty)
	})

	Convey("should return error when key not exist", t, func() {
		path, err := getKeyPath("not exist")
		So(err, ShouldNotBeNil)
		So(path, ShouldBeEmpty)
	})

	Convey("should return error when key not exist", t, func() {
		content, err := Text("not exist")
		So(err, ShouldNotBeNil)
		So(content, ShouldBeEmpty)

		json, err := JSON("not exist")
		So(err, ShouldNotBeNil)
		So(json, ShouldBeNil)

		bytes, err := Bytes("not exist")
		So(err, ShouldNotBeNil)
		So(bytes, ShouldBeNil)
	})
}

func TestConcurrentKey(t *testing.T) {
	var concurrent = 10
	var wg sync.WaitGroup
	wg.Add(concurrent)
	encryptDecrypt := func() {
		for i := 0; i < 10; i++ {
			_, err := Text("gcloud.json")
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			_, err = JSON("gcloud.json")
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			_, err = Bytes("gcloud.json")
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			//fmt.Print(text + "\n")
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go encryptDecrypt()
	}
	wg.Wait()
}
