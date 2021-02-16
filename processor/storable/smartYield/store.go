package smartYield

import (
	"database/sql"

	"github.com/pkg/errors"
)

func (s *Storable) storeProcessed(tx *sql.Tx) error {
	err := s.storeTokenBuyTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store token buy trades")
	}

	err = s.storeTokenSellTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store token sell trades")
	}

	err = s.storeJuniorBuyTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store junior buy trades")
	}

	err = s.storeJuniorRedeemTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store junior redeem trades")
	}

	err = s.storeSeniorBuyTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store senior buy trades")
	}

	err = s.storeSeniorRedeemTrades(tx)
	if err != nil {
		return errors.Wrap(err, "could not store senior redeem trades")
	}

	return nil
}
