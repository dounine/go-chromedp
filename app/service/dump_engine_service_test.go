package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDumpEngine(t *testing.T) {
	ass := assert.New(t)
	ctx := context.Background()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	dumpEngineService := NewDumpEngineService()
	err := dumpEngineService.Run(&ctx, "", "../../images")
	ass.NoError(err)
}

func TestQueryDumps(t *testing.T) {
	ass := assert.New(t)
	ctx := context.Background()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	dumpEngineService := NewDumpEngineService()
	dumps, err := dumpEngineService.Dumps(&ctx)
	ass.NoError(err)
	ass.NotEmpty(dumps)
	t.Logf("dumps: %+v", dumps)
}
