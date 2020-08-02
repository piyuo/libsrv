package usage

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/data"

	. "github.com/smartystreets/goconvey/convey"
)

func NewSample(ctx context.Context) (data.DB, error) {
	conn, err := data.FirestoreGlobalConnection(ctx)
	if err != nil {
		return nil, err
	}
	db := &data.BaseDB{
		conn: conn,
	}
	return db, nil
}

func TestUsage(t *testing.T) {
	Convey("Should count,add and remove usage", t, func() {
		ctx := context.Background()

		key := "test_usage"
		count := Count(ctx, key, time.Duration(24)*time.Hour)
		So(count, ShouldEqual, 0)
	})
}
