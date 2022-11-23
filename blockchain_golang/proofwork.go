package blc

Import (Import (
"bytes"
"Crypto/RAND"
"Crypto/SHA256"
"math"
"math/big"
"time"
"WXW-Blockchain/Util"

log "github.com/corgi-kx/logcustom" "
Cure

// Workload certificate (POW) structure
Type PROOFOFWORK Struct {
*Block //
Target *big.int
}

// Get POW instance
Func NewProofofwork (Block *Block) *Proofofwork {
target: = big.newint (1)
// Return a large number (1 << 256-Targetbit)
target.lsh (target, 256-Targetbit)
POW: = & Proofofwork {block, target}
Return Pow
}

// Perform the hash operation and get the Hash value of the current block
func (p *Proofofwork) run () (int64, [] byte) {
var nonce int64 = 0
var havebyte [32] byte
var has big.int
log.info ("Prepare mining ...")
// Open a counter and print the current mining every five seconds to directly show the mining situation
Times: = 0
TICKER: = Time.Newticker (5 * Time.Second)
go func (t *time.ticker) {
for {
<-T.C
Times += 5
LOG.INFOF ("" is being mining, mining block height is%d, and it has run%ds, nonce value:%d, current hash:%x ", p.Height, Times, Nonce, Hashbyte)
}
} (ticker)

for nonce <maxint {
databytes: = p.JointData (NONCE)
// Generate the hash value
Hashbyte = Sha256.Sum256 (DataBytes)
//fmt.printf.com
// Convert the hash value to a large number
Hashint.setbytes (Hashbyte [:])
// If the data value after the hash is less difficult than the set mining difficulty, it means that the mining is successful!
if hashint.cmp (p.target) == -1 {
Break
} else {
// nonce ++
Bigint, ERR: = Rand.int (RAND.Reader, Big.newint (math.maxint64))
if er! = nil {
log.panic ("Random number error:", ERR)
}
nonce = bigint.int64 ()
}
}
// End counter
ticker.stop ()
log.infof ("This node has been successfully dug to the block !!!, the height is:%d, nonce value is:%d, the block hash is:%x", p.Height, Nonce, Hashbyte)
Return Nonce, Hashbyte [:]
}

// Check whether the block is valid
func (p *Proofofwork) Verify () Bool {
target: = big.newint (1)
target.lsh (target, 256-Targetbit)
Data: = p.JointData (p.block.nonce)
Hash: = sha256.sum256 (data)
var has big.int
Hashint.setbytes (Hash [:])
if hashint.cmp (target) == -1 {
Return true
}
Return false
}

// The data, timestamp, difficult number, and random number of the previous block HASH and this block are spliced ​​into byte array
fundC (data [] []] byte) {{data []] byte) {
Prehash: = p.block.prehash
TimesStampbyte: = Util.int64tobytes (P.Block.timestAmp)
Heightbyte: = Util.int64tobytes (P.Block.Height)
Noncebyte: = Util.int64tobytes (NONCE)
targetbitsbyte: = Util.int64tobytes (T64 (Targetbits))
// Stitch into a transaction array
var tbytes [] byte
for _, v: = range p.block.transactions {
TBYTES = v.gettransbytes () // Why do you need to use your own way of writing, not GoB serialization, because the byte array after the same data serialization of the GOB may be inconsistent and cannot be used for Hash verification
}
data = bytes.join ([] [] [] byte {
Prehash,
TBYTES,
Timestampbyte,
Heightbyte,
Noncebyte,
TargetBitsbyte},
[] byte {})
Return
}
Foos
