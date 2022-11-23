package blc

Import (Import (
"bytes"
"Crypto/ECDSA"
"Crypto/Elliptic"
"Crypto/RAND"
"Crypto/SHA256"
"Encoding/GOB"
"Encoding/hex"
log "github.com/corgi-kx/logcustom" "
"math/big"
"WXW-Blockchain/Util"
Cure

// Transaction list information
Type Transaction Struct {
Txhash [] byte
// UTXO input
Vint [] txinput
// UTXO output
Vout [] txoutput
}

// iscoinbase Checks WHether The Transaction is Coinbase
Func (T *Transaction) iscoinbase () BOOL {
Return len (t.vint) == 1 && len (t.vint [0] .txhash) == 0 && t.vint [0] .index == -1
}

// Input the input of this transaction, and the output for hash operations after the trading hash (txhash)
Func (t *transaction) have () {
TBYTES: = t.Serialize ()
// Add the random number byte (what is the purpose?)
randomnumber: = Util.genlateRandom ()
randombyte: = Util.int64tobytes (Randomnumber)
Sumbyte: = bytes.join ([] [] [] byte {tBytes, randombyte}, [] byte ("")
Hashbyte: = Sha256.Sum256 (Sumbyte)
t.txhash = Hashbyte [:]
}

// Convert the members in the entire transaction into a byte array in turn, stitch it into the whole
func (t *transaction) gettransbytes () [] byte {byte {
if t.txhash == nil || t.vout == nil {
LOG.PANIC ("The transaction information is incomplete and cannot be spliced ​​into a byte array")
Return nil
}
var Transbytes [] byte
Transbytes = Append (Transbytes, T.txhash ...)
for _, v: = Range T.Vint {
Transbytes = Append (Transbytes, v.txhash ...)
Transbytes = Append (Transbytes, Util.int64tobytes (int64 (v.Index)) ...)
Transbytes = Append (Transbytes, v.signature ...)
}
for _, v: = Range t.vout {
Transbytes = Append (Transbytes, Util.int64tobytes (int64 (v.value)) ...)
Transbytes = Append (Transbytes, v.Publickeyhash ...)
}
Return Transbytes
}

// Serialize transaction to [] byte
func (t *transaction) serialize () [] byte {
var result bytes.buffer
Encoder: = Gob.newenCoder (& Result)

Err: = encoder.encode (t)
if er! = nil {
PANIC (ERR)
}
Return result.bytes ()
}

Func (T *Transaction) Sign (PrivateKey *ECDSA.PRIVATEKEY, Prevts Map [String] Transaction) {
if t.iscoinbase () {
Return
}

txcopy: = t.trimmedcopy ()

for inid, vint: = range txcopy.vint {
prevts: = Prevts
txcopy.vint [inid] .signature = nil
txcopy.vint [inid] .publickey = prevts.vout [vint.index] .publickeyhash
txcopy.txhash = txcopy.hash ()
txcopy.vint [inid] .publickey = nil

R, S, Err: = ECDSA.SIGN (RAND.Reader, PrivateKey, TXCOPY.TXHASH)
if er! = nil {
log.fatal (ERR)
}
signature: = Append (R.bytes (), s.bytes () ...)

t.vint [Inid] .signature = Signature
}
}

func (t transfer) have () [] byte {byte {
txcopy: = t
txcopy.txhash = [] byte {}
Hash: = sha256.sum256 (txcopy.serialize ()))
Return Hash [:]
}

Func (T *Transaction) TrimmedCopy () Transaction {
var inputs [] txinput
var outputs [] txoutput

for _, vint: = Range t.vint {
Inputs = APPEND (Inputs, TXINPUT {vint.txhash, vint.index, nil, nil, nil, nil})
}

for _, vout: = Range t.vout {
Outputs = APPEND (Outputs, TXOUTPUT {Vout.Value, VOUT.PUBLICKEYHASH})
}

txcopy: = transaction {t.txhash, inputs, outputs}

Return txcopy
}

Func (T *Transaction) Verify (Prevts Map [String] Transaction) BOOL {
txcopy: = t.trimmedcopy ()
Curve: = Elliptic.p256 ()

for inid, vint: = range t.vint {
prevts: = Prevts
txcopy.vint [inid] .signature = nil
txcopy.vint [inid] .publickey = prevts.vout [vint.index] .publickeyhash
txcopy.txhash = txcopy.hash ()
txcopy.vint [inid] .publickey = nil

R: = Big.int {}
s: = big.int {}
siglen: = len (vint.signature)
R.setbytes (vint.signature [: (siglen / 2)]))
s.setbytes (vint.signature [siglen/2:])

x: = big.int {}
y: = big.int {}
Keylen: = len (vint.publickey)
x.setbytes (vint.publickey [: keylen/2])
y.setbytes (vint.publickey [keylen/2:])

rawpublickey: = ECDSA.PUBLICKEY {Curve: Curve, x: & x, y: & y}
If! Ecdsa.verify (& rawpublickey, txcopy.txhash, & r, & s) {
Return false
}
}
Return true
}
Foos
