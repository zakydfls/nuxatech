package service

import (
	"context"
	"errors"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/model"
	"nuxatech-nextmedis/repository"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req *request.CreateAccountRequest) (*response.AccountResponse, error)
	GetAccount(ctx context.Context, id string) (*response.AccountResponse, error)
	Deposit(ctx context.Context, accountID string, req *request.TransactionRequest) (*response.TransactionResponse, error)
	Withdraw(ctx context.Context, accountID string, req *request.TransactionRequest) (*response.TransactionResponse, error)
}

type accountService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	validate        *validator.Validate
	mutexes         sync.Map
}

func (s *accountService) getLock(accountID string) *sync.Mutex {
	mutex, _ := s.mutexes.LoadOrStore(accountID, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

func (s *accountService) CreateAccount(ctx context.Context, req *request.CreateAccountRequest) (*response.AccountResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	account := &model.Account{
		UserID:    req.UserID,
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx := s.accountRepo.BeginTx(ctx)
	if err := s.accountRepo.CreateAccount(ctx, tx, account); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &response.AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*response.AccountResponse, error) {
	account, err := s.accountRepo.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}

func (s *accountService) Deposit(ctx context.Context, accountID string, req *request.TransactionRequest) (*response.TransactionResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	if err := s.validateAmount(req.Amount); err != nil {
		return nil, err
	}

	// Get account-specific lock
	lock := s.getLock(accountID)
	lock.Lock()
	defer lock.Unlock()

	// Start database transaction
	tx := s.accountRepo.BeginTx(ctx)
	if tx == nil {
		return nil, errors.New("failed to start transaction")
	}
	defer tx.Rollback()

	account, err := s.accountRepo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	now := time.Now().UnixMilli()
	transaction := &model.Transaction{
		AccountID:   accountID,
		Amount:      req.Amount,
		Type:        "deposit",
		Status:      "processing",
		Description: req.Description,
		CreatedAt:   now,
	}

	if err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return nil, err
	}

	newBalance := account.Balance + req.Amount
	if err := s.accountRepo.UpdateBalance(ctx, tx, accountID, newBalance); err != nil {
		if updateErr := tx.Model(&model.Transaction{}).
			Where("id = ?", transaction.ID).
			Update("status", "failed").Error; updateErr != nil {
			return nil, updateErr
		}
		return nil, err
	}
	if err := tx.Model(&model.Transaction{}).
		Where("id = ?", transaction.ID).
		Update("status", "success").Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	transaction.Status = "success"

	return &response.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Status:      transaction.Status,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
	}, nil
}

func (s *accountService) Withdraw(ctx context.Context, accountID string, req *request.TransactionRequest) (*response.TransactionResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	if err := s.validateAmount(req.Amount); err != nil {
		return nil, err
	}

	// Get account-specific lock
	lock := s.getLock(accountID)
	lock.Lock()
	defer lock.Unlock()

	// Start database transaction
	tx := s.accountRepo.BeginTx(ctx)
	if tx == nil {
		return nil, errors.New("failed to start transaction")
	}
	defer tx.Rollback()

	account, err := s.accountRepo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	now := time.Now().UnixMilli()
	transaction := &model.Transaction{
		AccountID:   accountID,
		Amount:      req.Amount,
		Type:        "withdrawal",
		Status:      "processing",
		Description: req.Description,
		CreatedAt:   now,
	}

	if err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return nil, err
	}

	newBalance := account.Balance - req.Amount
	if err := s.accountRepo.UpdateBalance(ctx, tx, accountID, newBalance); err != nil {
		if updateErr := tx.Model(&model.Transaction{}).
			Where("id = ?", transaction.ID).
			Update("status", "failed").Error; updateErr != nil {
			return nil, updateErr
		}
		return nil, err
	}
	if err := tx.Model(&model.Transaction{}).
		Where("id = ?", transaction.ID).
		Update("status", "success").Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	transaction.Status = "success"

	return &response.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Status:      transaction.Status,
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt,
	}, nil
}

func (s *accountService) validateAmount(amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	maxAmount := int64(1000000000)
	if amount > maxAmount {
		return errors.New("amount exceeds maximum allowed")
	}

	return nil
}

func NewAccountService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) AccountService {
	return &accountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		validate:        validator.New(),
	}
}
