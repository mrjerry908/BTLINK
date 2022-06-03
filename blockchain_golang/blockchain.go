package block

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"os"
	"time"
)

type blockchain struct {
	BD *database.BlockchainDB /
}


func NewBlockchain() *blockchain {
	blockchain := blockchain{}
	bd := database.New()
	blockchain.BD = bd
	return &blockchain
}


func (bc *blockchain) CreataGenesisTransaction(address string, value int, send Sender) {
	if !IsVaildBitcoinAddress(address) {
		log.Errorf("", address)
		return
	}

	txi := TXInput{[]byte{}, -1, nil, nil}

	wallets := NewWallets(bc.BD)
	genesisKeys, ok := wallets.Wallets[address]
	if !ok {
		log.Fatal
	}

	publicKeyHash := generatePublicKeyHash(genesisKeys.PublicKey)
	txo := TXOutput{value, publicKeyHash}
	ts := Transaction{nil, []TXInput{txi}, []TXOutput{txo}}
	ts.hash()
	tss := []Transaction{ts}

	bc.newGenesisBlockchain(tss)

	NewestBlockHeight = 1
	send.SendVersionToPeers(1)
	fmt.Println(

	utxos := UTXOHandle{bc}
	utxos.ResetUTXODataBase()
}


func (bc *blockchain) newGenesisBlockchain(transaction []Transaction) {

	if len(bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket)) != 0 {
		log.Fatal
	}

	genesisBlock := newGenesisBlock(transaction)

	bc.AddBlock(genesisBlock)
}


func (bc *blockchain) CreataRewardTransaction(address string) Transaction {
	if address == "" {
		log.Warn
		return Transaction{}
	}
	if !IsVaildBitcoinAddress(address) {
		log.Warnf("", address)
		return Transaction{}
	}

	publicKeyHash := getPublicKeyHashFromAddress(address)
	txo := TXOutput{TokenRewardNum, publicKeyHash}
	ts := Transaction{nil, nil, []TXOutput{txo}}
	ts.hash()
	return ts
}

func (bc *blockchain) CreateTransaction(from, to string, amount string, send Sender) {

	if len(bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket)) == 0 {
		log.Error
		return
	}

	if len(bc.BD.View([]byte(RewardAddrMapping), database.AddrBucket)) == 0 {
		log.Warn
	}

	fromSlice := []string{}
	toSlice := []string{}
	amountSlice := []int{}


	err := json.Unmarshal([]byte(from), &fromSlice)
	if err != nil {
		log.Error("json err:", err)
		return
	}
	err = json.Unmarshal([]byte(to), &toSlice)
	if err != nil {
		log.Error("json err:", err)
		return
	}
	err = json.Unmarshal([]byte(amount), &amountSlice)
	if err != nil {
		log.Error("json err:", err)
		return
	}
	if len(fromSlice) != len(toSlice) || len(fromSlice) != len(amountSlice) {
		log.Error
		return
	}

	for i, v := range fromSlice {
		if !IsVaildBitcoinAddress(v) {
			log.Errorf
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}

	for i, v := range toSlice {
		if !IsVaildBitcoinAddress(v) {
			log.Errorf
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}
	for i, v := range amountSlice {
		if v < 0 {
			log.Error
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}

	var tss []Transaction
	wallets := NewWallets(bc.BD)
	for index, fromAddress := range fromSlice {
		fromKeys, ok := wallets.Wallets[fromAddress]
		if !ok {
			log.Errorf(", fromAddress)
			continue
		}
		toKeysPublicKeyHash := getPublicKeyHashFromAddress(toSlice[index])
		if fromAddress == toSlice[index] {
			log.Errorf(" fromAddress)
			return
		}
		u := UTXOHandle{bc}

		utxos := u.findUTXOFromAddress(fromAddress)
		if len(utxos) == 0 {
			log.Errorf(", fromAddress)
			return

		if tss != nil {
			for _, ts := range tss {
			
				for index, vOut := range ts.Vout {
					if bytes.Compare(vOut.PublicKeyHash, generatePublicKeyHash(fromKeys.PublicKey)) != 0 {
						continue
					}
					for _, utxo := range utxos {
						if bytes.Equal(ts.TxHash, utxo.Hash) && index == utxo.Index {
							continue tagVout
						}
					}
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				}
		
				for _, vInt := range ts.Vint {
					for index, utxo := range utxos {
						if bytes.Equal(vInt.TxHash, utxo.Hash) && vInt.Index == utxo.Index {
							utxos = append(utxos[:index], utxos[index+1:]...)
						}
					}
				}

			}
		}

		newTXInput := []TXInput{}
		newTXOutput := []TXOutput{}
		var amount int
		for _, utxo := range utxos {
			amount += utxo.Vout.Value
			newTXInput = append(newTXInput, TXInput{utxo.Hash, utxo.Index, nil, fromKeys.PublicKey})
			if amount > amountSlice[index] {
				tfrom := TXOutput{}
				tfrom.Value = amount - amountSlice[index]
				tfrom.PublicKeyHash = generatePublicKeyHash(fromKeys.PublicKey)
				tTo := TXOutput{}
				tTo.Value = amountSlice[index]
				tTo.PublicKeyHash = toKeysPublicKeyHash
				newTXOutput = append(newTXOutput, tfrom)
				newTXOutput = append(newTXOutput, tTo)
				break
			} else if amount == amountSlice[index] {
				tTo := TXOutput{}
				tTo.Value = amountSlice[index]
				tTo.PublicKeyHash = toKeysPublicKeyHash
				newTXOutput = append(newTXOutput, tTo)
				break
			}
		}

		if amount < amountSlice[index] {
			log.Errorf("  index+1, fromAddress)
			continue
		}
		ts := Transaction{nil, newTXInput, newTXOutput[:]}
		ts.hash()
		tss = append(tss, ts)
	}
	if tss == nil {
		return
	}
	bc.signatureTransactions(tss, wallets)

	send.SendTransToPeers(tss)
}


func (bc *blockchain) Transfer(tss []Transaction, send Sender) {

	if !isGenesisTransaction(tss) {

		bc.verifyTransactionsSign(&tss)
		if len(tss) == 0 {

			return
		}
	}
	
	bc.VerifyTransBalance(&tss)
	if len(tss) == 0 {
	
		return
	}

	rewardTs := bc.CreataRewardTransaction(string(bc.BD.View([]byte(RewardAddrMapping), database.AddrBucket)))
	if rewardTs.TxHash != nil {
		tss = append(tss, rewardTs)
	}
	bc.addBlockchain(tss, send)
}


func (bc *blockchain) VerifyTransBalance(tss *[]Transaction) {

	var balance = map[string]int{}
	for i := range *tss {
		fromAddress := GetAddressFromPublicKey((*tss)[i].Vint[0].PublicKey)

		u := UTXOHandle{bc}
		utxos := u.findUTXOFromAddress(fromAddress)
		if len(utxos) == 0 {
			log.Warnf("%s ！", fromAddress)
			continue
		}
		aomunt := 0
		for _, v := range utxos {
			aomunt += v.Vout.Value
		}
		balance[fromAddress] = aomunt
	}

circle:
	for i := range *tss {
		fromAddress := GetAddressFromPublicKey((*tss)[i].Vint[0].PublicKey)
		u := UTXOHandle{bc}
		utxos := u.findUTXOFromAddress(fromAddress)
		var utxoAmount int //
		var voutAmount int //
		var costAmount int 
		/
		for _, vIn := range (*tss)[i].Vint {
			for _, vUTXO := range utxos {
				if bytes.Equal(vIn.TxHash, vUTXO.Hash) && vIn.Index == vUTXO.Index {
					utxoAmount += vUTXO.Vout.Value
				}
			}
		}
		for _, vOut := range (*tss)[i].Vout {
			if bytes.Equal(getPublicKeyHashFromAddress(fromAddress), vOut.PublicKeyHash) {
				voutAmount += vOut.Value
			}
		}
		costAmount = utxoAmount - voutAmount
		if _, ok := balance[fromAddress]; ok {
			balance[fromAddress] -= costAmount
			if balance[fromAddress] < 0 {
				log.Errorf("%s ", fromAddress)
				*tss = append((*tss)[:i], (*tss)[i+1:]...)
				balance[fromAddress] += costAmount
				goto circle
			}
		} else {
			log.Errorf("%s , fromAddress)
			*tss = append((*tss)[:i], (*tss)[i+1:]...)
			goto circle
		}
	}
	log.Debug
}

func (bc *blockchain) SetRewardAddress(address string) {
	bc.BD.Put([]byte(RewardAddrMapping), []byte(address), database.AddrBucket)
}


func (bc *blockchain) addBlockchain(transaction []Transaction, send Sender) {
	preBlockbyte := bc.BD.View(bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket), database.BlockBucket)
	preBlock := Block{}
	preBlock.Deserialize(preBlockbyte)
	height := preBlock.Height + 1

	nb, err := mineBlock(transaction, bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket), height)
	if err != nil {
		log.Warn(err)
		return
	}

	bc.AddBlock(nb)

	u := UTXOHandle{bc}
	u.Synchrodata(transaction)

	send.SendVersionToPeers(nb.Height)
}

func (bc *blockchain) AddBlock(block *Block) {
	bc.BD.Put(block.Hash, block.Serialize(), database.BlockBucket)
	bci := NewBlockchainIterator(bc)
	currentBlock := bci.Next()
	if currentBlock == nil || currentBlock.Height < block.Height {
		bc.BD.Put([]byte(LastBlockHashMapping), block.Hash, database.BlockBucket)
	}
}


func (bc *blockchain) signatureTransactions(tss []Transaction, wallets *wallets) {
	for i := range tss {
		copyTs := tss[i].customCopy()
		for index := range tss[i].Vint {

			bk := bitcoinKeys{nil, tss[i].Vint[index].PublicKey, nil}
			address := bk.getAddress()

			trans, err := bc.findTransaction(tss, tss[i].Vint[index].TxHash)
			if err != nil {
				log.Fatal(err)
			}
			copyTs.Vint[index].Signature = nil

	
			copyTs.TxHash = copyTs.hashSign()
			copyTs.Vint[index].PublicKey = nil
			privKey := wallets.Wallets[string(address)].PrivateKey

			tss[i].Vint[index].Signature = ellipticCurveSign(privKey, copyTs.TxHash)
		}
	}
}

func (bc *blockchain) verifyTransactionsSign(tss *[]Transaction) {
circle:
	for i := range *tss {
		copyTs := (*tss)[i].customCopy()
		for index, Vin := range (*tss)[i].Vint {
			findTs, err := bc.findTransaction(*tss, Vin.TxHash)
			if err != nil {
				log.Fatal(err)
			}

			if !bytes.Equal(findTs.Vout[Vin.Index].PublicKeyHash, generatePublicKeyHash(Vin.PublicKey)) {
				log.Errorf("", (*tss)[i].TxHash)
				*tss = append((*tss)[:i], (*tss)[i+1:]...)
				goto circle
			}
			copyTs.Vint[index].Signature = nil
			copyTs.Vint[index].PublicKey = findTs.Vout[Vin.Index].PublicKeyHash
			copyTs.TxHash = copyTs.hashSign()
			copyTs.Vint[index].PublicKey = nil
		
			if !ellipticCurveVerify(Vin.PublicKey, Vin.Signature, copyTs.TxHash) {
				log.Errorf("：%x", (*tss)[i].TxHash)
				*tss = append((*tss)[:i], (*tss)[i+1:]...)
				goto circle
			}
		}
	}
	
}


func (bc *blockchain) findTransaction(tss []Transaction, ID []byte) (Transaction, error) {

	if len(tss) != 0 {
		for _, tx := range tss {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return tx, nil
			}
		}
	}
	bci := NewBlockchainIterator(bc)

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return tx, nil
			}
		}
	
		var hashInt big.Int
		hashInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("FindTransaction err : Transaction is not found")
}


func (bc *blockchain) GetLastBlockHeight() int {
	bcl := NewBlockchainIterator(bc)
	lastblock := bcl.Next()
	if lastblock == nil {
		return 0
	}
	return lastblock.Height
}


func (bc *blockchain) GetBlockHashByHeight(height int) []byte {
	bcl := NewBlockchainIterator(bc)
	for {
		currentBlock := bcl.Next()
		if currentBlock == nil {
			return nil
		} else if currentBlock.Height == height {
			return currentBlock.Hash
		} else if isGenesisBlock(currentBlock) {
			return nil
		}
	}
}


func (bc *blockchain) GetBlockByHash(hash []byte) []byte {
	return bc.BD.View(hash, database.BlockBucket)
}


func (bc *blockchain) GetBalance(address string) int {
	if !IsVaildBitcoinAddress(address) {
		log.Errorf("：%s\n", address)
		os.Exit(0)
	}
	var balance int
	uHandle := UTXOHandle{bc}
	utxos := uHandle.findUTXOFromAddress(address)
	for _, v := range utxos {
		balance += v.Vout.Value
	}
	return balance
}


func (bc *blockchain) findAllUTXOs() map[string][]*UTXO {
	utxosMap := make(map[string][]*UTXO)
	txInputmap := make(map[string][]TXInput)
	bcIterator := NewBlockchainIterator(bc)
	for {
		currentBlock := bcIterator.Next()
		if currentBlock == nil {
			return nil
		}
		
		for i := len(currentBlock.Transactions) - 1; i >= 0; i-- {
			var utxos = []*UTXO{}
			ts := currentBlock.Transactions[i]
			for _, vInt := range ts.Vint {
				txInputmap[string(vInt.TxHash)] = append(txInputmap[string(vInt.TxHash)], vInt)
			}

		VoutTag:
			for index, vOut := range ts.Vout {
				if txInputmap[string(ts.TxHash)] == nil {
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				} else {
					for _, vIn := range txInputmap[string(ts.TxHash)] {
						if vIn.Index == index {
							continue VoutTag
						}
					}
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				}
				utxosMap[string(ts.TxHash)] = utxos
			}
		}

		if isGenesisBlock(currentBlock) {
			break
		}
	}
	return utxosMap
}

func (bc *blockchain) PrintAllBlockInfo() {
	blcIterator := NewBlockchainIterator(bc)
	for {
		block := blcIterator.Next()
		if block == nil {
			log.Error("")
			return
		}
		fmt.Println("========================================================================================================")
		fmt.Printf(hash         %x\n", block.Hash)
		fmt.Println("  	------------------------------------------------------------")
		for _, v := range block.Transactions {
			fmt.Printf("   	id:  %x\n", v.TxHash)
			fmt.Println("   	  tx_input：")
			for _, vIn := range v.Vint {
				fmt.Printf("			  %x\n", vIn.TxHash)
				fmt.Printf("		   %d\n", vIn.Index)
				fmt.Printf("			    %x\n", vIn.Signature)
				fmt.Printf("			    %x\n", vIn.PublicKey)
				fmt.Printf("			   %s\n", GetAddressFromPublicKey(vIn.PublicKey))
			}
			fmt.Println("  	  tx_output：")
			for index, vOut := range v.Vout {
				fmt.Printf("			   %d    \n", vOut.Value)
				fmt.Printf("			  %x    \n", vOut.PublicKeyHash)
				fmt.Printf("		   %s\n", GetAddressFromPublicKeyHash(vOut.PublicKeyHash))
				if len(v.Vout) != 1 && index != len(v.Vout)-1 {
					fmt.Println("			---------------")
				}
			}
		}
		fmt.Println("  	--------------------------------------------------------------------")
		fmt.Printf("          %s\n", time.Unix(block.TimeStamp, 0).Format("03:04:05 PM"))
		fmt.Printf("        %d\n", block.Height)
		fmt.Printf("           %d\n", block.Nonce)
		fmt.Printf("     %x\n", block.PreHash)
		var hashInt big.Int
		hashInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	fmt.Println("========================================================================================================")
}
