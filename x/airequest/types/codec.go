package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/oraichain/orai/x/airequest/exported"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO: Register the modules msgs
	cdc.RegisterConcrete(MsgSetKYCRequest{}, "airequest/SetKYCRequest", nil)
	cdc.RegisterConcrete(MsgSetClassificationRequest{}, "airequest/SetClassificationRequest", nil)
	cdc.RegisterConcrete(MsgSetOCRRequest{}, "airequest/SetOCRRequest", nil)
	cdc.RegisterConcrete(MsgSetPriceRequest{}, "airequest/SetPriceRequest", nil)
	cdc.RegisterInterface((*exported.AIRequestI)(nil), nil) // has to be pointer of interface
	cdc.RegisterConcrete(&AIRequest{}, "airequest/AIRequest", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
