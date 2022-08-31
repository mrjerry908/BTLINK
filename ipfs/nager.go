func (wm *WalletManager) Transfer(sk []byte, from, to, value string, stepLimit, nonce int64) (string, error) {
tx, hash := wm.CalculateTxHash(from, to, value, stepLimit, nonce)

sig, err := wm.signTransaction(hash[:], sk)
if err != nil {
return "", err
}
//log.Info(sig)

tx["signature"] = sig
//log.Info(tx)

ret, err := wm.WalletClient.Call_icx_sendTransaction(tx)
if err != nil {
log.Error(err)
return "", err
}

return ret, nil
}

//CreateNewWallet creates a wallet
func (wm *WalletManager) CreateNewWallet(name, password string) (*openwallet.Wallet, string, error) {
var (
err error
wallets[]*openwallet.Wallet
)

//Check if the wallet name exists
wallets, err = wm.GetWallets()
for _, w := range wallets {
if w.Alias ​​== name {
return nil, "", errors. New("The wallet's alias is duplicated!")
}
}

fmt.Printf("Create new wallet keystore...\n")

seed, err := hdkeychain.GenerateSeed(32)
if err != nil {
return nil, "", err
}

extSeed, err := hdkeystore.GetExtendSeed(seed, wm.Config.MasterKey)
if err != nil {
return nil, "", err
}

key, keyFile, err := hdkeystore.StoreHDKeyWithSeed(wm.Config.keyDir, name, password, extSeed, hdkeystore.StandardScryptN, hdkeystore.StandardScryptP)
if err != nil {
return nil, "", err
}

file.MkdirAll(wm.Config.dbPath)
file.MkdirAll(wm.Config.keyDir)

w := &openwallet.Wallet{
WalletID: key.KeyID,
Alias: key.Alias,
KeyFile: keyFile,
DBFile: filepath.Join(wm.Config.dbPath, key.FileName()+".db"),
}

w.SaveToDB()

return w, keyFile, nil
}

//GetWalletKeys Get the wallet list by loading the keystore file with the given file path
func (wm *WalletManager) GetWallets() ([]*openwallet.Wallet, error) {
wallets, err := openwallet.GetWalletsByKeyDir(wm.Config.keyDir)
if err != nil {
return nil, err
}

for _, w := range wallets {
w.DBFile = filepath.Join(wm.Config.dbPath, w.FileName()+".db")
}

return wallets, nil
}

func (wm *WalletManager) AddWalletInSummary(wid string, wallet *openwallet.Wallet) {
wm.WalletsInSum[wid] = wallet
}

//get wallet balance
func (wm *WalletManager) getWalletBalance(wallet *openwallet.Wallet) (decimal.Decimal, []*openwallet.Address, error) {
var (
synCount = 10
quit = make(chan struct{})
done = 0 //done marker
shouldDone = 0 //Total amount to be done
)

db, err := wallet.OpenDB()
if err != nil {
return decimal.NewFromFloat(0), nil, err
}
defer db.Close()

var addrs[]*openwallet.Address
db.All(&addrs)

balance, _ := decimal.NewFromString("0")
count := len(addrs)
if count <= 0 {
log.Std.Info("This wallet have 0 address!!!")
return decimal.NewFromFloat(0), nil, nil
} else {
log.Std.Info("wallet %s have %d addresses, please wait minutes to get wallet balance", wallet.Alias, count)
}

//production channel
producer := make(chan[]*openwallet.Address)
defer close(producer)

//consumer channel
worker := make(chan[]*openwallet.Address)
defer close(worker)

//statistic balance
go func(addrs chan []*openwallet.Address) {
for balances := range addrs {
//
//balances := <-addrs

for _, b := range balances {
addrB, _ := decimal.NewFromString(b.Balance)
balance = balance.Add(addrB)
}

// Accumulated number of completed threads
done++
if done == shouldDone {
close(quit) //Close the channel, which is equivalent to passing nil to the channel
}
}
}(worker)

/* Calculate the number of synCount threads, the number of internal runs */
//The number of loops in each thread, processed in parallel by synCount threads
runCount := count / synCount
otherCount := count % synCount

if runCount > 0 {
for i := 0; i < synCount; i++ {
			//start
//log.Std.Info("Start get balance thread[%d]", i)
start := i * runCount
end := (i+1)*runCount - 1
as := addrs[start:end]

go func(producer chan []*openwallet.Address, addrs []*openwallet.Address, wm *WalletManager) {
var bs []*openwallet.Address
for _, a := range addrs {
b, err := wm.WalletClient.Call_icx_getBalance(a.Address)
if err != nil {
log.Error(err)
continue
}
a.Balance = string(b)
bs = append(bs, a)
}

producer <- bs
}(producer, as, wm)

shouldDone++
}
}

if otherCount > 0 {
//
//log.Std.Info("Start get balance thread[REST]")
start := runCount * synCount
as := addrs[start:]

go func(producer chan []*openwallet.Address, addrs []*openwallet.Address, wm *WalletManager) {
var bs []*openwallet.Address
for _, a := range addrs {
b, err := wm.WalletClient.Call_icx_getBalance(a.Address)
if err != nil {
log.Error(err)
continue
}
a.Balance = string(b)
bs = append(bs, a)
}

producer <- bs
}(producer, as, wm)

shouldDone++
}

values ​​:= make([][]*openwallet.Address, 0)
outputAddress := make([]*openwallet.Address, 0)

//The following uses the production consumption mode
for {
var activeWorker chan<- []*openwallet.Address
var activeValue []*openwallet.Address

//When the data queue has data, release the top and activate consumption
if len(values) > 0 {
activeWorker = worker
activeValue = values[0]
}

select {
//The generator continuously generates data and inserts it at the end of the data queue
case pa := <-producer:
values ​​= append(values, pa)
//When the consumer is activated, transmit the data to the consumer and dequeue the top data
outputAddress = append(outputAddress, pa...)
//log.Debug("produced")
case activeWorker <- activeValue:
values ​​= values[1:]
//log.D
