package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/hako/durafmt"
	"github.com/shopspring/decimal"
)

func CleanUpHex(s string) string {
	s = strings.Replace(strings.TrimPrefix(s, "0x"), " ", "", -1)

	return strings.ToLower(s)
}

func ValidateAccount(accountAddress string) (string, error) {
	accountAddress = CleanUpHex(accountAddress)
	// check account length
	if len(accountAddress) != 40 {
		return "", errors.New("invalid account address")
	}

	return NormalizeAddress(accountAddress), nil
}

func AppendNotEmpty(slice []string, str string) []string {
	if str != "" {
		return append(slice, str)
	}

	return slice
}

// HexStrToBigIntStr transforms a hex sting like "0xff" to a big int string like "15". Arbitrary length values are possible.
func HexStrToBigIntStr(hexString string) (string, error) {
	value, err := HexStrToBigInt(hexString)
	return value.String(), err
}

// HexStrToBigInt transforms a hex sting like "0xff" to a big.Int. Arbitrary length values are possible.
func HexStrToBigInt(hexString string) (*big.Int, error) {
	value := new(big.Int)
	_, ok := value.SetString(Trim0x(hexString), 16)
	if !ok {
		return value, fmt.Errorf("could not transform hex string to big int: %s", hexString)
	}

	return value, nil
}

// Trim0x removes the "0x" prefix of hexes if it exists
func Trim0x(str string) string {
	return strings.TrimPrefix(str, "0x")
}

func Topic2Address(topic string) string {
	topic = Trim0x(topic)
	return NormalizeAddress(topic[24:])
}

func LogIsEvent(log web3types.Log, abi abi.ABI, event string) bool {
	if len(log.Topics) == 0 {
		return false
	}

	return CleanUpHex(log.Topics[0]) == CleanUpHex(abi.Events[event].ID.String())
}

func NormalizeAddress(addr string) string {
	return "0x" + Trim0x(strings.ToLower(addr))
}

func HumanDuration(seconds int64) string {
	return durafmt.Parse(time.Duration(seconds) * time.Second).String()
}

func PrettyPercent(d decimal.Decimal) string {
	return d.Mul(decimal.NewFromInt(100)).StringFixed(2)
}
