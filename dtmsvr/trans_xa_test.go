package dtmsvr

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yedf/dtm/common"
	"github.com/yedf/dtm/dtmcli"
	"github.com/yedf/dtm/examples"
)

func TestXa(t *testing.T) {

	xaNormal(t)
	xaRollback(t)
}

func xaNormal(t *testing.T) {
	xc := examples.XaClient
	gid := "xaNormal"
	err := xc.XaGlobalTransaction(gid, func(xa *dtmcli.Xa) error {
		req := examples.GenTransReq(30, false, false)
		resp, err := xa.CallBranch(req, examples.Busi+"/TransOutXa")
		common.CheckRestySuccess(resp, err)
		resp, err = xa.CallBranch(req, examples.Busi+"/TransInXa")
		common.CheckRestySuccess(resp, err)
		return nil
	})
	e2p(err)
	WaitTransProcessed(gid)
	assert.Equal(t, []string{"prepared", "succeed", "prepared", "succeed"}, getBranchesStatus(gid))
}

func xaRollback(t *testing.T) {
	xc := examples.XaClient
	gid := "xaRollback"
	err := xc.XaGlobalTransaction(gid, func(xa *dtmcli.Xa) error {
		req := examples.GenTransReq(30, false, true)
		resp, err := xa.CallBranch(req, examples.Busi+"/TransOutXa")
		common.CheckRestySuccess(resp, err)
		resp, err = xa.CallBranch(req, examples.Busi+"/TransInXa")
		common.CheckRestySuccess(resp, err)
		return nil
	})
	if err != nil {
		logrus.Errorf("global transaction failed, so rollback")
	}
	WaitTransProcessed(gid)
	assert.Equal(t, []string{"succeed", "prepared"}, getBranchesStatus(gid))
	assert.Equal(t, "failed", getTransStatus(gid))
}