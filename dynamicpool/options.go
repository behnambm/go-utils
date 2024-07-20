package dynamicpool

type Options func(pool *DynamicWorkerPool) error

//
//func WithErrorChannel(errCh chan error) Options {
//	return func(pool *DynamicWorkerPool) error {
//		pool.errCh = errCh
//		return nil
//	}
//}
