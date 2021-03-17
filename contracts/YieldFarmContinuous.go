// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// YieldFarmContinuousABI is the input ABI used to generate the binding from.
const YieldFarmContinuousABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_poolToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Claim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balanceAfter\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balanceAfter\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ackFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"balanceBefore\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claim\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastSoftPullTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"owed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pullRewardFromSource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardNotTransferred\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardRatePerSecond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardSource\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"setRewardRatePerSecond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"}],\"name\":\"setRewardsSource\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"softPullReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"userMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawAndClaim\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// YieldFarmContinuous is an auto generated Go binding around an Ethereum contract.
type YieldFarmContinuous struct {
	YieldFarmContinuousCaller     // Read-only binding to the contract
	YieldFarmContinuousTransactor // Write-only binding to the contract
	YieldFarmContinuousFilterer   // Log filterer for contract events
}

// YieldFarmContinuousCaller is an auto generated read-only Go binding around an Ethereum contract.
type YieldFarmContinuousCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// YieldFarmContinuousTransactor is an auto generated write-only Go binding around an Ethereum contract.
type YieldFarmContinuousTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// YieldFarmContinuousFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type YieldFarmContinuousFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// YieldFarmContinuousSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type YieldFarmContinuousSession struct {
	Contract     *YieldFarmContinuous // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// YieldFarmContinuousCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type YieldFarmContinuousCallerSession struct {
	Contract *YieldFarmContinuousCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// YieldFarmContinuousTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type YieldFarmContinuousTransactorSession struct {
	Contract     *YieldFarmContinuousTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// YieldFarmContinuousRaw is an auto generated low-level Go binding around an Ethereum contract.
type YieldFarmContinuousRaw struct {
	Contract *YieldFarmContinuous // Generic contract binding to access the raw methods on
}

// YieldFarmContinuousCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type YieldFarmContinuousCallerRaw struct {
	Contract *YieldFarmContinuousCaller // Generic read-only contract binding to access the raw methods on
}

// YieldFarmContinuousTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type YieldFarmContinuousTransactorRaw struct {
	Contract *YieldFarmContinuousTransactor // Generic write-only contract binding to access the raw methods on
}

// NewYieldFarmContinuous creates a new instance of YieldFarmContinuous, bound to a specific deployed contract.
func NewYieldFarmContinuous(address common.Address, backend bind.ContractBackend) (*YieldFarmContinuous, error) {
	contract, err := bindYieldFarmContinuous(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuous{YieldFarmContinuousCaller: YieldFarmContinuousCaller{contract: contract}, YieldFarmContinuousTransactor: YieldFarmContinuousTransactor{contract: contract}, YieldFarmContinuousFilterer: YieldFarmContinuousFilterer{contract: contract}}, nil
}

// NewYieldFarmContinuousCaller creates a new read-only instance of YieldFarmContinuous, bound to a specific deployed contract.
func NewYieldFarmContinuousCaller(address common.Address, caller bind.ContractCaller) (*YieldFarmContinuousCaller, error) {
	contract, err := bindYieldFarmContinuous(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousCaller{contract: contract}, nil
}

// NewYieldFarmContinuousTransactor creates a new write-only instance of YieldFarmContinuous, bound to a specific deployed contract.
func NewYieldFarmContinuousTransactor(address common.Address, transactor bind.ContractTransactor) (*YieldFarmContinuousTransactor, error) {
	contract, err := bindYieldFarmContinuous(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousTransactor{contract: contract}, nil
}

// NewYieldFarmContinuousFilterer creates a new log filterer instance of YieldFarmContinuous, bound to a specific deployed contract.
func NewYieldFarmContinuousFilterer(address common.Address, filterer bind.ContractFilterer) (*YieldFarmContinuousFilterer, error) {
	contract, err := bindYieldFarmContinuous(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousFilterer{contract: contract}, nil
}

// bindYieldFarmContinuous binds a generic wrapper to an already deployed contract.
func bindYieldFarmContinuous(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(YieldFarmContinuousABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_YieldFarmContinuous *YieldFarmContinuousRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _YieldFarmContinuous.Contract.YieldFarmContinuousCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_YieldFarmContinuous *YieldFarmContinuousRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.YieldFarmContinuousTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_YieldFarmContinuous *YieldFarmContinuousRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.YieldFarmContinuousTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_YieldFarmContinuous *YieldFarmContinuousCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _YieldFarmContinuous.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_YieldFarmContinuous *YieldFarmContinuousTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_YieldFarmContinuous *YieldFarmContinuousTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.contract.Transact(opts, method, params...)
}

// BalanceBefore is a free data retrieval call binding the contract method 0x94b5798a.
//
// Solidity: function balanceBefore() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) BalanceBefore(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "balanceBefore")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceBefore is a free data retrieval call binding the contract method 0x94b5798a.
//
// Solidity: function balanceBefore() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) BalanceBefore() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.BalanceBefore(&_YieldFarmContinuous.CallOpts)
}

// BalanceBefore is a free data retrieval call binding the contract method 0x94b5798a.
//
// Solidity: function balanceBefore() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) BalanceBefore() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.BalanceBefore(&_YieldFarmContinuous.CallOpts)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _YieldFarmContinuous.Contract.Balances(&_YieldFarmContinuous.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _YieldFarmContinuous.Contract.Balances(&_YieldFarmContinuous.CallOpts, arg0)
}

// CurrentMultiplier is a free data retrieval call binding the contract method 0x6fbaaa1e.
//
// Solidity: function currentMultiplier() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) CurrentMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "currentMultiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentMultiplier is a free data retrieval call binding the contract method 0x6fbaaa1e.
//
// Solidity: function currentMultiplier() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) CurrentMultiplier() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.CurrentMultiplier(&_YieldFarmContinuous.CallOpts)
}

// CurrentMultiplier is a free data retrieval call binding the contract method 0x6fbaaa1e.
//
// Solidity: function currentMultiplier() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) CurrentMultiplier() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.CurrentMultiplier(&_YieldFarmContinuous.CallOpts)
}

// LastSoftPullTs is a free data retrieval call binding the contract method 0x5ede634d.
//
// Solidity: function lastSoftPullTs() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) LastSoftPullTs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "lastSoftPullTs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastSoftPullTs is a free data retrieval call binding the contract method 0x5ede634d.
//
// Solidity: function lastSoftPullTs() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) LastSoftPullTs() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.LastSoftPullTs(&_YieldFarmContinuous.CallOpts)
}

// LastSoftPullTs is a free data retrieval call binding the contract method 0x5ede634d.
//
// Solidity: function lastSoftPullTs() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) LastSoftPullTs() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.LastSoftPullTs(&_YieldFarmContinuous.CallOpts)
}

// Owed is a free data retrieval call binding the contract method 0xdf18e047.
//
// Solidity: function owed(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) Owed(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "owed", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Owed is a free data retrieval call binding the contract method 0xdf18e047.
//
// Solidity: function owed(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) Owed(arg0 common.Address) (*big.Int, error) {
	return _YieldFarmContinuous.Contract.Owed(&_YieldFarmContinuous.CallOpts, arg0)
}

// Owed is a free data retrieval call binding the contract method 0xdf18e047.
//
// Solidity: function owed(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) Owed(arg0 common.Address) (*big.Int, error) {
	return _YieldFarmContinuous.Contract.Owed(&_YieldFarmContinuous.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousSession) Owner() (common.Address, error) {
	return _YieldFarmContinuous.Contract.Owner(&_YieldFarmContinuous.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) Owner() (common.Address, error) {
	return _YieldFarmContinuous.Contract.Owner(&_YieldFarmContinuous.CallOpts)
}

// PoolSize is a free data retrieval call binding the contract method 0x4ec18db9.
//
// Solidity: function poolSize() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) PoolSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "poolSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolSize is a free data retrieval call binding the contract method 0x4ec18db9.
//
// Solidity: function poolSize() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) PoolSize() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.PoolSize(&_YieldFarmContinuous.CallOpts)
}

// PoolSize is a free data retrieval call binding the contract method 0x4ec18db9.
//
// Solidity: function poolSize() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) PoolSize() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.PoolSize(&_YieldFarmContinuous.CallOpts)
}

// PoolToken is a free data retrieval call binding the contract method 0xcbdf382c.
//
// Solidity: function poolToken() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) PoolToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "poolToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolToken is a free data retrieval call binding the contract method 0xcbdf382c.
//
// Solidity: function poolToken() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousSession) PoolToken() (common.Address, error) {
	return _YieldFarmContinuous.Contract.PoolToken(&_YieldFarmContinuous.CallOpts)
}

// PoolToken is a free data retrieval call binding the contract method 0xcbdf382c.
//
// Solidity: function poolToken() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) PoolToken() (common.Address, error) {
	return _YieldFarmContinuous.Contract.PoolToken(&_YieldFarmContinuous.CallOpts)
}

// RewardNotTransferred is a free data retrieval call binding the contract method 0x5cb54206.
//
// Solidity: function rewardNotTransferred() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) RewardNotTransferred(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "rewardNotTransferred")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardNotTransferred is a free data retrieval call binding the contract method 0x5cb54206.
//
// Solidity: function rewardNotTransferred() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) RewardNotTransferred() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.RewardNotTransferred(&_YieldFarmContinuous.CallOpts)
}

// RewardNotTransferred is a free data retrieval call binding the contract method 0x5cb54206.
//
// Solidity: function rewardNotTransferred() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) RewardNotTransferred() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.RewardNotTransferred(&_YieldFarmContinuous.CallOpts)
}

// RewardRatePerSecond is a free data retrieval call binding the contract method 0x5d773b2c.
//
// Solidity: function rewardRatePerSecond() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) RewardRatePerSecond(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "rewardRatePerSecond")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardRatePerSecond is a free data retrieval call binding the contract method 0x5d773b2c.
//
// Solidity: function rewardRatePerSecond() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) RewardRatePerSecond() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.RewardRatePerSecond(&_YieldFarmContinuous.CallOpts)
}

// RewardRatePerSecond is a free data retrieval call binding the contract method 0x5d773b2c.
//
// Solidity: function rewardRatePerSecond() view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) RewardRatePerSecond() (*big.Int, error) {
	return _YieldFarmContinuous.Contract.RewardRatePerSecond(&_YieldFarmContinuous.CallOpts)
}

// RewardSource is a free data retrieval call binding the contract method 0x9cfbc002.
//
// Solidity: function rewardSource() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) RewardSource(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "rewardSource")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardSource is a free data retrieval call binding the contract method 0x9cfbc002.
//
// Solidity: function rewardSource() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousSession) RewardSource() (common.Address, error) {
	return _YieldFarmContinuous.Contract.RewardSource(&_YieldFarmContinuous.CallOpts)
}

// RewardSource is a free data retrieval call binding the contract method 0x9cfbc002.
//
// Solidity: function rewardSource() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) RewardSource() (common.Address, error) {
	return _YieldFarmContinuous.Contract.RewardSource(&_YieldFarmContinuous.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) RewardToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "rewardToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousSession) RewardToken() (common.Address, error) {
	return _YieldFarmContinuous.Contract.RewardToken(&_YieldFarmContinuous.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) RewardToken() (common.Address, error) {
	return _YieldFarmContinuous.Contract.RewardToken(&_YieldFarmContinuous.CallOpts)
}

// UserMultiplier is a free data retrieval call binding the contract method 0xb1a03b6b.
//
// Solidity: function userMultiplier(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCaller) UserMultiplier(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _YieldFarmContinuous.contract.Call(opts, &out, "userMultiplier", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserMultiplier is a free data retrieval call binding the contract method 0xb1a03b6b.
//
// Solidity: function userMultiplier(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) UserMultiplier(arg0 common.Address) (*big.Int, error) {
	return _YieldFarmContinuous.Contract.UserMultiplier(&_YieldFarmContinuous.CallOpts, arg0)
}

// UserMultiplier is a free data retrieval call binding the contract method 0xb1a03b6b.
//
// Solidity: function userMultiplier(address ) view returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousCallerSession) UserMultiplier(arg0 common.Address) (*big.Int, error) {
	return _YieldFarmContinuous.Contract.UserMultiplier(&_YieldFarmContinuous.CallOpts, arg0)
}

// AckFunds is a paid mutator transaction binding the contract method 0xacfd9325.
//
// Solidity: function ackFunds() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) AckFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "ackFunds")
}

// AckFunds is a paid mutator transaction binding the contract method 0xacfd9325.
//
// Solidity: function ackFunds() returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) AckFunds() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.AckFunds(&_YieldFarmContinuous.TransactOpts)
}

// AckFunds is a paid mutator transaction binding the contract method 0xacfd9325.
//
// Solidity: function ackFunds() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) AckFunds() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.AckFunds(&_YieldFarmContinuous.TransactOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x4e71d92d.
//
// Solidity: function claim() returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) Claim(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "claim")
}

// Claim is a paid mutator transaction binding the contract method 0x4e71d92d.
//
// Solidity: function claim() returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) Claim() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.Claim(&_YieldFarmContinuous.TransactOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x4e71d92d.
//
// Solidity: function claim() returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) Claim() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.Claim(&_YieldFarmContinuous.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.Deposit(&_YieldFarmContinuous.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.Deposit(&_YieldFarmContinuous.TransactOpts, amount)
}

// PullRewardFromSource is a paid mutator transaction binding the contract method 0x94094adf.
//
// Solidity: function pullRewardFromSource() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) PullRewardFromSource(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "pullRewardFromSource")
}

// PullRewardFromSource is a paid mutator transaction binding the contract method 0x94094adf.
//
// Solidity: function pullRewardFromSource() returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) PullRewardFromSource() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.PullRewardFromSource(&_YieldFarmContinuous.TransactOpts)
}

// PullRewardFromSource is a paid mutator transaction binding the contract method 0x94094adf.
//
// Solidity: function pullRewardFromSource() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) PullRewardFromSource() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.PullRewardFromSource(&_YieldFarmContinuous.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) RenounceOwnership() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.RenounceOwnership(&_YieldFarmContinuous.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.RenounceOwnership(&_YieldFarmContinuous.TransactOpts)
}

// SetRewardRatePerSecond is a paid mutator transaction binding the contract method 0xfdea76b4.
//
// Solidity: function setRewardRatePerSecond(uint256 rate) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) SetRewardRatePerSecond(opts *bind.TransactOpts, rate *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "setRewardRatePerSecond", rate)
}

// SetRewardRatePerSecond is a paid mutator transaction binding the contract method 0xfdea76b4.
//
// Solidity: function setRewardRatePerSecond(uint256 rate) returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) SetRewardRatePerSecond(rate *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.SetRewardRatePerSecond(&_YieldFarmContinuous.TransactOpts, rate)
}

// SetRewardRatePerSecond is a paid mutator transaction binding the contract method 0xfdea76b4.
//
// Solidity: function setRewardRatePerSecond(uint256 rate) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) SetRewardRatePerSecond(rate *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.SetRewardRatePerSecond(&_YieldFarmContinuous.TransactOpts, rate)
}

// SetRewardsSource is a paid mutator transaction binding the contract method 0xd1c76638.
//
// Solidity: function setRewardsSource(address src) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) SetRewardsSource(opts *bind.TransactOpts, src common.Address) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "setRewardsSource", src)
}

// SetRewardsSource is a paid mutator transaction binding the contract method 0xd1c76638.
//
// Solidity: function setRewardsSource(address src) returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) SetRewardsSource(src common.Address) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.SetRewardsSource(&_YieldFarmContinuous.TransactOpts, src)
}

// SetRewardsSource is a paid mutator transaction binding the contract method 0xd1c76638.
//
// Solidity: function setRewardsSource(address src) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) SetRewardsSource(src common.Address) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.SetRewardsSource(&_YieldFarmContinuous.TransactOpts, src)
}

// SoftPullReward is a paid mutator transaction binding the contract method 0x3378186a.
//
// Solidity: function softPullReward() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) SoftPullReward(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "softPullReward")
}

// SoftPullReward is a paid mutator transaction binding the contract method 0x3378186a.
//
// Solidity: function softPullReward() returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) SoftPullReward() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.SoftPullReward(&_YieldFarmContinuous.TransactOpts)
}

// SoftPullReward is a paid mutator transaction binding the contract method 0x3378186a.
//
// Solidity: function softPullReward() returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) SoftPullReward() (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.SoftPullReward(&_YieldFarmContinuous.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.TransferOwnership(&_YieldFarmContinuous.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.TransferOwnership(&_YieldFarmContinuous.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_YieldFarmContinuous *YieldFarmContinuousSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.Withdraw(&_YieldFarmContinuous.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.Withdraw(&_YieldFarmContinuous.TransactOpts, amount)
}

// WithdrawAndClaim is a paid mutator transaction binding the contract method 0x448a1047.
//
// Solidity: function withdrawAndClaim(uint256 amount) returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousTransactor) WithdrawAndClaim(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.contract.Transact(opts, "withdrawAndClaim", amount)
}

// WithdrawAndClaim is a paid mutator transaction binding the contract method 0x448a1047.
//
// Solidity: function withdrawAndClaim(uint256 amount) returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousSession) WithdrawAndClaim(amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.WithdrawAndClaim(&_YieldFarmContinuous.TransactOpts, amount)
}

// WithdrawAndClaim is a paid mutator transaction binding the contract method 0x448a1047.
//
// Solidity: function withdrawAndClaim(uint256 amount) returns(uint256)
func (_YieldFarmContinuous *YieldFarmContinuousTransactorSession) WithdrawAndClaim(amount *big.Int) (*types.Transaction, error) {
	return _YieldFarmContinuous.Contract.WithdrawAndClaim(&_YieldFarmContinuous.TransactOpts, amount)
}

// YieldFarmContinuousClaimIterator is returned from FilterClaim and is used to iterate over the raw logs and unpacked data for Claim events raised by the YieldFarmContinuous contract.
type YieldFarmContinuousClaimIterator struct {
	Event *YieldFarmContinuousClaim // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldFarmContinuousClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldFarmContinuousClaim)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldFarmContinuousClaim)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldFarmContinuousClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldFarmContinuousClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldFarmContinuousClaim represents a Claim event raised by the YieldFarmContinuous contract.
type YieldFarmContinuousClaim struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClaim is a free log retrieval operation binding the contract event 0x47cee97cb7acd717b3c0aa1435d004cd5b3c8c57d70dbceb4e4458bbd60e39d4.
//
// Solidity: event Claim(address indexed user, uint256 amount)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) FilterClaim(opts *bind.FilterOpts, user []common.Address) (*YieldFarmContinuousClaimIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.FilterLogs(opts, "Claim", userRule)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousClaimIterator{contract: _YieldFarmContinuous.contract, event: "Claim", logs: logs, sub: sub}, nil
}

// WatchClaim is a free log subscription operation binding the contract event 0x47cee97cb7acd717b3c0aa1435d004cd5b3c8c57d70dbceb4e4458bbd60e39d4.
//
// Solidity: event Claim(address indexed user, uint256 amount)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) WatchClaim(opts *bind.WatchOpts, sink chan<- *YieldFarmContinuousClaim, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.WatchLogs(opts, "Claim", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldFarmContinuousClaim)
				if err := _YieldFarmContinuous.contract.UnpackLog(event, "Claim", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaim is a log parse operation binding the contract event 0x47cee97cb7acd717b3c0aa1435d004cd5b3c8c57d70dbceb4e4458bbd60e39d4.
//
// Solidity: event Claim(address indexed user, uint256 amount)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) ParseClaim(log types.Log) (*YieldFarmContinuousClaim, error) {
	event := new(YieldFarmContinuousClaim)
	if err := _YieldFarmContinuous.contract.UnpackLog(event, "Claim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldFarmContinuousDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the YieldFarmContinuous contract.
type YieldFarmContinuousDepositIterator struct {
	Event *YieldFarmContinuousDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldFarmContinuousDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldFarmContinuousDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldFarmContinuousDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldFarmContinuousDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldFarmContinuousDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldFarmContinuousDeposit represents a Deposit event raised by the YieldFarmContinuous contract.
type YieldFarmContinuousDeposit struct {
	User         common.Address
	Amount       *big.Int
	BalanceAfter *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x90890809c654f11d6e72a28fa60149770a0d11ec6c92319d6ceb2bb0a4ea1a15.
//
// Solidity: event Deposit(address indexed user, uint256 amount, uint256 balanceAfter)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address) (*YieldFarmContinuousDepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.FilterLogs(opts, "Deposit", userRule)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousDepositIterator{contract: _YieldFarmContinuous.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x90890809c654f11d6e72a28fa60149770a0d11ec6c92319d6ceb2bb0a4ea1a15.
//
// Solidity: event Deposit(address indexed user, uint256 amount, uint256 balanceAfter)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *YieldFarmContinuousDeposit, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.WatchLogs(opts, "Deposit", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldFarmContinuousDeposit)
				if err := _YieldFarmContinuous.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposit is a log parse operation binding the contract event 0x90890809c654f11d6e72a28fa60149770a0d11ec6c92319d6ceb2bb0a4ea1a15.
//
// Solidity: event Deposit(address indexed user, uint256 amount, uint256 balanceAfter)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) ParseDeposit(log types.Log) (*YieldFarmContinuousDeposit, error) {
	event := new(YieldFarmContinuousDeposit)
	if err := _YieldFarmContinuous.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldFarmContinuousOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the YieldFarmContinuous contract.
type YieldFarmContinuousOwnershipTransferredIterator struct {
	Event *YieldFarmContinuousOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldFarmContinuousOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldFarmContinuousOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldFarmContinuousOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldFarmContinuousOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldFarmContinuousOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldFarmContinuousOwnershipTransferred represents a OwnershipTransferred event raised by the YieldFarmContinuous contract.
type YieldFarmContinuousOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*YieldFarmContinuousOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousOwnershipTransferredIterator{contract: _YieldFarmContinuous.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *YieldFarmContinuousOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldFarmContinuousOwnershipTransferred)
				if err := _YieldFarmContinuous.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) ParseOwnershipTransferred(log types.Log) (*YieldFarmContinuousOwnershipTransferred, error) {
	event := new(YieldFarmContinuousOwnershipTransferred)
	if err := _YieldFarmContinuous.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldFarmContinuousWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the YieldFarmContinuous contract.
type YieldFarmContinuousWithdrawIterator struct {
	Event *YieldFarmContinuousWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldFarmContinuousWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldFarmContinuousWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldFarmContinuousWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldFarmContinuousWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldFarmContinuousWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldFarmContinuousWithdraw represents a Withdraw event raised by the YieldFarmContinuous contract.
type YieldFarmContinuousWithdraw struct {
	User         common.Address
	Amount       *big.Int
	BalanceAfter *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed user, uint256 amount, uint256 balanceAfter)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address) (*YieldFarmContinuousWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.FilterLogs(opts, "Withdraw", userRule)
	if err != nil {
		return nil, err
	}
	return &YieldFarmContinuousWithdrawIterator{contract: _YieldFarmContinuous.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed user, uint256 amount, uint256 balanceAfter)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *YieldFarmContinuousWithdraw, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _YieldFarmContinuous.contract.WatchLogs(opts, "Withdraw", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldFarmContinuousWithdraw)
				if err := _YieldFarmContinuous.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdraw is a log parse operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed user, uint256 amount, uint256 balanceAfter)
func (_YieldFarmContinuous *YieldFarmContinuousFilterer) ParseWithdraw(log types.Log) (*YieldFarmContinuousWithdraw, error) {
	event := new(YieldFarmContinuousWithdraw)
	if err := _YieldFarmContinuous.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
