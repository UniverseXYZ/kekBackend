package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) handlePoolDetails(c *gin.Context) {
	pool := c.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		BadRequest(c, errors.New("invalid pool address"))
		return
	}

	var p types.SYPool

	err = a.db.QueryRow(`
		select protocol_id,
			   controller_address,
			   model_address,
			   provider_address,
			   sy_address,
			   oracle_address,
			   junior_bond_address,
			   senior_bond_address,
			   receipt_token_address,
			   underlying_address,
			   underlying_symbol,
			   underlying_decimals
		from smart_yield_pools p
		where sy_address = $1
	`, poolAddress).Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	var state types.SYPoolState
	err = a.db.QueryRow(`
			select block_number,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   (select count(*) from smart_yield_senior_buy where sy_address = pool_address ) as number_of_seniors,
				   coalesce((select avg(for_days) from smart_yield_senior_buy where sy_address = pool_address), 0) as avg_senior_buy,
				   (select count(*) from smart_yield_token_buy where sy_address = pool_address ) as number_of_juniors
			from smart_yield_state
			where pool_address = $1
			order by block_number desc
			limit 1
		`, p.SmartYieldAddress).Scan(&state.BlockNumber, &state.BlockTimestamp, &state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY, &state.OriginatorApy, &state.OriginatorNetApy, &state.NumberOfSeniors, &state.AvgSeniorMaturityDays, &state.NumberOfJuniors)
	if err != nil {
		Error(c, err)
		return
	}

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

	state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
	state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
	state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

	p.State = state

	OK(c, p)
}

func (a *API) handlePools(c *gin.Context) {
	protocolID := strings.ToLower(c.DefaultQuery("protocolID", "all"))

	var pools []types.SYPool

	query := `
		select protocol_id,
			   controller_address,
			   model_address,
			   provider_address,
			   sy_address,
			   oracle_address,
			   junior_bond_address,
			   senior_bond_address,
			   receipt_token_address,
			   underlying_address,
			   underlying_symbol,
			   underlying_decimals
		from smart_yield_pools p
		where 1 = 1 %s
	`

	var parameters []interface{}

	if protocolID == "all" {
		query = fmt.Sprintf(query, "")
	} else {
		protocolFilter := fmt.Sprintf("and protocol_id = $1")
		parameters = append(parameters, protocolID)
		query = fmt.Sprintf(query, protocolFilter)
	}

	rows, err := a.db.Query(query, parameters...)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	for rows.Next() {
		var p types.SYPool

		err := rows.Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals)
		if err != nil {
			Error(c, err)
			return
		}

		var state types.SYPoolState
		err = a.db.QueryRow(`
			select block_number,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   (select count(*) from smart_yield_senior_buy where sy_address = pool_address ) as number_of_seniors,
				   coalesce((select avg(for_days) from smart_yield_senior_buy where sy_address = pool_address), 0) as avg_senior_buy,
				   (select count(*) from smart_yield_token_buy where sy_address = pool_address ) as number_of_juniors
			from smart_yield_state
			where pool_address = $1
			order by block_number desc
			limit 1
		`, p.SmartYieldAddress).Scan(&state.BlockNumber, &state.BlockTimestamp, &state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY, &state.OriginatorApy, &state.OriginatorNetApy, &state.NumberOfSeniors, &state.AvgSeniorMaturityDays, &state.NumberOfJuniors)
		if err != nil && err != sql.ErrNoRows {
			Error(c, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

		state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
		state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
		state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

		p.State = state

		pools = append(pools, p)
	}

	OK(c, pools)
}
