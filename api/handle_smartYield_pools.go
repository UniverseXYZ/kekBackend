package api

import (
	"database/sql"
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
			   underlying_decimals,
			   coalesce((select pool_address from smart_yield_reward_pools where pool_token_address = sy_address ), '') as reward_pool_address
		from smart_yield_pools p
		where sy_address = $1
	`, poolAddress).Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals, &p.RewardPoolAddress)
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
			select included_in_block,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   number_of_seniors(pool_address)                      as number_of_seniors,
				   number_of_jtoken_holders(pool_address)               as number_of_juniors,
			       number_of_juniors_locked(pool_address)               as number_of_juniors_locked,
				    (abond_matures_at - (select extract (epoch from now())))::double precision / (60*60*24)     as avg_senior_buy,
				   coalesce(junior_liquidity_locked(pool_address), 0)   as junior_liquidity_locked
			from smart_yield_state
			where pool_address = $1
			order by included_in_block desc
			limit 1;
		`, p.SmartYieldAddress).Scan(&state.BlockNumber, &state.BlockTimestamp, &state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY, &state.OriginatorApy, &state.OriginatorNetApy, &state.NumberOfSeniors, &state.NumberOfJuniors, &state.NumberOfJuniorsLocked, &state.AvgSeniorMaturityDays, &state.JuniorLiquidityLocked)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

	state.JuniorLiquidityLocked = state.JuniorLiquidityLocked.Div(tenPowDec)
	state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
	state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
	state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

	p.State = state

	OK(c, p)
}

func (a *API) handlePools(c *gin.Context) {
	protocols := strings.ToLower(c.DefaultQuery("originator", "all"))
	underlyingSymbol := strings.ToUpper(c.DefaultQuery("underlyingSymbol", "all"))

	filters := new(Filters)
	if protocols != "all" {
		protocolsArray := strings.Split(protocols, ",")
		filters.Add("protocol_id", protocolsArray)
	}

	if underlyingSymbol != "ALL" {
		filters.Add("underlying_symbol", underlyingSymbol)
	}

	query, params := buildQueryWithFilter(`select protocol_id,
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
			   underlying_decimals,
			   coalesce((select pool_address from smart_yield_reward_pools where pool_token_address = sy_address ), '') as reward_pool_address
		from smart_yield_pools p
		%s
		%s %s`,
		filters,
		nil,
		nil)

	rows, err := a.db.Query(query, params...)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	var pools []types.SYPool
	for rows.Next() {
		var p types.SYPool

		err := rows.Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals, &p.RewardPoolAddress)
		if err != nil {
			Error(c, err)
			return
		}

		var state types.SYPoolState
		err = a.db.QueryRow(`
			select included_in_block,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   number_of_seniors(pool_address)                      as number_of_seniors,
				   number_of_jtoken_holders(pool_address)               as number_of_juniors,
			       number_of_juniors_locked(pool_address)               as number_of_juniors_locked,
				   (abond_matures_at - (select extract (epoch from now())))::double precision / (60*60*24)     as avg_senior_buy,
				   coalesce(junior_liquidity_locked(pool_address), 0)   as junior_liquidity_locked
			from smart_yield_state
			where pool_address = $1
			order by included_in_block desc
			limit 1;
		`, p.SmartYieldAddress).Scan(&state.BlockNumber, &state.BlockTimestamp, &state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY, &state.OriginatorApy, &state.OriginatorNetApy, &state.NumberOfSeniors, &state.NumberOfJuniors, &state.NumberOfJuniorsLocked, &state.AvgSeniorMaturityDays, &state.JuniorLiquidityLocked)
		if err != nil && err != sql.ErrNoRows {
			Error(c, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

		state.JuniorLiquidityLocked = state.JuniorLiquidityLocked.Div(tenPowDec)
		state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
		state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
		state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

		p.State = state

		pools = append(pools, p)
	}

	OK(c, pools)
}

func (a *API) handleRewardPools(c *gin.Context) {
	protocols := strings.ToLower(c.DefaultQuery("protocolId", "all"))
	underlyingSymbol := strings.ToUpper(c.DefaultQuery("underlyingSymbol", "all"))
	underlyingAddress := strings.ToLower(c.DefaultQuery("underlyingAddress", "all"))

	filters := new(Filters)
	if protocols != "all" {
		protocolsArray := strings.Split(protocols, ",")
		filters.Add("p.protocol_id", protocolsArray)
	}

	if underlyingSymbol != "ALL" {
		filters.Add("p.underlying_symbol", underlyingSymbol)
	}

	if underlyingAddress != "all" {
		filters.Add("p.underlying_address", utils.NormalizeAddress(underlyingAddress))
	}

	query, params := buildQueryWithFilter(`select 
				       r.pool_address,
				       r.pool_token_address,
				       r.reward_token_address,
				       p.underlying_decimals,
				       p.protocol_id,
				       p.underlying_symbol,
	                   p.underlying_address
				from smart_yield_reward_pools as r
				inner join smart_yield_pools as p
				on p.sy_address = r.pool_token_address %s 
				%s %s`,
		filters,
		nil,
		nil)
	var pools []types.SYRewardPool
	rows, err := a.db.Query(query, params...)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	for rows.Next() {
		var p types.SYRewardPool
		err := rows.Scan(&p.PoolAddress, &p.PoolTokenAddress, &p.RewardTokenAddress, &p.PoolTokenDecimals, &p.ProtocolID, &p.UnderlyingSymbol, &p.UnderlyingAddress)
		if err != nil {
			Error(c, err)
			return
		}

		pools = append(pools, p)
	}

	OK(c, pools)
}
