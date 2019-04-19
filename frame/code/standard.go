package code

import "google.golang.org/grpc/codes"

// ...
const (
	OK                 = Code(codes.OK)
	Canceled           = Code(codes.Canceled)
	Unknown            = Code(codes.Unknown)
	InvalidArgument    = Code(codes.InvalidArgument)    // 無效
	DeadlineExceeded   = Code(codes.DeadlineExceeded)   // 超時
	NotFound           = Code(codes.NotFound)           // 未找到
	AlreadyExists      = Code(codes.AlreadyExists)      // 已存在
	PermissionDenied   = Code(codes.PermissionDenied)   // 無權限
	ResourceExhausted  = Code(codes.ResourceExhausted)  // 資源耗盡
	FailedPrecondition = Code(codes.FailedPrecondition) // 請求失敗
	Aborted            = Code(codes.Aborted)            // 中止
	OutOfRange         = Code(codes.OutOfRange)         // 超出範圍
	Unimplemented      = Code(codes.Unimplemented)      // 未實現
	Internal           = Code(codes.Internal)           // 內部
	Unavailable        = Code(codes.Unavailable)        // 不可用
	DataLoss           = Code(codes.DataLoss)           // 資料遺失
	Unauthenticated    = Code(codes.Unauthenticated)    // 未認證
	Maintain           = Code(18)                       //維護中

	RoomLockError    = Code(1001)
	SelfInRoom       = Code(1002)
	SelfNotInRoom    = Code(1003)
	RoomFull         = Code(1004)
	RoomNotExistence = Code(1005)
)
