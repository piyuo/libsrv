package data

import (
    "testing"
     . "github.com/smartystreets/goconvey/convey"
)

func TestEnv(t *testing.T) {

    db,err := ProviderInstance().NewDB()
    Convey("NewDB should not have error", t, func() {
        So(err, ShouldBeNil)
    })
    Convey("Fail to get db", t, func() {
        So(db, ShouldNotBeNil)
      })
}
