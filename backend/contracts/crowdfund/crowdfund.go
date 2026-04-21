// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package crowdfund

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CrowdFundCampaign is an auto generated low-level Go binding around an user-defined struct.
type CrowdFundCampaign struct {
	Id        *big.Int
	Creator   common.Address
	Title     string
	Goal      *big.Int
	Pledged   *big.Int
	Deadline  *big.Int
	Withdrawn bool
}

// CrowdFundMetaData contains all meta data concerning the CrowdFund contract.
var CrowdFundMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"MAX_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"campaigns\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"goal\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pledged\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contributions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createCampaign\",\"inputs\":[{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"goal\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"duration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fund\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getCampaign\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCrowdFund.Campaign\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"goal\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pledged\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextCampaignId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"refund\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CampaignCreated\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"title\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"goal\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Funded\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"funder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Refunded\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"funder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CampaignEnded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CampaignNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CampaignStillActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FundsAlreadyWithdrawn\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GoalAlreadyReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GoalNotReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidContribution\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidGoal\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotCampaignCreator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToRefund\",\"inputs\":[]}]",
	Bin: "0x608060405234801561000f575f80fd5b50610e178061001d5f395ff3fe60806040526004361061008f575f3560e01c80635598f8cc116100575780635598f8cc146101715780637903a7561461019d578063b1724b46146101b1578063b6a6d177146101c7578063ca1d209d146101dc575f80fd5b8063141961bc14610093578063278ecde1146100ce5780632e1a7d4d146100ef5780633020580b1461010e5780633d891f591461013b575b5f80fd5b34801561009e575f80fd5b506100b26100ad366004610a42565b6101ef565b6040516100c59796959493929190610a9c565b60405180910390f35b3480156100d9575f80fd5b506100ed6100e8366004610a42565b6102be565b005b3480156100fa575f80fd5b506100ed610109366004610a42565b61045f565b348015610119575f80fd5b5061012d610128366004610aed565b610627565b6040519081526020016100c5565b348015610146575f80fd5b5061012d610155366004610b65565b600260209081525f928352604080842090915290825290205481565b34801561017c575f80fd5b5061019061018b366004610a42565b6107d8565b6040516100c59190610b9e565b3480156101a8575f80fd5b5061012d5f5481565b3480156101bc575f80fd5b5061012d62278d0081565b3480156101d2575f80fd5b5061012d610e1081565b6100ed6101ea366004610a42565b610944565b600160208190525f91825260409091208054918101546002820180546001600160a01b03909216929161022190610c0c565b80601f016020809104026020016040519081016040528092919081815260200182805461024d90610c0c565b80156102985780601f1061026f57610100808354040283529160200191610298565b820191905f5260205f20905b81548152906001019060200180831161027b57829003601f168201915b505050506003830154600484015460058501546006909501549394919390925060ff1687565b5f818152600160208190526040909120908101546001600160a01b03166102f8576040516337f9b70b60e11b815260040160405180910390fd5b806005015442101561031d57604051634e5b565b60e11b815260040160405180910390fd5b80600301548160040154106103455760405163465c128f60e11b815260040160405180910390fd5b5f8281526002602090815260408083203384529091528120549081900361037f5760405163f76aef6560e01b815260040160405180910390fd5b5f8381526002602090815260408083203380855292528083208390555183908381818185875af1925050503d805f81146103d4576040519150601f19603f3d011682016040523d82523d5f602084013e6103d9565b606091505b50509050806104215760405162461bcd60e51b815260206004820152600f60248201526e1514905394d1915497d19052531151608a1b60448201526064015b60405180910390fd5b604051828152339085907f7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4906020015b60405180910390a350505050565b5f818152600160208190526040909120908101546001600160a01b0316610499576040516337f9b70b60e11b815260040160405180910390fd5b60018101546001600160a01b031633146104c6576040516339a0685760e21b815260040160405180910390fd5b80600501544210156104eb57604051634e5b565b60e11b815260040160405180910390fd5b806003015481600401541015610514576040516378c754c960e01b815260040160405180910390fd5b600681015460ff161561053a57604051633f72752760e21b815260040160405180910390fd5b60068101805460ff191660019081179091556004820154908201546040515f916001600160a01b03169083908381818185875af1925050503d805f811461059c576040519150601f19603f3d011682016040523d82523d5f602084013e6105a1565b606091505b50509050806105e45760405162461bcd60e51b815260206004820152600f60248201526e1514905394d1915497d19052531151608a1b6044820152606401610418565b60018301546040518381526001600160a01b039091169085907fcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef20237290602001610451565b5f825f0361064857604051639b60eb4d60e01b815260040160405180910390fd5b610e1082108061065a575062278d0082115b1561067857604051637616640160e01b815260040160405180910390fd5b5f8054908061068683610c58565b9190505590506040518060e00160405280828152602001336001600160a01b0316815260200186868080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201829052509385525050506020820186905260408201526060016106fc8442610c70565b81525f60209182018190528381526001808352604091829020845181559284015190830180546001600160a01b0319166001600160a01b03909216919091179055820151600282019061074f9082610ce3565b50606082015160038201556080820151600482015560a0820151600582015560c0909101516006909101805460ff191691151591909117905533817fdc26653af5b99b2da33e2ad69ee6600d9aeccc82b034501db4338309615ca2388787876107b88842610c70565b6040516107c89493929190610da3565b60405180910390a3949350505050565b61081d6040518060e001604052805f81526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f151581525090565b5f828152600160208181526040808420815160e08101835281548152938101546001600160a01b0316928401929092526002820180549184019161086090610c0c565b80601f016020809104026020016040519081016040528092919081815260200182805461088c90610c0c565b80156108d75780601f106108ae576101008083540402835291602001916108d7565b820191905f5260205f20905b8154815290600101906020018083116108ba57829003601f168201915b50505091835250506003820154602080830191909152600483015460408301526005830154606083015260069092015460ff1615156080909101528101519091506001600160a01b031661093e576040516337f9b70b60e11b815260040160405180910390fd5b92915050565b5f818152600160208190526040909120908101546001600160a01b031661097e576040516337f9b70b60e11b815260040160405180910390fd5b806005015442106109a25760405163154eb81560e21b815260040160405180910390fd5b345f036109c25760405163652122d960e01b815260040160405180910390fd5b34816004015f8282546109d59190610c70565b90915550505f82815260026020908152604080832033845290915281208054349290610a02908490610c70565b9091555050604051348152339083907f38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca69269060200160405180910390a35050565b5f60208284031215610a52575f80fd5b5035919050565b5f81518084525f5b81811015610a7d57602081850181015186830182015201610a61565b505f602082860101526020601f19601f83011685010191505092915050565b8781526001600160a01b038716602082015260e0604082018190525f90610ac590830188610a59565b90508560608301528460808301528360a083015282151560c083015298975050505050505050565b5f805f8060608587031215610b00575f80fd5b843567ffffffffffffffff80821115610b17575f80fd5b818701915087601f830112610b2a575f80fd5b813581811115610b38575f80fd5b886020828501011115610b49575f80fd5b6020928301999098509187013596604001359550909350505050565b5f8060408385031215610b76575f80fd5b8235915060208301356001600160a01b0381168114610b93575f80fd5b809150509250929050565b602081528151602082015260018060a01b0360208301511660408201525f604083015160e06060840152610bd6610100840182610a59565b905060608401516080840152608084015160a084015260a084015160c084015260c0840151151560e08401528091505092915050565b600181811c90821680610c2057607f821691505b602082108103610c3e57634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52601160045260245ffd5b5f60018201610c6957610c69610c44565b5060010190565b8082018082111561093e5761093e610c44565b634e487b7160e01b5f52604160045260245ffd5b601f821115610cde57805f5260205f20601f840160051c81016020851015610cbc5750805b601f840160051c820191505b81811015610cdb575f8155600101610cc8565b50505b505050565b815167ffffffffffffffff811115610cfd57610cfd610c83565b610d1181610d0b8454610c0c565b84610c97565b602080601f831160018114610d44575f8415610d2d5750858301515b5f19600386901b1c1916600185901b178555610d9b565b5f85815260208120601f198616915b82811015610d7257888601518255948401946001909101908401610d53565b5085821015610d8f57878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b60608152836060820152838560808301375f608085830101525f6080601f19601f87011683010190508360208301528260408301529594505050505056fea2646970667358221220c55ff0ebfd18710fe688b1c404f4f5033554397dc9fa58dfa36b98b758afb56f64736f6c63430008180033",
}

// CrowdFundABI is the input ABI used to generate the binding from.
// Deprecated: Use CrowdFundMetaData.ABI instead.
var CrowdFundABI = CrowdFundMetaData.ABI

// CrowdFundBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrowdFundMetaData.Bin instead.
var CrowdFundBin = CrowdFundMetaData.Bin

// DeployCrowdFund deploys a new Ethereum contract, binding an instance of CrowdFund to it.
func DeployCrowdFund(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrowdFund, error) {
	parsed, err := CrowdFundMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrowdFundBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrowdFund{CrowdFundCaller: CrowdFundCaller{contract: contract}, CrowdFundTransactor: CrowdFundTransactor{contract: contract}, CrowdFundFilterer: CrowdFundFilterer{contract: contract}}, nil
}

// CrowdFund is an auto generated Go binding around an Ethereum contract.
type CrowdFund struct {
	CrowdFundCaller     // Read-only binding to the contract
	CrowdFundTransactor // Write-only binding to the contract
	CrowdFundFilterer   // Log filterer for contract events
}

// CrowdFundCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrowdFundCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdFundTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrowdFundTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdFundFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrowdFundFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdFundSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrowdFundSession struct {
	Contract     *CrowdFund        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrowdFundCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrowdFundCallerSession struct {
	Contract *CrowdFundCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// CrowdFundTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrowdFundTransactorSession struct {
	Contract     *CrowdFundTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CrowdFundRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrowdFundRaw struct {
	Contract *CrowdFund // Generic contract binding to access the raw methods on
}

// CrowdFundCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrowdFundCallerRaw struct {
	Contract *CrowdFundCaller // Generic read-only contract binding to access the raw methods on
}

// CrowdFundTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrowdFundTransactorRaw struct {
	Contract *CrowdFundTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrowdFund creates a new instance of CrowdFund, bound to a specific deployed contract.
func NewCrowdFund(address common.Address, backend bind.ContractBackend) (*CrowdFund, error) {
	contract, err := bindCrowdFund(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrowdFund{CrowdFundCaller: CrowdFundCaller{contract: contract}, CrowdFundTransactor: CrowdFundTransactor{contract: contract}, CrowdFundFilterer: CrowdFundFilterer{contract: contract}}, nil
}

// NewCrowdFundCaller creates a new read-only instance of CrowdFund, bound to a specific deployed contract.
func NewCrowdFundCaller(address common.Address, caller bind.ContractCaller) (*CrowdFundCaller, error) {
	contract, err := bindCrowdFund(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrowdFundCaller{contract: contract}, nil
}

// NewCrowdFundTransactor creates a new write-only instance of CrowdFund, bound to a specific deployed contract.
func NewCrowdFundTransactor(address common.Address, transactor bind.ContractTransactor) (*CrowdFundTransactor, error) {
	contract, err := bindCrowdFund(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrowdFundTransactor{contract: contract}, nil
}

// NewCrowdFundFilterer creates a new log filterer instance of CrowdFund, bound to a specific deployed contract.
func NewCrowdFundFilterer(address common.Address, filterer bind.ContractFilterer) (*CrowdFundFilterer, error) {
	contract, err := bindCrowdFund(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrowdFundFilterer{contract: contract}, nil
}

// bindCrowdFund binds a generic wrapper to an already deployed contract.
func bindCrowdFund(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrowdFundMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrowdFund *CrowdFundRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrowdFund.Contract.CrowdFundCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrowdFund *CrowdFundRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrowdFund.Contract.CrowdFundTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrowdFund *CrowdFundRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrowdFund.Contract.CrowdFundTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrowdFund *CrowdFundCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrowdFund.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrowdFund *CrowdFundTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrowdFund.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrowdFund *CrowdFundTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrowdFund.Contract.contract.Transact(opts, method, params...)
}

// MAXDURATION is a free data retrieval call binding the contract method 0xb1724b46.
//
// Solidity: function MAX_DURATION() view returns(uint256)
func (_CrowdFund *CrowdFundCaller) MAXDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFund.contract.Call(opts, &out, "MAX_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXDURATION is a free data retrieval call binding the contract method 0xb1724b46.
//
// Solidity: function MAX_DURATION() view returns(uint256)
func (_CrowdFund *CrowdFundSession) MAXDURATION() (*big.Int, error) {
	return _CrowdFund.Contract.MAXDURATION(&_CrowdFund.CallOpts)
}

// MAXDURATION is a free data retrieval call binding the contract method 0xb1724b46.
//
// Solidity: function MAX_DURATION() view returns(uint256)
func (_CrowdFund *CrowdFundCallerSession) MAXDURATION() (*big.Int, error) {
	return _CrowdFund.Contract.MAXDURATION(&_CrowdFund.CallOpts)
}

// MINDURATION is a free data retrieval call binding the contract method 0xb6a6d177.
//
// Solidity: function MIN_DURATION() view returns(uint256)
func (_CrowdFund *CrowdFundCaller) MINDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFund.contract.Call(opts, &out, "MIN_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINDURATION is a free data retrieval call binding the contract method 0xb6a6d177.
//
// Solidity: function MIN_DURATION() view returns(uint256)
func (_CrowdFund *CrowdFundSession) MINDURATION() (*big.Int, error) {
	return _CrowdFund.Contract.MINDURATION(&_CrowdFund.CallOpts)
}

// MINDURATION is a free data retrieval call binding the contract method 0xb6a6d177.
//
// Solidity: function MIN_DURATION() view returns(uint256)
func (_CrowdFund *CrowdFundCallerSession) MINDURATION() (*big.Int, error) {
	return _CrowdFund.Contract.MINDURATION(&_CrowdFund.CallOpts)
}

// Campaigns is a free data retrieval call binding the contract method 0x141961bc.
//
// Solidity: function campaigns(uint256 ) view returns(uint256 id, address creator, string title, uint256 goal, uint256 pledged, uint256 deadline, bool withdrawn)
func (_CrowdFund *CrowdFundCaller) Campaigns(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id        *big.Int
	Creator   common.Address
	Title     string
	Goal      *big.Int
	Pledged   *big.Int
	Deadline  *big.Int
	Withdrawn bool
}, error) {
	var out []interface{}
	err := _CrowdFund.contract.Call(opts, &out, "campaigns", arg0)

	outstruct := new(struct {
		Id        *big.Int
		Creator   common.Address
		Title     string
		Goal      *big.Int
		Pledged   *big.Int
		Deadline  *big.Int
		Withdrawn bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Creator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Title = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Goal = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Pledged = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Deadline = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Withdrawn = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// Campaigns is a free data retrieval call binding the contract method 0x141961bc.
//
// Solidity: function campaigns(uint256 ) view returns(uint256 id, address creator, string title, uint256 goal, uint256 pledged, uint256 deadline, bool withdrawn)
func (_CrowdFund *CrowdFundSession) Campaigns(arg0 *big.Int) (struct {
	Id        *big.Int
	Creator   common.Address
	Title     string
	Goal      *big.Int
	Pledged   *big.Int
	Deadline  *big.Int
	Withdrawn bool
}, error) {
	return _CrowdFund.Contract.Campaigns(&_CrowdFund.CallOpts, arg0)
}

// Campaigns is a free data retrieval call binding the contract method 0x141961bc.
//
// Solidity: function campaigns(uint256 ) view returns(uint256 id, address creator, string title, uint256 goal, uint256 pledged, uint256 deadline, bool withdrawn)
func (_CrowdFund *CrowdFundCallerSession) Campaigns(arg0 *big.Int) (struct {
	Id        *big.Int
	Creator   common.Address
	Title     string
	Goal      *big.Int
	Pledged   *big.Int
	Deadline  *big.Int
	Withdrawn bool
}, error) {
	return _CrowdFund.Contract.Campaigns(&_CrowdFund.CallOpts, arg0)
}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_CrowdFund *CrowdFundCaller) Contributions(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFund.contract.Call(opts, &out, "contributions", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_CrowdFund *CrowdFundSession) Contributions(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _CrowdFund.Contract.Contributions(&_CrowdFund.CallOpts, arg0, arg1)
}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_CrowdFund *CrowdFundCallerSession) Contributions(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _CrowdFund.Contract.Contributions(&_CrowdFund.CallOpts, arg0, arg1)
}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns((uint256,address,string,uint256,uint256,uint256,bool))
func (_CrowdFund *CrowdFundCaller) GetCampaign(opts *bind.CallOpts, campaignId *big.Int) (CrowdFundCampaign, error) {
	var out []interface{}
	err := _CrowdFund.contract.Call(opts, &out, "getCampaign", campaignId)

	if err != nil {
		return *new(CrowdFundCampaign), err
	}

	out0 := *abi.ConvertType(out[0], new(CrowdFundCampaign)).(*CrowdFundCampaign)

	return out0, err

}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns((uint256,address,string,uint256,uint256,uint256,bool))
func (_CrowdFund *CrowdFundSession) GetCampaign(campaignId *big.Int) (CrowdFundCampaign, error) {
	return _CrowdFund.Contract.GetCampaign(&_CrowdFund.CallOpts, campaignId)
}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns((uint256,address,string,uint256,uint256,uint256,bool))
func (_CrowdFund *CrowdFundCallerSession) GetCampaign(campaignId *big.Int) (CrowdFundCampaign, error) {
	return _CrowdFund.Contract.GetCampaign(&_CrowdFund.CallOpts, campaignId)
}

// NextCampaignId is a free data retrieval call binding the contract method 0x7903a756.
//
// Solidity: function nextCampaignId() view returns(uint256)
func (_CrowdFund *CrowdFundCaller) NextCampaignId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFund.contract.Call(opts, &out, "nextCampaignId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextCampaignId is a free data retrieval call binding the contract method 0x7903a756.
//
// Solidity: function nextCampaignId() view returns(uint256)
func (_CrowdFund *CrowdFundSession) NextCampaignId() (*big.Int, error) {
	return _CrowdFund.Contract.NextCampaignId(&_CrowdFund.CallOpts)
}

// NextCampaignId is a free data retrieval call binding the contract method 0x7903a756.
//
// Solidity: function nextCampaignId() view returns(uint256)
func (_CrowdFund *CrowdFundCallerSession) NextCampaignId() (*big.Int, error) {
	return _CrowdFund.Contract.NextCampaignId(&_CrowdFund.CallOpts)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0x3020580b.
//
// Solidity: function createCampaign(string title, uint256 goal, uint256 duration) returns(uint256 campaignId)
func (_CrowdFund *CrowdFundTransactor) CreateCampaign(opts *bind.TransactOpts, title string, goal *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _CrowdFund.contract.Transact(opts, "createCampaign", title, goal, duration)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0x3020580b.
//
// Solidity: function createCampaign(string title, uint256 goal, uint256 duration) returns(uint256 campaignId)
func (_CrowdFund *CrowdFundSession) CreateCampaign(title string, goal *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.CreateCampaign(&_CrowdFund.TransactOpts, title, goal, duration)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0x3020580b.
//
// Solidity: function createCampaign(string title, uint256 goal, uint256 duration) returns(uint256 campaignId)
func (_CrowdFund *CrowdFundTransactorSession) CreateCampaign(title string, goal *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.CreateCampaign(&_CrowdFund.TransactOpts, title, goal, duration)
}

// Fund is a paid mutator transaction binding the contract method 0xca1d209d.
//
// Solidity: function fund(uint256 campaignId) payable returns()
func (_CrowdFund *CrowdFundTransactor) Fund(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.contract.Transact(opts, "fund", campaignId)
}

// Fund is a paid mutator transaction binding the contract method 0xca1d209d.
//
// Solidity: function fund(uint256 campaignId) payable returns()
func (_CrowdFund *CrowdFundSession) Fund(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.Fund(&_CrowdFund.TransactOpts, campaignId)
}

// Fund is a paid mutator transaction binding the contract method 0xca1d209d.
//
// Solidity: function fund(uint256 campaignId) payable returns()
func (_CrowdFund *CrowdFundTransactorSession) Fund(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.Fund(&_CrowdFund.TransactOpts, campaignId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 campaignId) returns()
func (_CrowdFund *CrowdFundTransactor) Refund(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.contract.Transact(opts, "refund", campaignId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 campaignId) returns()
func (_CrowdFund *CrowdFundSession) Refund(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.Refund(&_CrowdFund.TransactOpts, campaignId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 campaignId) returns()
func (_CrowdFund *CrowdFundTransactorSession) Refund(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.Refund(&_CrowdFund.TransactOpts, campaignId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 campaignId) returns()
func (_CrowdFund *CrowdFundTransactor) Withdraw(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.contract.Transact(opts, "withdraw", campaignId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 campaignId) returns()
func (_CrowdFund *CrowdFundSession) Withdraw(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.Withdraw(&_CrowdFund.TransactOpts, campaignId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 campaignId) returns()
func (_CrowdFund *CrowdFundTransactorSession) Withdraw(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFund.Contract.Withdraw(&_CrowdFund.TransactOpts, campaignId)
}

// CrowdFundCampaignCreatedIterator is returned from FilterCampaignCreated and is used to iterate over the raw logs and unpacked data for CampaignCreated events raised by the CrowdFund contract.
type CrowdFundCampaignCreatedIterator struct {
	Event *CrowdFundCampaignCreated // Event containing the contract specifics and raw log

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
func (it *CrowdFundCampaignCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundCampaignCreated)
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
		it.Event = new(CrowdFundCampaignCreated)
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
func (it *CrowdFundCampaignCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundCampaignCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundCampaignCreated represents a CampaignCreated event raised by the CrowdFund contract.
type CrowdFundCampaignCreated struct {
	CampaignId *big.Int
	Creator    common.Address
	Title      string
	Goal       *big.Int
	Deadline   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCampaignCreated is a free log retrieval operation binding the contract event 0xdc26653af5b99b2da33e2ad69ee6600d9aeccc82b034501db4338309615ca238.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed creator, string title, uint256 goal, uint256 deadline)
func (_CrowdFund *CrowdFundFilterer) FilterCampaignCreated(opts *bind.FilterOpts, campaignId []*big.Int, creator []common.Address) (*CrowdFundCampaignCreatedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _CrowdFund.contract.FilterLogs(opts, "CampaignCreated", campaignIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundCampaignCreatedIterator{contract: _CrowdFund.contract, event: "CampaignCreated", logs: logs, sub: sub}, nil
}

// WatchCampaignCreated is a free log subscription operation binding the contract event 0xdc26653af5b99b2da33e2ad69ee6600d9aeccc82b034501db4338309615ca238.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed creator, string title, uint256 goal, uint256 deadline)
func (_CrowdFund *CrowdFundFilterer) WatchCampaignCreated(opts *bind.WatchOpts, sink chan<- *CrowdFundCampaignCreated, campaignId []*big.Int, creator []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _CrowdFund.contract.WatchLogs(opts, "CampaignCreated", campaignIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundCampaignCreated)
				if err := _CrowdFund.contract.UnpackLog(event, "CampaignCreated", log); err != nil {
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

// ParseCampaignCreated is a log parse operation binding the contract event 0xdc26653af5b99b2da33e2ad69ee6600d9aeccc82b034501db4338309615ca238.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed creator, string title, uint256 goal, uint256 deadline)
func (_CrowdFund *CrowdFundFilterer) ParseCampaignCreated(log types.Log) (*CrowdFundCampaignCreated, error) {
	event := new(CrowdFundCampaignCreated)
	if err := _CrowdFund.contract.UnpackLog(event, "CampaignCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundFundedIterator is returned from FilterFunded and is used to iterate over the raw logs and unpacked data for Funded events raised by the CrowdFund contract.
type CrowdFundFundedIterator struct {
	Event *CrowdFundFunded // Event containing the contract specifics and raw log

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
func (it *CrowdFundFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundFunded)
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
		it.Event = new(CrowdFundFunded)
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
func (it *CrowdFundFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundFunded represents a Funded event raised by the CrowdFund contract.
type CrowdFundFunded struct {
	CampaignId *big.Int
	Funder     common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterFunded is a free log retrieval operation binding the contract event 0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926.
//
// Solidity: event Funded(uint256 indexed campaignId, address indexed funder, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) FilterFunded(opts *bind.FilterOpts, campaignId []*big.Int, funder []common.Address) (*CrowdFundFundedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}

	logs, sub, err := _CrowdFund.contract.FilterLogs(opts, "Funded", campaignIdRule, funderRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundFundedIterator{contract: _CrowdFund.contract, event: "Funded", logs: logs, sub: sub}, nil
}

// WatchFunded is a free log subscription operation binding the contract event 0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926.
//
// Solidity: event Funded(uint256 indexed campaignId, address indexed funder, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) WatchFunded(opts *bind.WatchOpts, sink chan<- *CrowdFundFunded, campaignId []*big.Int, funder []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}

	logs, sub, err := _CrowdFund.contract.WatchLogs(opts, "Funded", campaignIdRule, funderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundFunded)
				if err := _CrowdFund.contract.UnpackLog(event, "Funded", log); err != nil {
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

// ParseFunded is a log parse operation binding the contract event 0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926.
//
// Solidity: event Funded(uint256 indexed campaignId, address indexed funder, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) ParseFunded(log types.Log) (*CrowdFundFunded, error) {
	event := new(CrowdFundFunded)
	if err := _CrowdFund.contract.UnpackLog(event, "Funded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the CrowdFund contract.
type CrowdFundRefundedIterator struct {
	Event *CrowdFundRefunded // Event containing the contract specifics and raw log

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
func (it *CrowdFundRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundRefunded)
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
		it.Event = new(CrowdFundRefunded)
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
func (it *CrowdFundRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundRefunded represents a Refunded event raised by the CrowdFund contract.
type CrowdFundRefunded struct {
	CampaignId *big.Int
	Funder     common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed campaignId, address indexed funder, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) FilterRefunded(opts *bind.FilterOpts, campaignId []*big.Int, funder []common.Address) (*CrowdFundRefundedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}

	logs, sub, err := _CrowdFund.contract.FilterLogs(opts, "Refunded", campaignIdRule, funderRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundRefundedIterator{contract: _CrowdFund.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed campaignId, address indexed funder, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *CrowdFundRefunded, campaignId []*big.Int, funder []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}

	logs, sub, err := _CrowdFund.contract.WatchLogs(opts, "Refunded", campaignIdRule, funderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundRefunded)
				if err := _CrowdFund.contract.UnpackLog(event, "Refunded", log); err != nil {
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

// ParseRefunded is a log parse operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed campaignId, address indexed funder, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) ParseRefunded(log types.Log) (*CrowdFundRefunded, error) {
	event := new(CrowdFundRefunded)
	if err := _CrowdFund.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the CrowdFund contract.
type CrowdFundWithdrawnIterator struct {
	Event *CrowdFundWithdrawn // Event containing the contract specifics and raw log

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
func (it *CrowdFundWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundWithdrawn)
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
		it.Event = new(CrowdFundWithdrawn)
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
func (it *CrowdFundWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundWithdrawn represents a Withdrawn event raised by the CrowdFund contract.
type CrowdFundWithdrawn struct {
	CampaignId *big.Int
	Creator    common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372.
//
// Solidity: event Withdrawn(uint256 indexed campaignId, address indexed creator, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) FilterWithdrawn(opts *bind.FilterOpts, campaignId []*big.Int, creator []common.Address) (*CrowdFundWithdrawnIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _CrowdFund.contract.FilterLogs(opts, "Withdrawn", campaignIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundWithdrawnIterator{contract: _CrowdFund.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372.
//
// Solidity: event Withdrawn(uint256 indexed campaignId, address indexed creator, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *CrowdFundWithdrawn, campaignId []*big.Int, creator []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _CrowdFund.contract.WatchLogs(opts, "Withdrawn", campaignIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundWithdrawn)
				if err := _CrowdFund.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372.
//
// Solidity: event Withdrawn(uint256 indexed campaignId, address indexed creator, uint256 amount)
func (_CrowdFund *CrowdFundFilterer) ParseWithdrawn(log types.Log) (*CrowdFundWithdrawn, error) {
	event := new(CrowdFundWithdrawn)
	if err := _CrowdFund.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
