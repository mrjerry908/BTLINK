*
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package icon

import (
"fmt"
"path/filepath"
"strings"
"time"

"github.com/blocktree/go-owcrypt"
"github.com/blocktree/openwallet/v2/common/file"
"github.com/shopspring/decimal"
)

/*
The tool can read the default configuration data of each currency wallet,
The configuration data of the currency wallet is placed in conf/{symbol}.conf, for example: ADA.conf, BTC.conf, ETH.conf.
Executing the wmd wallet -s <symbol> command will first check whether the configuration file of the currency wallet exists.
None: ConfigFlow is executed, configuration file is initialized.
Yes: Execute regular commands.
Users can also modify the configuration file through wmd config -s.
Or execute wmd config flow to perform the configuration initialization process again.
*/

const (
// currency
Symbol = "ICX"
MasterKey = "Icon seed"
CurveType = owcrypt.ECC_CURVE_SECP256K1
)

type WalletConfig struct {
// currency
Symbol string
MasterKey string

keyDir string
//Address export path
addressDir string
//Configuration file path
configFilePath string
//Configuration file name
configFileName string
//local database file path
dbPath string
// backup path
backupDir string
//Wallet service API
ServerAPI string
//The minimum transfer amount of the address
minTransfer decimal.Decimal
	//Transaction Fees
fees decimal.Decimal
//summary threshold
Threshold decimal.Decimal
//summary address
SumAddress string
//Total execution interval time
CycleSeconds time.Duration
//default configuration content
DefaultConfig string
//curve type
CurveType uint32
//stepLimit miner fee cap
StepLimit int64
}

func NewConfig(symbol string, masterKey string) *WalletConfig {
c := WalletConfig{}

// currency
c.Symbol = symbol
c.MasterKey = masterKey
c.CurveType = CurveType
//key backup path
c.keyDir = filepath.Join("data", strings.ToLower(c.Symbol), "key")
//Address export path
c.addressDir = filepath.Join("data", strings.ToLower(c.Symbol), "address")
//blockchain data
//blockchainDir = filepath.Join("data", strings.ToLower(Symbol), "blockchain")
//Configuration file path
c.configFilePath = filepath.Join("conf")
//Configuration file name
c.configFileName = c.Symbol + ".ini"
//local database file path
c.dbPath = filepath.Join("data", strings.ToLower(c.Symbol), "db")
// backup path
c.backupDir = filepath.Join("data", strings.ToLower(c.Symbol), "backup")
//Wallet service API
c.ServerAPI = ""
c.StepLimit = 100000
//path of wallet installation
//The minimum transfer amount of the address
c.minTransfer = decimal.Zero
	//Transaction Fees
c.fees = decimal.Zero
//summary threshold
c.Threshold = decimal.NewFromFloat(5)
//summary address
c.SumAddress = ""
//Total execution interval time
c.CycleSeconds = time.Second * 30
//default configuration content
c.DefaultConfig = `
# node api url
apiUrl = ""
# transaction fix fees
fees = "0.001"
# transaction max step limit,
stepLimit = 100000
# the minimum amount could transfer of address
minTransfer = "0"
# the safe address that wallet send money to.
sumAddress = ""
# when wallet's balance is over this value, the wallet willl send money to [sumAddress]
#unit is ICX
threshold = ""
# summary task timer cycle time, sample: 1h, 1h1m , 2m, 30s, 3m20s etc...
cycleSeconds = "30s"
`
return &c
}

//printConfig Print config information
func (wc *WalletConfig) PrintConfig() error {
wc.InitConfig()
//read configuration
absFile := filepath.Join(wc.configFilePath, wc.configFileName)

fmt.Printf("-------------------------------------------- --------------\n")
file.PrintFile(absFile)
fmt.Printf("-------------------------------------------- --------------\n")

return nil

}

//initConfig initializes the configuration file
func (wc *WalletConfig) InitConfig() {
//read configuration
absFile := filepath.Join(wc.configFilePath, wc.configFileName)
if !file.Exists(absFile) {
file.MkdirAll(wc.configFilePath)
file.WriteFile(absFile, []byte(wc.DefaultConfig), false)
}
