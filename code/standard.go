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

	// system error
	Maintain          = 18 //維護中
	GameTypeError     = 19
	AccountTypeError  = 20
	NoThisGameAccount = 21
	NoThisPlayer      = 22
	NoMoneyToBet      = 23

	// ulg error
	AuthorizedError    = 50
	ExchangeError      = 51
	GetUserError       = 52
	NoExchange         = 53
	ULGInfoFormatError = 54
	NoCheckoutError    = 55
	CheckoutError      = 56

	RoomLock         = 1001
	SelfInRoom       = 1002
	SelfNotInRoom    = 1003
	RoomFull         = 1004
	RoomNotExistence = 1005
	AlreadyInGame    = 1006
)
