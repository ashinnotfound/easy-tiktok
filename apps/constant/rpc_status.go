package constant

type rpcStatus struct {
	statusOk     *int32
	statusFailed *int32
}

var RPC_STATUS rpcStatus

func init() {
	RPC_STATUS.statusOk = new(int32)
	RPC_STATUS.statusFailed = new(int32)
	*RPC_STATUS.statusOk = 0
	*RPC_STATUS.statusFailed = 400
}

func (x *rpcStatus) StatusOK() *int32 {
	return x.statusOk
}

func (x *rpcStatus) StatusFailed() *int32 {
	return x.statusFailed
}
