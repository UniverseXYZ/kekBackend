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
const GovernanceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"CancellationProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"CancellationProposalStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"power\",\"type\":\"uint256\"}],\"name\":\"CancellationProposalVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"CancellationProposalVoteCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalQueued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"power\",\"type\":\"uint256\"}],\"name\":\"Vote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"VoteCanceled\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptanceThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activeDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancelProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancelVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancelVoteCancellationProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"cancellationProposals\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"createTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"forVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"againstVotes\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"}],\"name\":\"castVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"executeCancellationProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"getActions\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"getCancellationProposalReceipt\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"hasVoted\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"}],\"internalType\":\"structGovernance.Receipt\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"getProposalQuorum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"getReceipt\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"hasVoted\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"}],\"internalType\":\"structGovernance.Receipt\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gracePeriodDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"barnAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastProposalId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"latestProposalIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minQuorum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"createTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"forVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"againstVotes\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"canceled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"queue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queueDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"queuedTransactions\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setAcceptanceThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setActiveDuration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setGracePeriodDuration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"quorum\",\"type\":\"uint256\"}],\"name\":\"setMinQuorum\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setQueueDuration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"setWarmUpDuration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"startCancellationProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"state\",\"outputs\":[{\"internalType\":\"enumGovernance.ProposalState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"support\",\"type\":\"bool\"}],\"name\":\"voteCancellationProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"warmUpDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

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

// AcceptanceThreshold is a free data retrieval call binding the contract method 0xb0edbb9b.
//
// Solidity: function acceptanceThreshold() view returns(uint256)
func (_Governance *GovernanceCaller) AcceptanceThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "acceptanceThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AcceptanceThreshold is a free data retrieval call binding the contract method 0xb0edbb9b.
//
// Solidity: function acceptanceThreshold() view returns(uint256)
func (_Governance *GovernanceSession) AcceptanceThreshold() (*big.Int, error) {
	return _Governance.Contract.AcceptanceThreshold(&_Governance.CallOpts)
}

// AcceptanceThreshold is a free data retrieval call binding the contract method 0xb0edbb9b.
//
// Solidity: function acceptanceThreshold() view returns(uint256)
func (_Governance *GovernanceCallerSession) AcceptanceThreshold() (*big.Int, error) {
	return _Governance.Contract.AcceptanceThreshold(&_Governance.CallOpts)
}

// ActiveDuration is a free data retrieval call binding the contract method 0x3d05f009.
//
// Solidity: function activeDuration() view returns(uint256)
func (_Governance *GovernanceCaller) ActiveDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "activeDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveDuration is a free data retrieval call binding the contract method 0x3d05f009.
//
// Solidity: function activeDuration() view returns(uint256)
func (_Governance *GovernanceSession) ActiveDuration() (*big.Int, error) {
	return _Governance.Contract.ActiveDuration(&_Governance.CallOpts)
}

// ActiveDuration is a free data retrieval call binding the contract method 0x3d05f009.
//
// Solidity: function activeDuration() view returns(uint256)
func (_Governance *GovernanceCallerSession) ActiveDuration() (*big.Int, error) {
	return _Governance.Contract.ActiveDuration(&_Governance.CallOpts)
}

// CancellationProposals is a free data retrieval call binding the contract method 0x3349b563.
//
// Solidity: function cancellationProposals(uint256 ) view returns(address creator, uint256 createTime, uint256 forVotes, uint256 againstVotes)
func (_Governance *GovernanceCaller) CancellationProposals(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Creator      common.Address
	CreateTime   *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
}, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "cancellationProposals", arg0)

	outstruct := new(struct {
		Creator      common.Address
		CreateTime   *big.Int
		ForVotes     *big.Int
		AgainstVotes *big.Int
	})

	outstruct.Creator = out[0].(common.Address)
	outstruct.CreateTime = out[1].(*big.Int)
	outstruct.ForVotes = out[2].(*big.Int)
	outstruct.AgainstVotes = out[3].(*big.Int)

	return *outstruct, err

}

// CancellationProposals is a free data retrieval call binding the contract method 0x3349b563.
//
// Solidity: function cancellationProposals(uint256 ) view returns(address creator, uint256 createTime, uint256 forVotes, uint256 againstVotes)
func (_Governance *GovernanceSession) CancellationProposals(arg0 *big.Int) (struct {
	Creator      common.Address
	CreateTime   *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
}, error) {
	return _Governance.Contract.CancellationProposals(&_Governance.CallOpts, arg0)
}

// CancellationProposals is a free data retrieval call binding the contract method 0x3349b563.
//
// Solidity: function cancellationProposals(uint256 ) view returns(address creator, uint256 createTime, uint256 forVotes, uint256 againstVotes)
func (_Governance *GovernanceCallerSession) CancellationProposals(arg0 *big.Int) (struct {
	Creator      common.Address
	CreateTime   *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
}, error) {
	return _Governance.Contract.CancellationProposals(&_Governance.CallOpts, arg0)
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

// GetCancellationProposalReceipt is a free data retrieval call binding the contract method 0x454aee55.
//
// Solidity: function getCancellationProposalReceipt(uint256 proposalId, address voter) view returns((bool,uint256,bool))
func (_Governance *GovernanceCaller) GetCancellationProposalReceipt(opts *bind.CallOpts, proposalId *big.Int, voter common.Address) (GovernanceReceipt, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "getCancellationProposalReceipt", proposalId, voter)

	if err != nil {
		return *new(GovernanceReceipt), err
	}

	out0 := *abi.ConvertType(out[0], new(GovernanceReceipt)).(*GovernanceReceipt)

	return out0, err

}

// GetCancellationProposalReceipt is a free data retrieval call binding the contract method 0x454aee55.
//
// Solidity: function getCancellationProposalReceipt(uint256 proposalId, address voter) view returns((bool,uint256,bool))
func (_Governance *GovernanceSession) GetCancellationProposalReceipt(proposalId *big.Int, voter common.Address) (GovernanceReceipt, error) {
	return _Governance.Contract.GetCancellationProposalReceipt(&_Governance.CallOpts, proposalId, voter)
}

// GetCancellationProposalReceipt is a free data retrieval call binding the contract method 0x454aee55.
//
// Solidity: function getCancellationProposalReceipt(uint256 proposalId, address voter) view returns((bool,uint256,bool))
func (_Governance *GovernanceCallerSession) GetCancellationProposalReceipt(proposalId *big.Int, voter common.Address) (GovernanceReceipt, error) {
	return _Governance.Contract.GetCancellationProposalReceipt(&_Governance.CallOpts, proposalId, voter)
}

// GetProposalQuorum is a free data retrieval call binding the contract method 0xd0cd595e.
//
// Solidity: function getProposalQuorum(uint256 proposalId) view returns(uint256)
func (_Governance *GovernanceCaller) GetProposalQuorum(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "getProposalQuorum", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProposalQuorum is a free data retrieval call binding the contract method 0xd0cd595e.
//
// Solidity: function getProposalQuorum(uint256 proposalId) view returns(uint256)
func (_Governance *GovernanceSession) GetProposalQuorum(proposalId *big.Int) (*big.Int, error) {
	return _Governance.Contract.GetProposalQuorum(&_Governance.CallOpts, proposalId)
}

// GetProposalQuorum is a free data retrieval call binding the contract method 0xd0cd595e.
//
// Solidity: function getProposalQuorum(uint256 proposalId) view returns(uint256)
func (_Governance *GovernanceCallerSession) GetProposalQuorum(proposalId *big.Int) (*big.Int, error) {
	return _Governance.Contract.GetProposalQuorum(&_Governance.CallOpts, proposalId)
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

// GracePeriodDuration is a free data retrieval call binding the contract method 0xc099f575.
//
// Solidity: function gracePeriodDuration() view returns(uint256)
func (_Governance *GovernanceCaller) GracePeriodDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "gracePeriodDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GracePeriodDuration is a free data retrieval call binding the contract method 0xc099f575.
//
// Solidity: function gracePeriodDuration() view returns(uint256)
func (_Governance *GovernanceSession) GracePeriodDuration() (*big.Int, error) {
	return _Governance.Contract.GracePeriodDuration(&_Governance.CallOpts)
}

// GracePeriodDuration is a free data retrieval call binding the contract method 0xc099f575.
//
// Solidity: function gracePeriodDuration() view returns(uint256)
func (_Governance *GovernanceCallerSession) GracePeriodDuration() (*big.Int, error) {
	return _Governance.Contract.GracePeriodDuration(&_Governance.CallOpts)
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() view returns(bool)
func (_Governance *GovernanceCaller) IsActive(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "isActive")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() view returns(bool)
func (_Governance *GovernanceSession) IsActive() (bool, error) {
	return _Governance.Contract.IsActive(&_Governance.CallOpts)
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() view returns(bool)
func (_Governance *GovernanceCallerSession) IsActive() (bool, error) {
	return _Governance.Contract.IsActive(&_Governance.CallOpts)
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

// MinQuorum is a free data retrieval call binding the contract method 0xb5a127e5.
//
// Solidity: function minQuorum() view returns(uint256)
func (_Governance *GovernanceCaller) MinQuorum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "minQuorum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinQuorum is a free data retrieval call binding the contract method 0xb5a127e5.
//
// Solidity: function minQuorum() view returns(uint256)
func (_Governance *GovernanceSession) MinQuorum() (*big.Int, error) {
	return _Governance.Contract.MinQuorum(&_Governance.CallOpts)
}

// MinQuorum is a free data retrieval call binding the contract method 0xb5a127e5.
//
// Solidity: function minQuorum() view returns(uint256)
func (_Governance *GovernanceCallerSession) MinQuorum() (*big.Int, error) {
	return _Governance.Contract.MinQuorum(&_Governance.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) view returns(uint256 id, address proposer, string description, string title, uint256 createTime, uint256 eta, uint256 forVotes, uint256 againstVotes, bool canceled, bool executed)
func (_Governance *GovernanceCaller) Proposals(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
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
	outstruct.Eta = out[5].(*big.Int)
	outstruct.ForVotes = out[6].(*big.Int)
	outstruct.AgainstVotes = out[7].(*big.Int)
	outstruct.Canceled = out[8].(bool)
	outstruct.Executed = out[9].(bool)

	return *outstruct, err

}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) view returns(uint256 id, address proposer, string description, string title, uint256 createTime, uint256 eta, uint256 forVotes, uint256 againstVotes, bool canceled, bool executed)
func (_Governance *GovernanceSession) Proposals(arg0 *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
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
// Solidity: function proposals(uint256 ) view returns(uint256 id, address proposer, string description, string title, uint256 createTime, uint256 eta, uint256 forVotes, uint256 againstVotes, bool canceled, bool executed)
func (_Governance *GovernanceCallerSession) Proposals(arg0 *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Description  string
	Title        string
	CreateTime   *big.Int
	Eta          *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	return _Governance.Contract.Proposals(&_Governance.CallOpts, arg0)
}

// QueueDuration is a free data retrieval call binding the contract method 0x2e8e34e1.
//
// Solidity: function queueDuration() view returns(uint256)
func (_Governance *GovernanceCaller) QueueDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "queueDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QueueDuration is a free data retrieval call binding the contract method 0x2e8e34e1.
//
// Solidity: function queueDuration() view returns(uint256)
func (_Governance *GovernanceSession) QueueDuration() (*big.Int, error) {
	return _Governance.Contract.QueueDuration(&_Governance.CallOpts)
}

// QueueDuration is a free data retrieval call binding the contract method 0x2e8e34e1.
//
// Solidity: function queueDuration() view returns(uint256)
func (_Governance *GovernanceCallerSession) QueueDuration() (*big.Int, error) {
	return _Governance.Contract.QueueDuration(&_Governance.CallOpts)
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

// WarmUpDuration is a free data retrieval call binding the contract method 0x5f2e9f60.
//
// Solidity: function warmUpDuration() view returns(uint256)
func (_Governance *GovernanceCaller) WarmUpDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Governance.contract.Call(opts, &out, "warmUpDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WarmUpDuration is a free data retrieval call binding the contract method 0x5f2e9f60.
//
// Solidity: function warmUpDuration() view returns(uint256)
func (_Governance *GovernanceSession) WarmUpDuration() (*big.Int, error) {
	return _Governance.Contract.WarmUpDuration(&_Governance.CallOpts)
}

// WarmUpDuration is a free data retrieval call binding the contract method 0x5f2e9f60.
//
// Solidity: function warmUpDuration() view returns(uint256)
func (_Governance *GovernanceCallerSession) WarmUpDuration() (*big.Int, error) {
	return _Governance.Contract.WarmUpDuration(&_Governance.CallOpts)
}

// Activate is a paid mutator transaction binding the contract method 0x0f15f4c0.
//
// Solidity: function activate() returns()
func (_Governance *GovernanceTransactor) Activate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "activate")
}

// Activate is a paid mutator transaction binding the contract method 0x0f15f4c0.
//
// Solidity: function activate() returns()
func (_Governance *GovernanceSession) Activate() (*types.Transaction, error) {
	return _Governance.Contract.Activate(&_Governance.TransactOpts)
}

// Activate is a paid mutator transaction binding the contract method 0x0f15f4c0.
//
// Solidity: function activate() returns()
func (_Governance *GovernanceTransactorSession) Activate() (*types.Transaction, error) {
	return _Governance.Contract.Activate(&_Governance.TransactOpts)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) CancelProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "cancelProposal", proposalId)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 proposalId) returns()
func (_Governance *GovernanceSession) CancelProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.CancelProposal(&_Governance.TransactOpts, proposalId)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) CancelProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.CancelProposal(&_Governance.TransactOpts, proposalId)
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

// CancelVoteCancellationProposal is a paid mutator transaction binding the contract method 0x93bffa04.
//
// Solidity: function cancelVoteCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) CancelVoteCancellationProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "cancelVoteCancellationProposal", proposalId)
}

// CancelVoteCancellationProposal is a paid mutator transaction binding the contract method 0x93bffa04.
//
// Solidity: function cancelVoteCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceSession) CancelVoteCancellationProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.CancelVoteCancellationProposal(&_Governance.TransactOpts, proposalId)
}

// CancelVoteCancellationProposal is a paid mutator transaction binding the contract method 0x93bffa04.
//
// Solidity: function cancelVoteCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) CancelVoteCancellationProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.CancelVoteCancellationProposal(&_Governance.TransactOpts, proposalId)
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

// ExecuteCancellationProposal is a paid mutator transaction binding the contract method 0xbf8d68cf.
//
// Solidity: function executeCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) ExecuteCancellationProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "executeCancellationProposal", proposalId)
}

// ExecuteCancellationProposal is a paid mutator transaction binding the contract method 0xbf8d68cf.
//
// Solidity: function executeCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceSession) ExecuteCancellationProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.ExecuteCancellationProposal(&_Governance.TransactOpts, proposalId)
}

// ExecuteCancellationProposal is a paid mutator transaction binding the contract method 0xbf8d68cf.
//
// Solidity: function executeCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) ExecuteCancellationProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.ExecuteCancellationProposal(&_Governance.TransactOpts, proposalId)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address barnAddr) returns()
func (_Governance *GovernanceTransactor) Initialize(opts *bind.TransactOpts, barnAddr common.Address) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "initialize", barnAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address barnAddr) returns()
func (_Governance *GovernanceSession) Initialize(barnAddr common.Address) (*types.Transaction, error) {
	return _Governance.Contract.Initialize(&_Governance.TransactOpts, barnAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address barnAddr) returns()
func (_Governance *GovernanceTransactorSession) Initialize(barnAddr common.Address) (*types.Transaction, error) {
	return _Governance.Contract.Initialize(&_Governance.TransactOpts, barnAddr)
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

// SetAcceptanceThreshold is a paid mutator transaction binding the contract method 0xd1291f19.
//
// Solidity: function setAcceptanceThreshold(uint256 threshold) returns()
func (_Governance *GovernanceTransactor) SetAcceptanceThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setAcceptanceThreshold", threshold)
}

// SetAcceptanceThreshold is a paid mutator transaction binding the contract method 0xd1291f19.
//
// Solidity: function setAcceptanceThreshold(uint256 threshold) returns()
func (_Governance *GovernanceSession) SetAcceptanceThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetAcceptanceThreshold(&_Governance.TransactOpts, threshold)
}

// SetAcceptanceThreshold is a paid mutator transaction binding the contract method 0xd1291f19.
//
// Solidity: function setAcceptanceThreshold(uint256 threshold) returns()
func (_Governance *GovernanceTransactorSession) SetAcceptanceThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetAcceptanceThreshold(&_Governance.TransactOpts, threshold)
}

// SetActiveDuration is a paid mutator transaction binding the contract method 0x24cd62d3.
//
// Solidity: function setActiveDuration(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetActiveDuration(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setActiveDuration", period)
}

// SetActiveDuration is a paid mutator transaction binding the contract method 0x24cd62d3.
//
// Solidity: function setActiveDuration(uint256 period) returns()
func (_Governance *GovernanceSession) SetActiveDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetActiveDuration(&_Governance.TransactOpts, period)
}

// SetActiveDuration is a paid mutator transaction binding the contract method 0x24cd62d3.
//
// Solidity: function setActiveDuration(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetActiveDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetActiveDuration(&_Governance.TransactOpts, period)
}

// SetGracePeriodDuration is a paid mutator transaction binding the contract method 0x342d067a.
//
// Solidity: function setGracePeriodDuration(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetGracePeriodDuration(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setGracePeriodDuration", period)
}

// SetGracePeriodDuration is a paid mutator transaction binding the contract method 0x342d067a.
//
// Solidity: function setGracePeriodDuration(uint256 period) returns()
func (_Governance *GovernanceSession) SetGracePeriodDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetGracePeriodDuration(&_Governance.TransactOpts, period)
}

// SetGracePeriodDuration is a paid mutator transaction binding the contract method 0x342d067a.
//
// Solidity: function setGracePeriodDuration(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetGracePeriodDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetGracePeriodDuration(&_Governance.TransactOpts, period)
}

// SetMinQuorum is a paid mutator transaction binding the contract method 0x563909de.
//
// Solidity: function setMinQuorum(uint256 quorum) returns()
func (_Governance *GovernanceTransactor) SetMinQuorum(opts *bind.TransactOpts, quorum *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setMinQuorum", quorum)
}

// SetMinQuorum is a paid mutator transaction binding the contract method 0x563909de.
//
// Solidity: function setMinQuorum(uint256 quorum) returns()
func (_Governance *GovernanceSession) SetMinQuorum(quorum *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetMinQuorum(&_Governance.TransactOpts, quorum)
}

// SetMinQuorum is a paid mutator transaction binding the contract method 0x563909de.
//
// Solidity: function setMinQuorum(uint256 quorum) returns()
func (_Governance *GovernanceTransactorSession) SetMinQuorum(quorum *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetMinQuorum(&_Governance.TransactOpts, quorum)
}

// SetQueueDuration is a paid mutator transaction binding the contract method 0x53e5056a.
//
// Solidity: function setQueueDuration(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetQueueDuration(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setQueueDuration", period)
}

// SetQueueDuration is a paid mutator transaction binding the contract method 0x53e5056a.
//
// Solidity: function setQueueDuration(uint256 period) returns()
func (_Governance *GovernanceSession) SetQueueDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetQueueDuration(&_Governance.TransactOpts, period)
}

// SetQueueDuration is a paid mutator transaction binding the contract method 0x53e5056a.
//
// Solidity: function setQueueDuration(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetQueueDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetQueueDuration(&_Governance.TransactOpts, period)
}

// SetWarmUpDuration is a paid mutator transaction binding the contract method 0x984690db.
//
// Solidity: function setWarmUpDuration(uint256 period) returns()
func (_Governance *GovernanceTransactor) SetWarmUpDuration(opts *bind.TransactOpts, period *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "setWarmUpDuration", period)
}

// SetWarmUpDuration is a paid mutator transaction binding the contract method 0x984690db.
//
// Solidity: function setWarmUpDuration(uint256 period) returns()
func (_Governance *GovernanceSession) SetWarmUpDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetWarmUpDuration(&_Governance.TransactOpts, period)
}

// SetWarmUpDuration is a paid mutator transaction binding the contract method 0x984690db.
//
// Solidity: function setWarmUpDuration(uint256 period) returns()
func (_Governance *GovernanceTransactorSession) SetWarmUpDuration(period *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.SetWarmUpDuration(&_Governance.TransactOpts, period)
}

// StartCancellationProposal is a paid mutator transaction binding the contract method 0xa31c3443.
//
// Solidity: function startCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactor) StartCancellationProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "startCancellationProposal", proposalId)
}

// StartCancellationProposal is a paid mutator transaction binding the contract method 0xa31c3443.
//
// Solidity: function startCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceSession) StartCancellationProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.StartCancellationProposal(&_Governance.TransactOpts, proposalId)
}

// StartCancellationProposal is a paid mutator transaction binding the contract method 0xa31c3443.
//
// Solidity: function startCancellationProposal(uint256 proposalId) returns()
func (_Governance *GovernanceTransactorSession) StartCancellationProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _Governance.Contract.StartCancellationProposal(&_Governance.TransactOpts, proposalId)
}

// VoteCancellationProposal is a paid mutator transaction binding the contract method 0x15789ef0.
//
// Solidity: function voteCancellationProposal(uint256 proposalId, bool support) returns()
func (_Governance *GovernanceTransactor) VoteCancellationProposal(opts *bind.TransactOpts, proposalId *big.Int, support bool) (*types.Transaction, error) {
	return _Governance.contract.Transact(opts, "voteCancellationProposal", proposalId, support)
}

// VoteCancellationProposal is a paid mutator transaction binding the contract method 0x15789ef0.
//
// Solidity: function voteCancellationProposal(uint256 proposalId, bool support) returns()
func (_Governance *GovernanceSession) VoteCancellationProposal(proposalId *big.Int, support bool) (*types.Transaction, error) {
	return _Governance.Contract.VoteCancellationProposal(&_Governance.TransactOpts, proposalId, support)
}

// VoteCancellationProposal is a paid mutator transaction binding the contract method 0x15789ef0.
//
// Solidity: function voteCancellationProposal(uint256 proposalId, bool support) returns()
func (_Governance *GovernanceTransactorSession) VoteCancellationProposal(proposalId *big.Int, support bool) (*types.Transaction, error) {
	return _Governance.Contract.VoteCancellationProposal(&_Governance.TransactOpts, proposalId, support)
}

// GovernanceCancellationProposalExecutedIterator is returned from FilterCancellationProposalExecuted and is used to iterate over the raw logs and unpacked data for CancellationProposalExecuted events raised by the Governance contract.
type GovernanceCancellationProposalExecutedIterator struct {
	Event *GovernanceCancellationProposalExecuted // Event containing the contract specifics and raw log

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
func (it *GovernanceCancellationProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceCancellationProposalExecuted)
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
		it.Event = new(GovernanceCancellationProposalExecuted)
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
func (it *GovernanceCancellationProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceCancellationProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceCancellationProposalExecuted represents a CancellationProposalExecuted event raised by the Governance contract.
type GovernanceCancellationProposalExecuted struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancellationProposalExecuted is a free log retrieval operation binding the contract event 0xcf7eb172d2813aba9e8e46ce649fb4a18eb29f141f4878d37bf7e988ec05b515.
//
// Solidity: event CancellationProposalExecuted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) FilterCancellationProposalExecuted(opts *bind.FilterOpts, proposalId []*big.Int) (*GovernanceCancellationProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "CancellationProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceCancellationProposalExecutedIterator{contract: _Governance.contract, event: "CancellationProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchCancellationProposalExecuted is a free log subscription operation binding the contract event 0xcf7eb172d2813aba9e8e46ce649fb4a18eb29f141f4878d37bf7e988ec05b515.
//
// Solidity: event CancellationProposalExecuted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) WatchCancellationProposalExecuted(opts *bind.WatchOpts, sink chan<- *GovernanceCancellationProposalExecuted, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "CancellationProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceCancellationProposalExecuted)
				if err := _Governance.contract.UnpackLog(event, "CancellationProposalExecuted", log); err != nil {
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

// ParseCancellationProposalExecuted is a log parse operation binding the contract event 0xcf7eb172d2813aba9e8e46ce649fb4a18eb29f141f4878d37bf7e988ec05b515.
//
// Solidity: event CancellationProposalExecuted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) ParseCancellationProposalExecuted(log types.Log) (*GovernanceCancellationProposalExecuted, error) {
	event := new(GovernanceCancellationProposalExecuted)
	if err := _Governance.contract.UnpackLog(event, "CancellationProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceCancellationProposalStartedIterator is returned from FilterCancellationProposalStarted and is used to iterate over the raw logs and unpacked data for CancellationProposalStarted events raised by the Governance contract.
type GovernanceCancellationProposalStartedIterator struct {
	Event *GovernanceCancellationProposalStarted // Event containing the contract specifics and raw log

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
func (it *GovernanceCancellationProposalStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceCancellationProposalStarted)
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
		it.Event = new(GovernanceCancellationProposalStarted)
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
func (it *GovernanceCancellationProposalStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceCancellationProposalStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceCancellationProposalStarted represents a CancellationProposalStarted event raised by the Governance contract.
type GovernanceCancellationProposalStarted struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancellationProposalStarted is a free log retrieval operation binding the contract event 0x16678c7f63f307d21f71413a77de0aa67a7673745d03209e05f19dcea38e00b0.
//
// Solidity: event CancellationProposalStarted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) FilterCancellationProposalStarted(opts *bind.FilterOpts, proposalId []*big.Int) (*GovernanceCancellationProposalStartedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "CancellationProposalStarted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceCancellationProposalStartedIterator{contract: _Governance.contract, event: "CancellationProposalStarted", logs: logs, sub: sub}, nil
}

// WatchCancellationProposalStarted is a free log subscription operation binding the contract event 0x16678c7f63f307d21f71413a77de0aa67a7673745d03209e05f19dcea38e00b0.
//
// Solidity: event CancellationProposalStarted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) WatchCancellationProposalStarted(opts *bind.WatchOpts, sink chan<- *GovernanceCancellationProposalStarted, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "CancellationProposalStarted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceCancellationProposalStarted)
				if err := _Governance.contract.UnpackLog(event, "CancellationProposalStarted", log); err != nil {
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

// ParseCancellationProposalStarted is a log parse operation binding the contract event 0x16678c7f63f307d21f71413a77de0aa67a7673745d03209e05f19dcea38e00b0.
//
// Solidity: event CancellationProposalStarted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) ParseCancellationProposalStarted(log types.Log) (*GovernanceCancellationProposalStarted, error) {
	event := new(GovernanceCancellationProposalStarted)
	if err := _Governance.contract.UnpackLog(event, "CancellationProposalStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceCancellationProposalVoteIterator is returned from FilterCancellationProposalVote and is used to iterate over the raw logs and unpacked data for CancellationProposalVote events raised by the Governance contract.
type GovernanceCancellationProposalVoteIterator struct {
	Event *GovernanceCancellationProposalVote // Event containing the contract specifics and raw log

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
func (it *GovernanceCancellationProposalVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceCancellationProposalVote)
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
		it.Event = new(GovernanceCancellationProposalVote)
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
func (it *GovernanceCancellationProposalVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceCancellationProposalVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceCancellationProposalVote represents a CancellationProposalVote event raised by the Governance contract.
type GovernanceCancellationProposalVote struct {
	ProposalId *big.Int
	User       common.Address
	Support    bool
	Power      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancellationProposalVote is a free log retrieval operation binding the contract event 0xfff6296a70cd3a166136f34901255c3cd9cbf6dcfe8e5909126a4172b826c809.
//
// Solidity: event CancellationProposalVote(uint256 indexed proposalId, address indexed user, bool support, uint256 power)
func (_Governance *GovernanceFilterer) FilterCancellationProposalVote(opts *bind.FilterOpts, proposalId []*big.Int, user []common.Address) (*GovernanceCancellationProposalVoteIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "CancellationProposalVote", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceCancellationProposalVoteIterator{contract: _Governance.contract, event: "CancellationProposalVote", logs: logs, sub: sub}, nil
}

// WatchCancellationProposalVote is a free log subscription operation binding the contract event 0xfff6296a70cd3a166136f34901255c3cd9cbf6dcfe8e5909126a4172b826c809.
//
// Solidity: event CancellationProposalVote(uint256 indexed proposalId, address indexed user, bool support, uint256 power)
func (_Governance *GovernanceFilterer) WatchCancellationProposalVote(opts *bind.WatchOpts, sink chan<- *GovernanceCancellationProposalVote, proposalId []*big.Int, user []common.Address) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "CancellationProposalVote", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceCancellationProposalVote)
				if err := _Governance.contract.UnpackLog(event, "CancellationProposalVote", log); err != nil {
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

// ParseCancellationProposalVote is a log parse operation binding the contract event 0xfff6296a70cd3a166136f34901255c3cd9cbf6dcfe8e5909126a4172b826c809.
//
// Solidity: event CancellationProposalVote(uint256 indexed proposalId, address indexed user, bool support, uint256 power)
func (_Governance *GovernanceFilterer) ParseCancellationProposalVote(log types.Log) (*GovernanceCancellationProposalVote, error) {
	event := new(GovernanceCancellationProposalVote)
	if err := _Governance.contract.UnpackLog(event, "CancellationProposalVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceCancellationProposalVoteCancelledIterator is returned from FilterCancellationProposalVoteCancelled and is used to iterate over the raw logs and unpacked data for CancellationProposalVoteCancelled events raised by the Governance contract.
type GovernanceCancellationProposalVoteCancelledIterator struct {
	Event *GovernanceCancellationProposalVoteCancelled // Event containing the contract specifics and raw log

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
func (it *GovernanceCancellationProposalVoteCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceCancellationProposalVoteCancelled)
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
		it.Event = new(GovernanceCancellationProposalVoteCancelled)
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
func (it *GovernanceCancellationProposalVoteCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceCancellationProposalVoteCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceCancellationProposalVoteCancelled represents a CancellationProposalVoteCancelled event raised by the Governance contract.
type GovernanceCancellationProposalVoteCancelled struct {
	ProposalId *big.Int
	User       common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancellationProposalVoteCancelled is a free log retrieval operation binding the contract event 0x5e8aa098e39964b153d61768be281eb6183261d780f23758fca70e9711cb7610.
//
// Solidity: event CancellationProposalVoteCancelled(uint256 indexed proposalId, address indexed user)
func (_Governance *GovernanceFilterer) FilterCancellationProposalVoteCancelled(opts *bind.FilterOpts, proposalId []*big.Int, user []common.Address) (*GovernanceCancellationProposalVoteCancelledIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "CancellationProposalVoteCancelled", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceCancellationProposalVoteCancelledIterator{contract: _Governance.contract, event: "CancellationProposalVoteCancelled", logs: logs, sub: sub}, nil
}

// WatchCancellationProposalVoteCancelled is a free log subscription operation binding the contract event 0x5e8aa098e39964b153d61768be281eb6183261d780f23758fca70e9711cb7610.
//
// Solidity: event CancellationProposalVoteCancelled(uint256 indexed proposalId, address indexed user)
func (_Governance *GovernanceFilterer) WatchCancellationProposalVoteCancelled(opts *bind.WatchOpts, sink chan<- *GovernanceCancellationProposalVoteCancelled, proposalId []*big.Int, user []common.Address) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "CancellationProposalVoteCancelled", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceCancellationProposalVoteCancelled)
				if err := _Governance.contract.UnpackLog(event, "CancellationProposalVoteCancelled", log); err != nil {
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

// ParseCancellationProposalVoteCancelled is a log parse operation binding the contract event 0x5e8aa098e39964b153d61768be281eb6183261d780f23758fca70e9711cb7610.
//
// Solidity: event CancellationProposalVoteCancelled(uint256 indexed proposalId, address indexed user)
func (_Governance *GovernanceFilterer) ParseCancellationProposalVoteCancelled(log types.Log) (*GovernanceCancellationProposalVoteCancelled, error) {
	event := new(GovernanceCancellationProposalVoteCancelled)
	if err := _Governance.contract.UnpackLog(event, "CancellationProposalVoteCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceProposalCanceledIterator is returned from FilterProposalCanceled and is used to iterate over the raw logs and unpacked data for ProposalCanceled events raised by the Governance contract.
type GovernanceProposalCanceledIterator struct {
	Event *GovernanceProposalCanceled // Event containing the contract specifics and raw log

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
func (it *GovernanceProposalCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceProposalCanceled)
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
		it.Event = new(GovernanceProposalCanceled)
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
func (it *GovernanceProposalCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceProposalCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceProposalCanceled represents a ProposalCanceled event raised by the Governance contract.
type GovernanceProposalCanceled struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalCanceled is a free log retrieval operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) FilterProposalCanceled(opts *bind.FilterOpts, proposalId []*big.Int) (*GovernanceProposalCanceledIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "ProposalCanceled", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceProposalCanceledIterator{contract: _Governance.contract, event: "ProposalCanceled", logs: logs, sub: sub}, nil
}

// WatchProposalCanceled is a free log subscription operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) WatchProposalCanceled(opts *bind.WatchOpts, sink chan<- *GovernanceProposalCanceled, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "ProposalCanceled", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceProposalCanceled)
				if err := _Governance.contract.UnpackLog(event, "ProposalCanceled", log); err != nil {
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

// ParseProposalCanceled is a log parse operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) ParseProposalCanceled(log types.Log) (*GovernanceProposalCanceled, error) {
	event := new(GovernanceProposalCanceled)
	if err := _Governance.contract.UnpackLog(event, "ProposalCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceProposalCreatedIterator is returned from FilterProposalCreated and is used to iterate over the raw logs and unpacked data for ProposalCreated events raised by the Governance contract.
type GovernanceProposalCreatedIterator struct {
	Event *GovernanceProposalCreated // Event containing the contract specifics and raw log

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
func (it *GovernanceProposalCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceProposalCreated)
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
		it.Event = new(GovernanceProposalCreated)
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
func (it *GovernanceProposalCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceProposalCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceProposalCreated represents a ProposalCreated event raised by the Governance contract.
type GovernanceProposalCreated struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalCreated is a free log retrieval operation binding the contract event 0xc2c021f5d73c63c481d336fbbafec58f694fc45095f00b02d2deb8cca59afe07.
//
// Solidity: event ProposalCreated(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) FilterProposalCreated(opts *bind.FilterOpts, proposalId []*big.Int) (*GovernanceProposalCreatedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "ProposalCreated", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceProposalCreatedIterator{contract: _Governance.contract, event: "ProposalCreated", logs: logs, sub: sub}, nil
}

// WatchProposalCreated is a free log subscription operation binding the contract event 0xc2c021f5d73c63c481d336fbbafec58f694fc45095f00b02d2deb8cca59afe07.
//
// Solidity: event ProposalCreated(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) WatchProposalCreated(opts *bind.WatchOpts, sink chan<- *GovernanceProposalCreated, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "ProposalCreated", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceProposalCreated)
				if err := _Governance.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
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

// ParseProposalCreated is a log parse operation binding the contract event 0xc2c021f5d73c63c481d336fbbafec58f694fc45095f00b02d2deb8cca59afe07.
//
// Solidity: event ProposalCreated(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) ParseProposalCreated(log types.Log) (*GovernanceProposalCreated, error) {
	event := new(GovernanceProposalCreated)
	if err := _Governance.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the Governance contract.
type GovernanceProposalExecutedIterator struct {
	Event *GovernanceProposalExecuted // Event containing the contract specifics and raw log

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
func (it *GovernanceProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceProposalExecuted)
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
		it.Event = new(GovernanceProposalExecuted)
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
func (it *GovernanceProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceProposalExecuted represents a ProposalExecuted event raised by the Governance contract.
type GovernanceProposalExecuted struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) FilterProposalExecuted(opts *bind.FilterOpts, proposalId []*big.Int) (*GovernanceProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceProposalExecutedIterator{contract: _Governance.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *GovernanceProposalExecuted, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceProposalExecuted)
				if err := _Governance.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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

// ParseProposalExecuted is a log parse operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) ParseProposalExecuted(log types.Log) (*GovernanceProposalExecuted, error) {
	event := new(GovernanceProposalExecuted)
	if err := _Governance.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceProposalQueuedIterator is returned from FilterProposalQueued and is used to iterate over the raw logs and unpacked data for ProposalQueued events raised by the Governance contract.
type GovernanceProposalQueuedIterator struct {
	Event *GovernanceProposalQueued // Event containing the contract specifics and raw log

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
func (it *GovernanceProposalQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceProposalQueued)
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
		it.Event = new(GovernanceProposalQueued)
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
func (it *GovernanceProposalQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceProposalQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceProposalQueued represents a ProposalQueued event raised by the Governance contract.
type GovernanceProposalQueued struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalQueued is a free log retrieval operation binding the contract event 0x3358bd34aca93e3ad9a243de48c96b8a820ec804097b77ee85179c1bcdfe9e9f.
//
// Solidity: event ProposalQueued(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) FilterProposalQueued(opts *bind.FilterOpts, proposalId []*big.Int) (*GovernanceProposalQueuedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "ProposalQueued", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceProposalQueuedIterator{contract: _Governance.contract, event: "ProposalQueued", logs: logs, sub: sub}, nil
}

// WatchProposalQueued is a free log subscription operation binding the contract event 0x3358bd34aca93e3ad9a243de48c96b8a820ec804097b77ee85179c1bcdfe9e9f.
//
// Solidity: event ProposalQueued(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) WatchProposalQueued(opts *bind.WatchOpts, sink chan<- *GovernanceProposalQueued, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "ProposalQueued", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceProposalQueued)
				if err := _Governance.contract.UnpackLog(event, "ProposalQueued", log); err != nil {
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

// ParseProposalQueued is a log parse operation binding the contract event 0x3358bd34aca93e3ad9a243de48c96b8a820ec804097b77ee85179c1bcdfe9e9f.
//
// Solidity: event ProposalQueued(uint256 indexed proposalId)
func (_Governance *GovernanceFilterer) ParseProposalQueued(log types.Log) (*GovernanceProposalQueued, error) {
	event := new(GovernanceProposalQueued)
	if err := _Governance.contract.UnpackLog(event, "ProposalQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceVoteIterator is returned from FilterVote and is used to iterate over the raw logs and unpacked data for Vote events raised by the Governance contract.
type GovernanceVoteIterator struct {
	Event *GovernanceVote // Event containing the contract specifics and raw log

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
func (it *GovernanceVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceVote)
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
		it.Event = new(GovernanceVote)
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
func (it *GovernanceVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceVote represents a Vote event raised by the Governance contract.
type GovernanceVote struct {
	ProposalId *big.Int
	User       common.Address
	Support    bool
	Power      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVote is a free log retrieval operation binding the contract event 0x88d35328232823f54954b6627e9f732371656f6daa40cb1b01b27dc7875a7b47.
//
// Solidity: event Vote(uint256 indexed proposalId, address indexed user, bool support, uint256 power)
func (_Governance *GovernanceFilterer) FilterVote(opts *bind.FilterOpts, proposalId []*big.Int, user []common.Address) (*GovernanceVoteIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "Vote", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceVoteIterator{contract: _Governance.contract, event: "Vote", logs: logs, sub: sub}, nil
}

// WatchVote is a free log subscription operation binding the contract event 0x88d35328232823f54954b6627e9f732371656f6daa40cb1b01b27dc7875a7b47.
//
// Solidity: event Vote(uint256 indexed proposalId, address indexed user, bool support, uint256 power)
func (_Governance *GovernanceFilterer) WatchVote(opts *bind.WatchOpts, sink chan<- *GovernanceVote, proposalId []*big.Int, user []common.Address) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "Vote", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceVote)
				if err := _Governance.contract.UnpackLog(event, "Vote", log); err != nil {
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

// ParseVote is a log parse operation binding the contract event 0x88d35328232823f54954b6627e9f732371656f6daa40cb1b01b27dc7875a7b47.
//
// Solidity: event Vote(uint256 indexed proposalId, address indexed user, bool support, uint256 power)
func (_Governance *GovernanceFilterer) ParseVote(log types.Log) (*GovernanceVote, error) {
	event := new(GovernanceVote)
	if err := _Governance.contract.UnpackLog(event, "Vote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovernanceVoteCanceledIterator is returned from FilterVoteCanceled and is used to iterate over the raw logs and unpacked data for VoteCanceled events raised by the Governance contract.
type GovernanceVoteCanceledIterator struct {
	Event *GovernanceVoteCanceled // Event containing the contract specifics and raw log

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
func (it *GovernanceVoteCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovernanceVoteCanceled)
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
		it.Event = new(GovernanceVoteCanceled)
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
func (it *GovernanceVoteCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovernanceVoteCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovernanceVoteCanceled represents a VoteCanceled event raised by the Governance contract.
type GovernanceVoteCanceled struct {
	ProposalId *big.Int
	User       common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteCanceled is a free log retrieval operation binding the contract event 0x12beef84830227673717dd5522ee1228a8004e88dc2678d8740f582264efb2b6.
//
// Solidity: event VoteCanceled(uint256 indexed proposalId, address indexed user)
func (_Governance *GovernanceFilterer) FilterVoteCanceled(opts *bind.FilterOpts, proposalId []*big.Int, user []common.Address) (*GovernanceVoteCanceledIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.FilterLogs(opts, "VoteCanceled", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &GovernanceVoteCanceledIterator{contract: _Governance.contract, event: "VoteCanceled", logs: logs, sub: sub}, nil
}

// WatchVoteCanceled is a free log subscription operation binding the contract event 0x12beef84830227673717dd5522ee1228a8004e88dc2678d8740f582264efb2b6.
//
// Solidity: event VoteCanceled(uint256 indexed proposalId, address indexed user)
func (_Governance *GovernanceFilterer) WatchVoteCanceled(opts *bind.WatchOpts, sink chan<- *GovernanceVoteCanceled, proposalId []*big.Int, user []common.Address) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Governance.contract.WatchLogs(opts, "VoteCanceled", proposalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovernanceVoteCanceled)
				if err := _Governance.contract.UnpackLog(event, "VoteCanceled", log); err != nil {
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

// ParseVoteCanceled is a log parse operation binding the contract event 0x12beef84830227673717dd5522ee1228a8004e88dc2678d8740f582264efb2b6.
//
// Solidity: event VoteCanceled(uint256 indexed proposalId, address indexed user)
func (_Governance *GovernanceFilterer) ParseVoteCanceled(log types.Log) (*GovernanceVoteCanceled, error) {
	event := new(GovernanceVoteCanceled)
	if err := _Governance.contract.UnpackLog(event, "VoteCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
