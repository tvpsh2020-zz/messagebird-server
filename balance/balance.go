package balance

import (
	"github.com/messagebird/go-rest-api/balance"
	"github.com/tvpsh2020/messagebird-server/config"
)

// GetBalance will show your prepaid account balance
func GetBalance() (*balance.Balance, error) {
	balance, err := balance.Read(config.MBClient)

	if err != nil {
		return nil, err
	}

	return balance, nil
}
