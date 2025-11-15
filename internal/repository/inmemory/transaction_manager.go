package inmemory

import "context"

type TransactionManager struct{}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

// WithinTransaction выполняет функцию в "транзакции" (для in-memory просто выполняет функцию)
func (tm *TransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}



