// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xerror"
)

// Error type.
const (
	ER_KEY_SIGNATURE_VERIFY_FAILED                 int = 2101
	ER_HDKEY_DERIVE_HARD_FROM_PUBLIC               int = 3101
	ER_HDKEY_PRIV_EXTKEY_INVALID                   int = 3102
	ER_HDKEY_CHILD_INVALID                         int = 3103
	ER_HDKEY_SERIALIZED_KEY_WRONG_SIZE             int = 3104
	ER_HDKEY_CHECKSUM_MISMATCH                     int = 3105
	ER_HDKEY_DERIVE_PATH_INVALID                   int = 3106
	ER_ADDRESS_CHECKSUM_MISMATCH                   int = 4101
	ER_ADDRESS_TYPE_UNKNOWN                        int = 4102
	ER_ADDRESS_FORMAT_MALFORMED                    int = 4103
	ER_ADDRESS_SIZE_MALFORMED                      int = 4104
	ER_ADDRESS_WITNESS_VERSION_UNSUPPORTED         int = 4105
	ER_SCRIPT_TYPE_UNKNOWN                         int = 5100
	ER_SCRIPT_STANDARD_ADDRESS_TYPE_UNSUPPORTED    int = 5105
	ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED        int = 5106
	ER_TRANSACTION_SIGN_OUT_INDEX                  int = 6000
	ER_TRANSACTION_SIGN_REDEEM_EMPTY               int = 6001
	ER_TRANSACTION_BUILDER_AMOUNT_NOT_ENOUGH_ERROR int = 6101
	ER_TRANSACTION_BUILDER_FROM_EMPTY              int = 6102
	ER_TRANSACTION_BUILDER_CHANGETO_EMPTY          int = 6103
	ER_TRANSACTION_BUILDER_SENDTO_EMPTY            int = 6104
	ER_TRANSACTION_BUILDER_SIGN_KEY_EMPTY          int = 6105
	ER_TRANSACTION_BUILDER_MIN_FEE_NOT_ENOUGH      int = 6106
	ER_TRANSACTION_BUILDER_FEE_TOO_HIGH            int = 6107
	ER_TRANSACTION_PARTIALLY_MAGIC_MISMATCH        int = 6201
	ER_MICROPAYMENT_LOCKTIME_MISMATCH              int = 6301
	ER_MICROPAYMENT_REFUND_BOND_MISMATCH           int = 6302
)

// Errors -- the jump table of error.
var Errors = map[int]*xerror.Error{
	ER_KEY_SIGNATURE_VERIFY_FAILED:                 {Num: ER_KEY_SIGNATURE_VERIFY_FAILED, State: "TKS00", Message: "key.signature.verify.failed"},
	ER_HDKEY_DERIVE_HARD_FROM_PUBLIC:               {Num: ER_HDKEY_DERIVE_HARD_FROM_PUBLIC, State: "THK00", Message: "hdkey.derive.a.hardened.key.from.public.key.error"},
	ER_HDKEY_PRIV_EXTKEY_INVALID:                   {Num: ER_HDKEY_PRIV_EXTKEY_INVALID, State: "THK00", Message: "hdkey.unable.to.create.private.keys.from.a.public.extened.key"},
	ER_HDKEY_CHILD_INVALID:                         {Num: ER_HDKEY_CHILD_INVALID, State: "THK00", Message: "hdkey.at.this.index.is.invalid"},
	ER_HDKEY_CHECKSUM_MISMATCH:                     {Num: ER_HDKEY_CHECKSUM_MISMATCH, State: "THK00", Message: "hdkey.derive.path.invalid[%v]"},
	ER_HDKEY_DERIVE_PATH_INVALID:                   {Num: ER_HDKEY_DERIVE_PATH_INVALID, State: "THK00", Message: "hdkey.checksum.mismatch"},
	ER_ADDRESS_CHECKSUM_MISMATCH:                   {Num: ER_ADDRESS_CHECKSUM_MISMATCH, State: "THK00", Message: "address.checksum.mismatch"},
	ER_ADDRESS_TYPE_UNKNOWN:                        {Num: ER_ADDRESS_TYPE_UNKNOWN, State: "TADDR0", Message: "address.unknown.type[%v]"},
	ER_ADDRESS_FORMAT_MALFORMED:                    {Num: ER_ADDRESS_FORMAT_MALFORMED, State: "TADDR0", Message: "address.unknown.format[%v]"},
	ER_ADDRESS_SIZE_MALFORMED:                      {Num: ER_ADDRESS_SIZE_MALFORMED, State: "TADDR0", Message: "address.size[%v].invalid"},
	ER_ADDRESS_WITNESS_VERSION_UNSUPPORTED:         {Num: ER_ADDRESS_WITNESS_VERSION_UNSUPPORTED, State: "TADDR0", Message: "address.witness.address.version[%v].unsupported"},
	ER_SCRIPT_TYPE_UNKNOWN:                         {Num: ER_SCRIPT_TYPE_UNKNOWN, State: "TS000", Message: "script.unknow.type[%v]"},
	ER_SCRIPT_STANDARD_ADDRESS_TYPE_UNSUPPORTED:    {Num: ER_SCRIPT_STANDARD_ADDRESS_TYPE_UNSUPPORTED, State: "TS000", Message: "script.standard.unsupported.address.type[%v]"},
	ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED:        {Num: ER_SCRIPT_STANDARD_PUBKEYS_LE_NREQUIRED, State: "TS000", Message: "script.standard.pubkeys[%v].less.than.nrequired[%v]"},
	ER_TRANSACTION_SIGN_OUT_INDEX:                  {Num: ER_TRANSACTION_SIGN_OUT_INDEX, State: "TTX00", Message: "transaction.sign.idx[%v].out.index[%v]"},
	ER_TRANSACTION_SIGN_REDEEM_EMPTY:               {Num: ER_TRANSACTION_SIGN_REDEEM_EMPTY, State: "TTX00", Message: "transaction.sign.idx[%v].redeem.can.not.be.nil.since.keys[%v]>1"},
	ER_TRANSACTION_BUILDER_AMOUNT_NOT_ENOUGH_ERROR: {Num: ER_TRANSACTION_BUILDER_AMOUNT_NOT_ENOUGH_ERROR, State: "TTB00", Message: "transaction.builder.amount.totalout[%v].more.than.totalin[%v]"},
	ER_TRANSACTION_BUILDER_FROM_EMPTY:              {Num: ER_TRANSACTION_BUILDER_FROM_EMPTY, State: "TTB00", Message: "transaction.builder.from.is.empty.at.group.idx[%v]"},
	ER_TRANSACTION_BUILDER_CHANGETO_EMPTY:          {Num: ER_TRANSACTION_BUILDER_CHANGETO_EMPTY, State: "TTB00", Message: "transaction.builder.changeto.is.empty"},
	ER_TRANSACTION_BUILDER_SENDTO_EMPTY:            {Num: ER_TRANSACTION_BUILDER_SENDTO_EMPTY, State: "TTB00", Message: "transaction.builder.sendto.is.empty.at.group.idx[%v]"},
	ER_TRANSACTION_BUILDER_SIGN_KEY_EMPTY:          {Num: ER_TRANSACTION_BUILDER_SIGN_KEY_EMPTY, State: "TTB00", Message: "transaction.builder.sign.but.key.is.empty.at.input.idx[%v]"},
	ER_TRANSACTION_BUILDER_MIN_FEE_NOT_ENOUGH:      {Num: ER_TRANSACTION_BUILDER_MIN_FEE_NOT_ENOUGH, State: "TTB00", Message: "transaction.builder.min.fee[%v].not.enough.from.change.value[%v]"},
	ER_TRANSACTION_BUILDER_FEE_TOO_HIGH:            {Num: ER_TRANSACTION_BUILDER_FEE_TOO_HIGH, State: "TTB00", Message: "transaction.builder.fee[%v].too.high.than.max.fee[%v]"},
	ER_TRANSACTION_PARTIALLY_MAGIC_MISMATCH:        {Num: ER_TRANSACTION_PARTIALLY_MAGIC_MISMATCH, State: "TTP00", Message: "transaction.partially.request.magic.mismatch.want[%x].got[%x]"},
	ER_MICROPAYMENT_LOCKTIME_MISMATCH:              {Num: ER_MICROPAYMENT_LOCKTIME_MISMATCH, State: "TM000", Message: "micropayment.locktime.mismatch.want[%v].got[%v]"},
	ER_MICROPAYMENT_REFUND_BOND_MISMATCH:           {Num: ER_MICROPAYMENT_REFUND_BOND_MISMATCH, State: "TM000", Message: "micropayment.refund.bond.mismatch"},
}
