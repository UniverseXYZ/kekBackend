//go:generate abigen --abi ../abis/Governance.json --pkg contracts --out Governance.go --type Governance
//go:generate abigen --abi ../abis/YieldFarmContinuous.json --pkg contracts --out YieldFarmContinuous.go --type YieldFarmContinuous
package contracts
