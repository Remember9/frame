package encode

import "github.com/Remember9/frame/util/xerrors"

const (
	SuccessCode = 0 // 成功返回码
)

var (
	HealthError = xerrors.New(500, "invalid health")

	DataExistsError = xerrors.New(1010001, "已存在")

	ErrAesKeyLengthSixteen = xerrors.New(1020001, "a sixteen or twenty-four or thirty-two length secret key is required")
	ErrAesIv               = xerrors.New(1020002, "a sixteen-length ivAes is required")
	ErrAesPaddingSize      = xerrors.New(1020003, "padding size error please check the secret key or iv")
)
