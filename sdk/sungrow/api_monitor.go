package sungrow

import (
	"errors"
)

type PageRequest struct {
	CurPage int `json:"curPage"`
	Size    int `json:"size,omitempty"`
}

// 电站列表信息查询

type PowerStationListRequest struct {
	Request

	PageRequest
	PsName    string `json:"ps_name,omitempty"`
	PsType    string `json:"ps_type,omitempty"`
	ShareType string `json:"share_type,omitempty"`
	ValidFlag string `json:"valid_flag,omitempty"`
	OrgID     string `json:"org_id,omitempty"`
}

func (sdk *SungrowSDK) GetPowerStationList(req PowerStationListRequest) error {
	if req.CurPage < 0 {
		return errors.New("curPage 参数必须大于0")
	}

	// appkey 与 token 由 callOnce 统一注入，token 为空或失效时自动获取并缓存
	var out any
	if err := sdk.do(PathGetPsList, &req, &out); err != nil {
		return err
	}

	return nil
}
