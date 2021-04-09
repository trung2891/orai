package airequest

import (
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
)

const (
	ModuleName                     = types.ModuleName
	RouterKey                      = types.RouterKey
	StoreKey                       = types.StoreKey
	QuerierRoute                   = types.QuerierRoute
	DefaultParamspace              = types.DefaultParamspace
	EventTypeSetAIRequest          = types.EventTypeSetAIRequest
	AttributeRequestValidator      = types.AttributeRequestValidator
	EventTypeRequestWithData       = types.EventTypeRequestWithData
	AttributeRequestDSources       = types.AttributeRequestDSources
	AttributeRequestTCases         = types.AttributeRequestTCases
	AttributeRequestID             = types.AttributeRequestID
	AttributeContract              = types.AttributeContract
	AttributeRequestCreator        = types.AttributeRequestCreator
	AttributeRequestValidatorCount = types.AttributeRequestValidatorCount
	AttributeRequestInput          = types.AttributeRequestInput
	AttributeRequestExpectedOutput = types.AttributeRequestExpectedOutput
)

var (
	NewKeeper          = keeper.NewKeeper
	NewQuerier         = keeper.NewQuerier
	NewMsgSetAIRequest = types.NewMsgSetAIRequest
	ModuleCdc          = types.ModuleCdc
	RegisterCodec      = types.RegisterCodec
	NewGenesisState    = types.NewGenesisState
	RequestKeyPrefix   = types.RequestKeyPrefix
	ErrRequestNotFound = types.ErrRequestNotFound
	NewAIRequest       = types.NewAIRequest
)

type (
	Keeper          = keeper.Keeper
	MsgSetAIRequest = types.MsgSetAIRequest
	AIRequest       = types.AIRequest
	GenesisState    = types.GenesisState
)
