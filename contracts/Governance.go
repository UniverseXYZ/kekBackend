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

// GovernanceReceipt is an auto generated low-level Go binding around an user-defined struct.
type GovernanceReceipt struct {
	HasVoted bool
	Votes    *big.Int
	Support  bool
}

// GovernanceABI is the input ABI used to generate the binding from.
const GovernanceABI = "[{\"inputs\":[],\"name\":\"ACTIVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GRACE_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINIMUM_FOR_VOTES_THRESHOLD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINIMUM_QUORUM\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"QUEUE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WARM_UP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"abdicate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newGuardian\",\"type\":\"address\"}],\"name\":\"anoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancelVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"}],\"name\":\"castVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"getActions\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"getReceipt\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"hasVoted\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"}],\"internalType\":\"structGovernance.Receipt\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"guardian\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"barnAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"guardianAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastProposalId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"latestProposalIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposalMaxOperations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"createTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quorum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"forVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"againstVotes\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"canceled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"queue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"queuedTransactions\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setActivePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setGracePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"quorum\",\"type\":\"uint256\"}],\"name\":\"setMinimumQuorum\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setMinimumThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setQueuePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setWarmUpPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"startVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"state\",\"outputs\":[{\"internalType\":\"enumGovernance.ProposalState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Governance is an auto generated Go binding around an Ethereum contract.
type Governance struct {
	GovernanceCaller     // Read-only binding to the contract
	GovernanceTransactor // Write-only binding to the contract
	GovernanceFilterer   // Log filterer for contract events
}

// GovernanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type GovernanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovernanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GovernanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovernanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GovernanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovernanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GovernanceSession struct {
	Contract     *Governance       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GovernanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GovernanceCallerSession struct {
	Contract *GovernanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// GovernanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GovernanceTransactorSession struct {
	Contract     *GovernanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// GovernanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type GovernanceRaw struct {
	Contract *Governance // Generic contract binding to access the raw methods on
}

// GovernanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GovernanceCallerRaw struct {
	Contract *GovernanceCaller // Generic read-only contract binding to access the raw methods on
}

// GovernanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GovernanceTransactorRaw struct {
	Contract *GovernanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGovernance creates a new instance of Governance, bound to a specific deployed contract.
func NewGovernance(address common.Address, backend bind.ContractBackend) (*Governance, error) {
	contract, err := bindGovernance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Governance{GovernanceCaller: GovernanceCaller{contract: contract}, GovernanceTransactor: GovernanceTransactor{contract: contract}, GovernanceFilterer: GovernanceFilterer{contract: contract}}, nil
}

// NewGovernanceCaller creates a new read-only instance of Governance, bound to a specific deployed contract.
func NewGovernanceCaller(address common.Address, caller bind.ContractCaller) (*GovernanceCaller, error) {
	contract, err := bindGovernance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GovernanceCaller{contract: contract}, nil
}

// NewGovernanceTransactor creates a new write-only instance of Governance, bound to a specific deployed contract.
func NewGovernanceTransactor(address common.Address, transactor bind.ContractTransactor) (*GovernanceTransactor, error) {
	contract, err := bindGovernance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GovernanceTransactor{contract: contract}, nil
}

// NewGovernanceFilterer creates a new log filterer instance of Governance, bound to a specific deployed contract.
func NewGovernanceFilterer(address common.Address, filterer bind.ContractFilterer) (*GovernanceFilterer, error) {
	contract, err := bindGovernance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GovernanceFilterer{contract: contract}, nil
}

// bindGovernance binds a generic wrapper to an already deployed contract.
func bindGovernance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GovernanceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Governance *GovernanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Governance.Contract.GovernanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Governance *GovernanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Governance.Contract.GovernanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Governance *GovernanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Governance.Contract.GovernanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Governance *GovernanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Governance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Governance *GovernanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Governance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Governance *GovernanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Governance.Contract.contract.Transact(opts, method, params...)
}

// ACTIVE is a free data retrieval call binding the contract method 0xc90bd047.
//
// Solidity: function ACTIVE() view returns(uint256)
func (_Governance *GovernanceCaller) ACTIVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "ACTIVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ACTIVE is a free data retrieval call binding the contract method 0xc90bd047.
//
// Solidity: function ACTIVE() view returns(uint256)
func (_Governance *GovernanceSession) ACTIVE() (*big.Int, error) {
	return _Governance.Contract.ACTIVE(&_Governance.CallOpts)
}

// ACTIVE is a free data retrieval call binding the contract method 0xc90bd047.
//
// Solidity: function ACTIVE() view returns(uint256)
func (_Governance *GovernanceCallerSession) ACTIVE() (*big.Int, error) {
	return _Governance.Contract.ACTIVE(&_Governance.CallOpts)
}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint256)
func (_Governance *GovernanceCaller) GRACEPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "GRACE_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint256)
func (_Governance *GovernanceSession) GRACEPERIOD() (*big.Int, error) {
	return _Governance.Contract.GRACEPERIOD(&_Governance.CallOpts)
}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint256)
func (_Governance *GovernanceCallerSession) GRACEPERIOD() (*big.Int, error) {
	return _Governance.Contract.GRACEPERIOD(&_Governance.CallOpts)
}

// MINIMUMFORVOTESTHRESHOLD is a free data retrieval call binding the contract method 0xb17ac1e1.
//
// Solidity: function MINIMUM_FOR_VOTES_THRESHOLD() view returns(uint256)
func (_Governance *GovernanceCaller) MINIMUMFORVOTESTHRESHOLD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "MINIMUM_FOR_VOTES_THRESHOLD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINIMUMFORVOTESTHRESHOLD is a free data retrieval call binding the contract method 0xb17ac1e1.
//
// Solidity: function MINIMUM_FOR_VOTES_THRESHOLD() view returns(uint256)
func (_Governance *GovernanceSession) MINIMUMFORVOTESTHRESHOLD() (*big.Int, error) {
	return _Governance.Contract.MINIMUMFORVOTESTHRESHOLD(&_Governance.CallOpts)
}

// MINIMUMFORVOTESTHRESHOLD is a free data retrieval call binding the contract method 0xb17ac1e1.
//
// Solidity: function MINIMUM_FOR_VOTES_THRESHOLD() view returns(uint256)
func (_Governance *GovernanceCallerSession) MINIMUMFORVOTESTHRESHOLD() (*big.Int, error) {
	return _Governance.Contract.MINIMUMFORVOTESTHRESHOLD(&_Governance.CallOpts)
}

// MINIMUMQUORUM is a free data retrieval call binding the contract method 0xb159beac.
//
// Solidity: function MINIMUM_QUORUM() view returns(uint256)
func (_Governance *GovernanceCaller) MINIMUMQUORUM(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "MINIMUM_QUORUM")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINIMUMQUORUM is a free data retrieval call binding the contract method 0xb159beac.
//
// Solidity: function MINIMUM_QUORUM() view returns(uint256)
func (_Governance *GovernanceSession) MINIMUMQUORUM() (*big.Int, error) {
	return _Governance.Contract.MINIMUMQUORUM(&_Governance.CallOpts)
}

// MINIMUMQUORUM is a free data retrieval call binding the contract method 0xb159beac.
//
// Solidity: function MINIMUM_QUORUM() view returns(uint256)
func (_Governance *GovernanceCallerSession) MINIMUMQUORUM() (*big.Int, error) {
	return _Governance.Contract.MINIMUMQUORUM(&_Governance.CallOpts)
}

// QUEUE is a free data retrieval call binding the contract method 0x1a3d203d.
//
// Solidity: function QUEUE() view returns(uint256)
func (_Governance *GovernanceCaller) QUEUE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "QUEUE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QUEUE is a free data retrieval call binding the contract method 0x1a3d203d.
//
// Solidity: function QUEUE() view returns(uint256)
func (_Governance *GovernanceSession) QUEUE() (*big.Int, error) {
	return _Governance.Contract.QUEUE(&_Governance.CallOpts)
}

// QUEUE is a free data retrieval call binding the contract method 0x1a3d203d.
//
// Solidity: function QUEUE() view returns(uint256)
func (_Governance *GovernanceCallerSession) QUEUE() (*big.Int, error) {
	return _Governance.Contract.QUEUE(&_Governance.CallOpts)
}

// WARMUP is a free data retrieval call binding the contract method 0x1e9d5fad.
//
// Solidity: function WARM_UP() view returns(uint256)
func (_Governance *GovernanceCaller) WARMUP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "WARM_UP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WARMUP is a free data retrieval call binding the contract method 0x1e9d5fad.
//
// Solidity: function WARM_UP() view returns(uint256)
func (_Governance *GovernanceSession) WARMUP() (*big.Int, error) {
	return _Governance.Contract.WARMUP(&_Governance.CallOpts)
}

// WARMUP is a free data retrieval call binding the contract method 0x1e9d5fad.
//
// Solidity: function WARM_UP() view returns(uint256)
func (_Governance *GovernanceCallerSession) WARMUP() (*big.Int, error) {
	return _Governance.Contract.WARMUP(&_Governance.CallOpts)
}

// GetActions is a free data retrieval call binding the contract method 0x328dd982.
//
// Solidity: function getActions(uint256 proposalId) view returns(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas)
func (_Governance *GovernanceCaller) GetActions(opts *bind.CallOpts, proposalId *big.Int) (struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "getActions", proposalId)

	outstruct := new(struct {
		Targets    []common.Address
		Values     []*big.Int
		Signatures []string
		Calldatas  [][]byte
	})

	outstruct.Targets = out[0].([]common.Address)
	outstruct.Values = out[1].([]*big.Int)
	outstruct.Signatures = out[2].([]string)
	outstruct.Calldatas = out[3].([][]byte)

	return *outstruct, err

}

// GetActions is a free data retrieval call binding the contract method 0x328dd982.
//
// Solidity: function getActions(uint256 proposalId) view returns(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas)
func (_Governance *GovernanceSession) GetActions(proposalId *big.Int) (struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}, error) {
	return _Governance.Contract.GetActions(&_Governance.CallOpts, proposalId)
}

// GetActions is a free data retrieval call binding the contract method 0x328dd982.
//
// Solidity: function getActions(uint256 proposalId) view returns(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas)
func (_Governance *GovernanceCallerSession) GetActions(proposalId *big.Int) (struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}, error) {
	return _Governance.Contract.GetActions(&_Governance.CallOpts, proposalId)
}

// GetReceipt is a free data retrieval call binding the contract method 0xe23a9a52.
//
// Solidity: function getReceipt(uint256 proposalId, address voter) view returns((bool,uint256,bool))
func (_Governance *GovernanceCaller) GetReceipt(opts *bind.CallOpts, proposalId *big.Int, voter common.Address) (GovernanceReceipt, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "getReceipt", proposalId, voter)

	if err != nil {
		return *new(GovernanceReceipt), err
	}

	out0 := *abi.ConvertType(out[0], new(GovernanceReceipt)).(*GovernanceReceipt)

	return out0, err

}

// GetReceipt is a free data retrieval call binding the contract method 0xe23a9a52.
//
// Solidity: function getReceipt(uint256 proposalId, address voter) view returns((bool,uint256,bool))
func (_Governance *GovernanceSession) GetReceipt(proposalId *big.Int, voter common.Address) (GovernanceReceipt, error) {
	return _Governance.Contract.GetReceipt(&_Governance.CallOpts, proposalId, voter)
}

// GetReceipt is a free data retrieval call binding the contract method 0xe23a9a52.
//
// Solidity: function getReceipt(uint256 proposalId, address voter) view returns((bool,uint256,bool))
func (_Governance *GovernanceCallerSession) GetReceipt(proposalId *big.Int, voter common.Address) (GovernanceReceipt, error) {
	return _Governance.Contract.GetReceipt(&_Governance.CallOpts, proposalId, voter)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_Governance *GovernanceCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_Governance *GovernanceSession) Guardian() (common.Address, error) {
	return _Governance.Contract.Guardian(&_Governance.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_Governance *GovernanceCallerSession) Guardian() (common.Address, error) {
	return _Governance.Contract.Guardian(&_Governance.CallOpts)
}

// LastProposalId is a free data retrieval call binding the contract method 0x74cb3041.
//
// Solidity: function lastProposalId() view returns(uint256)
func (_Governance *GovernanceCaller) LastProposalId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "lastProposalId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastProposalId is a free data retrieval call binding the contract method 0x74cb3041.
//
// Solidity: function lastProposalId() view returns(uint256)
func (_Governance *GovernanceSession) LastProposalId() (*big.Int, error) {
	return _Governance.Contract.LastProposalId(&_Governance.CallOpts)
}

// LastProposalId is a free data retrieval call binding the contract method 0x74cb3041.
//
// Solidity: function lastProposalId() view returns(uint256)
func (_Governance *GovernanceCallerSession) LastProposalId() (*big.Int, error) {
	return _Governance.Contract.LastProposalId(&_Governance.CallOpts)
}

// LatestProposalIds is a free data retrieval call binding the contract method 0x17977c61.
//
// Solidity: function latestProposalIds(address ) view returns(uint256)
func (_Governance *GovernanceCaller) LatestProposalIds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "latestProposalIds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestProposalIds is a free data retrieval call binding the contract method 0x17977c61.
//
// Solidity: function latestProposalIds(address ) view returns(uint256)
func (_Governance *GovernanceSession) LatestProposalIds(arg0 common.Address) (*big.Int, error) {
	return _Governance.Contract.LatestProposalIds(&_Governance.CallOpts, arg0)
}

// LatestProposalIds is a free data retrieval call binding the contract method 0x17977c61.
//
// Solidity: function latestProposalIds(address ) view returns(uint256)
func (_Governance *GovernanceCallerSession) LatestProposalIds(arg0 common.Address) (*big.Int, error) {
	return _Governance.Contract.LatestProposalIds(&_Governance.CallOpts, arg0)
}

// ProposalMaxOperations is a free data retrieval call binding the contract method 0x7bdbe4d0.
//
// Solidity: function proposalMaxOperations() pure returns(uint256)
func (_Governance *GovernanceCaller) ProposalMaxOperations(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "proposalMaxOperations")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalMaxOperations is a free data retrieval call binding the contract method 0x7bdbe4d0.
//
// Solidity: function proposalMaxOperations() pure returns(uint256)
func (_Governance *GovernanceSession) ProposalMaxOperations() (*big.Int, error) {
	return _Governance.Contract.ProposalMaxOperations(&_Governance.CallOpts)
}

// ProposalMaxOperations is a free data retrieval call binding the contract method 0x7bdbe4d0.
//
// Solidity: function proposalMaxOperations() pure returns(uint256)
func (_Governance *GovernanceCallerSession) ProposalMaxOperations() (*big.Int, error) {
	return _Governance.Contract.ProposalMaxOperations(&_Governance.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) view returns(uint256 id, address proposer, string description, string title, uint256 createTime, uint256 startTime, uint256 quorum, uint256 eta, uint256 forVotes, uint256 againstVotes, bool canceled, bool executed)
func (_Governance *GovernanceCaller) Proposals(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
	StartTime    *big.Int
	Quorum       *big.Int
	Eta          *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "proposals", arg0)

	outstruct := new(struct {
		Id           *big.Int
		Proposer     common.Address
		Description  string
		Title        string
		CreateTime   *big.Int
		StartTime    *big.Int
		Quorum       *big.Int
		Eta          *big.Int
		ForVotes     *big.Int
		AgainstVotes *big.Int
		Canceled     bool
		Executed     bool
	})

	outstruct.Id = out[0].(*big.Int)
	outstruct.Proposer = out[1].(common.Address)
	outstruct.Description = out[2].(string)
	outstruct.Title = out[3].(string)
	outstruct.CreateTime = out[4].(*big.Int)
	outstruct.StartTime = out[5].(*big.Int)
	outstruct.Quorum = out[6].(*big.Int)
	outstruct.Eta = out[7].(*big.Int)
	outstruct.ForVotes = out[8].(*big.Int)
	outstruct.AgainstVotes = out[9].(*big.Int)
	outstruct.Canceled = out[10].(bool)
	outstruct.Executed = out[11].(bool)

	return *outstruct, err

}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) view returns(uint256 id, address proposer, string description, string title, uint256 createTime, uint256 startTime, uint256 quorum, uint256 eta, uint256 forVotes, uint256 againstVotes, bool canceled, bool executed)
func (_Governance *GovernanceSession) Proposals(arg0 *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
	StartTime    *big.Int
	Quorum       *big.Int
	Eta          *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	return _Governance.Contract.Proposals(&_Governance.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) view returns(uint256 id, address proposer, string description, string title, uint256 createTime, uint256 startTime, uint256 quorum, uint256 eta, uint256 forVotes, uint256 againstVotes, bool canceled, bool executed)
func (_Governance *GovernanceCallerSession) Proposals(arg0 *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
	StartTime    *big.Int
	Quorum       *big.Int
	Eta          *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	return _Governance.Contract.Proposals(&_Governance.CallOpts, arg0)
}

// QueuedTransactions is a free data retrieval call binding the contract method 0xf2b06537.
//
// Solidity: function queuedTransactions(bytes32 ) view returns(bool)
func (_Governance *GovernanceCaller) QueuedTransactions(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "queuedTransactions", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// QueuedTransactions is a free data retrieval call binding the contract method 0xf2b06537.
//
// Solidity: function queuedTransactions(bytes32 ) view returns(bool)
func (_Governance *GovernanceSession) QueuedTransactions(arg0 [32]byte) (bool, error) {
	return _Governance.Contract.QueuedTransactions(&_Governance.CallOpts, arg0)
}

// QueuedTransactions is a free data retrieval call binding the contract method 0xf2b06537.
//
// Solidity: function queuedTransactions(bytes32 ) view returns(bool)
func (_Governance *GovernanceCallerSession) QueuedTransactions(arg0 [32]byte) (bool, error) {
	return _Governance.Contract.QueuedTransactions(&_Governance.CallOpts, arg0)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_Governance *GovernanceCaller) State(opts *bind.CallOpts, proposalId *big.Int) (uint8, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "state", proposalId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_Governance *GovernanceSession) State(proposalId *big.Int) (uint8, error) {
	return _Governance.Contract.State(&_Governance.CallOpts, proposalId)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_Governance *GovernanceCallerSession) State(proposalId *big.Int) (uint8, error) {
	return _Governance.Contract.State(&_Governance.CallOpts, proposalId)
}

// Abdicate is a paid mutator transaction binding the contract method 0x314e99a2.
//
// Solidity: function abdicate() returns()
func (_Governance *GovernanceTransactor) Abdicate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "abdicate")
}

// Abdicate is a paid mutator transaction binding the contract method 0x314e99a2.
//
// Solidity: function abdicate() returns()
func (_Governance *GovernanceSession) Abdicate() (*types.Transaction, error) {
	return _Governance.Contract.Abdicate(&_Governance.TransactOpts)
}

// Abdicate is a paid mutator transaction binding the contract method 0x314e99a2.
//
// Solidity: function abdicate() returns()
func (_Governance *GovernanceTransactorSession) Abdicate() (*types.Transaction, error) {
	return _Governance.Contract.Abdicate(&_Governance.TransactOpts)
}

// Anoint is a paid mutator transaction binding the contract method 0x3addc3b0.
//
// Solidity: function anoint(address newGuardian) returns()
func (_Governance *GovernanceTransactor) Anoint(opts *bind.TransactOpts, newGuardian common.Address) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "anoint", newGuardian)
}

// Anoint is a paid mutator transaction binding the contract method 0x3addc3b0.
//
// Solidity: function anoint(address newGuardian) returns()
func (_Governance *GovernanceSession) Anoint(newGuardian common.Address) (*types.Transaction, error) {
	return _Governance.Contract.Anoint(&_Governance.TransactOpts, newGuardian)
}

// Anoint is a paid mutator transaction binding the contract method 0x3addc3b0.
//
// Solidity: function anoint(address newGuardian) returns()
func (_Governance *GovernanceTransactorSession) Anoint(newGuardian common.Address) (*types.Transaction, error) {
	return _Governance.Contract.Anoint(&_Governance.TransactOpts, newGuardian)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) Cancel(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "cancel", proposalId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 proposalId) returns()
func (_Governance *GovernanceSession) Cancel(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.Cancel(&_Governance.TransactOpts, proposalId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) Cancel(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.Cancel(&_Governance.TransactOpts, proposalId)
}

// CancelVote is a paid mutator transaction binding the contract method 0xbacbe2da.
//
// Solidity: function cancelVote(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) CancelVote(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "cancelVote", proposalId)
}

// CancelVote is a paid mutator transaction binding the contract method 0xbacbe2da.
//
// Solidity: function cancelVote(uint256 proposalId) returns()
func (_Governance *GovernanceSession) CancelVote(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.CancelVote(&_Governance.TransactOpts, proposalId)
}

// CancelVote is a paid mutator transaction binding the contract method 0xbacbe2da.
//
// Solidity: function cancelVote(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) CancelVote(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.CancelVote(&_Governance.TransactOpts, proposalId)
}

// CastVote is a paid mutator transaction binding the contract method 0x15373e3d.
//
// Solidity: function castVote(uint256 proposalId, bool support) returns()
func (_Governance *GovernanceTransactor) CastVote(opts *bind.TransactOpts, proposalId *big.Int, support bool) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "castVote", proposalId, support)
}

// CastVote is a paid mutator transaction binding the contract method 0x15373e3d.
//
// Solidity: function castVote(uint256 proposalId, bool support) returns()
func (_Governance *GovernanceSession) CastVote(proposalId *big.Int, support bool) (*types.Transaction, error) {
	return _Governance.Contract.CastVote(&_Governance.TransactOpts, proposalId, support)
}

// CastVote is a paid mutator transaction binding the contract method 0x15373e3d.
//
// Solidity: function castVote(uint256 proposalId, bool support) returns()
func (_Governance *GovernanceTransactorSession) CastVote(proposalId *big.Int, support bool) (*types.Transaction, error) {
	return _Governance.Contract.CastVote(&_Governance.TransactOpts, proposalId, support)
}

// Execute is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 proposalId) payable returns()
func (_Governance *GovernanceTransactor) Execute(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "execute", proposalId)
}

// Execute is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 proposalId) payable returns()
func (_Governance *GovernanceSession) Execute(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.Execute(&_Governance.TransactOpts, proposalId)
}

// Execute is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 proposalId) payable returns()
func (_Governance *GovernanceTransactorSession) Execute(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.Execute(&_Governance.TransactOpts, proposalId)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address barnAddr, address guardianAddress) returns()
func (_Governance *GovernanceTransactor) Initialize(opts *bind.TransactOpts, barnAddr common.Address, guardianAddress common.Address) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "initialize", barnAddr, guardianAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address barnAddr, address guardianAddress) returns()
func (_Governance *GovernanceSession) Initialize(barnAddr common.Address, guardianAddress common.Address) (*types.Transaction, error) {
	return _Governance.Contract.Initialize(&_Governance.TransactOpts, barnAddr, guardianAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address barnAddr, address guardianAddress) returns()
func (_Governance *GovernanceTransactorSession) Initialize(barnAddr common.Address, guardianAddress common.Address) (*types.Transaction, error) {
	return _Governance.Contract.Initialize(&_Governance.TransactOpts, barnAddr, guardianAddress)
}

// Propose is a paid mutator transaction binding the contract method 0x490145c8.
//
// Solidity: function propose(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, string description, string title) returns(uint256)
func (_Governance *GovernanceTransactor) Propose(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, signatures []string, calldatas [][]byte, description string, title string) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "propose", targets, values, signatures, calldatas, description, title)
}

// Propose is a paid mutator transaction binding the contract method 0x490145c8.
//
// Solidity: function propose(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, string description, string title) returns(uint256)
func (_Governance *GovernanceSession) Propose(targets []common.Address, values []*big.Int, signatures []string, calldatas [][]byte, description string, title string) (*types.Transaction, error) {
	return _Governance.Contract.Propose(&_Governance.TransactOpts, targets, values, signatures, calldatas, description, title)
}

// Propose is a paid mutator transaction binding the contract method 0x490145c8.
//
// Solidity: function propose(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, string description, string title) returns(uint256)
func (_Governance *GovernanceTransactorSession) Propose(targets []common.Address, values []*big.Int, signatures []string, calldatas [][]byte, description string, title string) (*types.Transaction, error) {
	return _Governance.Contract.Propose(&_Governance.TransactOpts, targets, values, signatures, calldatas, description, title)
}

// Queue is a paid mutator transaction binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) Queue(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "queue", proposalId)
}

// Queue is a paid mutator transaction binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 proposalId) returns()
func (_Governance *GovernanceSession) Queue(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.Queue(&_Governance.TransactOpts, proposalId)
}

// Queue is a paid mutator transaction binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) Queue(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.Queue(&_Governance.TransactOpts, proposalId)
}

// SetActivePeriod is a paid mutator transaction binding the contract method 0xcd820edc.
//
// Solidity: function setActivePeriod(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetActivePeriod(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setActivePeriod", period)
}

// SetActivePeriod is a paid mutator transaction binding the contract method 0xcd820edc.
//
// Solidity: function setActivePeriod(uint256 period) returns()
func (_Governance *GovernanceSession) SetActivePeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetActivePeriod(&_Governance.TransactOpts, period)
}

// SetActivePeriod is a paid mutator transaction binding the contract method 0xcd820edc.
//
// Solidity: function setActivePeriod(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetActivePeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetActivePeriod(&_Governance.TransactOpts, period)
}

// SetGracePeriod is a paid mutator transaction binding the contract method 0xf2f65960.
//
// Solidity: function setGracePeriod(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetGracePeriod(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setGracePeriod", period)
}

// SetGracePeriod is a paid mutator transaction binding the contract method 0xf2f65960.
//
// Solidity: function setGracePeriod(uint256 period) returns()
func (_Governance *GovernanceSession) SetGracePeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetGracePeriod(&_Governance.TransactOpts, period)
}

// SetGracePeriod is a paid mutator transaction binding the contract method 0xf2f65960.
//
// Solidity: function setGracePeriod(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetGracePeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetGracePeriod(&_Governance.TransactOpts, period)
}

// SetMinimumQuorum is a paid mutator transaction binding the contract method 0x6f698fb5.
//
// Solidity: function setMinimumQuorum(uint256 quorum) returns()
func (_Governance *GovernanceTransactor) SetMinimumQuorum(opts *bind.TransactOpts, quorum *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setMinimumQuorum", quorum)
}

// SetMinimumQuorum is a paid mutator transaction binding the contract method 0x6f698fb5.
//
// Solidity: function setMinimumQuorum(uint256 quorum) returns()
func (_Governance *GovernanceSession) SetMinimumQuorum(quorum *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetMinimumQuorum(&_Governance.TransactOpts, quorum)
}

// SetMinimumQuorum is a paid mutator transaction binding the contract method 0x6f698fb5.
//
// Solidity: function setMinimumQuorum(uint256 quorum) returns()
func (_Governance *GovernanceTransactorSession) SetMinimumQuorum(quorum *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetMinimumQuorum(&_Governance.TransactOpts, quorum)
}

// SetMinimumThreshold is a paid mutator transaction binding the contract method 0x67058d29.
//
// Solidity: function setMinimumThreshold(uint256 threshold) returns()
func (_Governance *GovernanceTransactor) SetMinimumThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setMinimumThreshold", threshold)
}

// SetMinimumThreshold is a paid mutator transaction binding the contract method 0x67058d29.
//
// Solidity: function setMinimumThreshold(uint256 threshold) returns()
func (_Governance *GovernanceSession) SetMinimumThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetMinimumThreshold(&_Governance.TransactOpts, threshold)
}

// SetMinimumThreshold is a paid mutator transaction binding the contract method 0x67058d29.
//
// Solidity: function setMinimumThreshold(uint256 threshold) returns()
func (_Governance *GovernanceTransactorSession) SetMinimumThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetMinimumThreshold(&_Governance.TransactOpts, threshold)
}

// SetQueuePeriod is a paid mutator transaction binding the contract method 0x63bd97c3.
//
// Solidity: function setQueuePeriod(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetQueuePeriod(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setQueuePeriod", period)
}

// SetQueuePeriod is a paid mutator transaction binding the contract method 0x63bd97c3.
//
// Solidity: function setQueuePeriod(uint256 period) returns()
func (_Governance *GovernanceSession) SetQueuePeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetQueuePeriod(&_Governance.TransactOpts, period)
}

// SetQueuePeriod is a paid mutator transaction binding the contract method 0x63bd97c3.
//
// Solidity: function setQueuePeriod(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetQueuePeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetQueuePeriod(&_Governance.TransactOpts, period)
}

// SetWarmUpPeriod is a paid mutator transaction binding the contract method 0xf9a571a2.
//
// Solidity: function setWarmUpPeriod(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetWarmUpPeriod(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setWarmUpPeriod", period)
}

// SetWarmUpPeriod is a paid mutator transaction binding the contract method 0xf9a571a2.
//
// Solidity: function setWarmUpPeriod(uint256 period) returns()
func (_Governance *GovernanceSession) SetWarmUpPeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetWarmUpPeriod(&_Governance.TransactOpts, period)
}

// SetWarmUpPeriod is a paid mutator transaction binding the contract method 0xf9a571a2.
//
// Solidity: function setWarmUpPeriod(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetWarmUpPeriod(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetWarmUpPeriod(&_Governance.TransactOpts, period)
}

// StartVote is a paid mutator transaction binding the contract method 0x3f220524.
//
// Solidity: function startVote(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) StartVote(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "startVote", proposalId)
}

// StartVote is a paid mutator transaction binding the contract method 0x3f220524.
//
// Solidity: function startVote(uint256 proposalId) returns()
func (_Governance *GovernanceSession) StartVote(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.StartVote(&_Governance.TransactOpts, proposalId)
}

// StartVote is a paid mutator transaction binding the contract method 0x3f220524.
//
// Solidity: function startVote(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) StartVote(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.StartVote(&_Governance.TransactOpts, proposalId)
}
