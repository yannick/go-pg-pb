package timestamp

import (
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10/types"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func tsappender(in []byte, v reflect.Value, flags int) []byte {
	pbts := v.Interface().(timestamppb.Timestamp)
	tm, err := ptypes.Timestamp(&pbts)
	if err != nil {
		panic(err)
	}

	return types.AppendTime(in, tm, flags)

}

func tscanner(v reflect.Value, rd types.Reader, n int) error {
	if !v.CanSet() {
		return fmt.Errorf("pg: Scan(nonsettable %s)", v.Type())
	}

	tm, err := types.ScanTime(rd, n)
	if err != nil {
		return err
	}
	pt, err := ptypes.TimestampProto(tm)
	if err != nil {
		return err
	}

	ptr := v.Addr().Interface().(*timestamppb.Timestamp)
	*ptr = *pt
	return nil
}

func init() {
	var x timestamppb.Timestamp
	types.RegisterAppender(x, tsappender)
	types.RegisterScanner(x, tscanner)
}
