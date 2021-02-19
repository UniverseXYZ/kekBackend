package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/barnbridge-backend/api/types"
)

func (a *API) handlePools(c *gin.Context) {
	protocolID := strings.ToLower(c.DefaultQuery("protocolID", "all"))

	var pools []types.SYPool

	query := `select protocol_id,
					controller_address,
					model_address,
					provider_address,
					sy_address,
					oracle_address,
					junior_bond_address,
					senior_bond_address,
					receipt_token_address,
					underlying_address,
					underlying_symbol
					,underlying_decimals
	from smart_yield_pools where 1=1
	%s`

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

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	for rows.Next() {
		var p types.SYPool

		err := rows.Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals)
		if err != nil {
			Error(c, err)
			return
		}

		pools = append(pools, p)
	}
	OK(c, pools)
}
