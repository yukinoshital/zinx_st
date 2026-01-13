package ziface

type Iroute interface {
	PreHandle(Irequest)
	Handle(Irequest)
	AfterHandle(Irequest)
}