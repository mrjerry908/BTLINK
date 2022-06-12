package client
	
	import (
	"errors"
	qtypes "github.com/btlink/qbase/types"
	
	qcliacc "github.com/btlink/qbase/client/account"
	"github.com/btlink/qbase/client/context"
	qclitx "github.com/btlink/qbase/client/tx"
	"github.com/btlink/qbase/txs"
	atxs "github.com/btlink/qos/module/approve/txs"
	approvetypes "github.com/btlink/qos/module/approve/types"
	"github.com/btlink/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	)
	
	type operateType int
	
	const (
	createType operateType = iota
	increaseType
	decreaseType
	useType
	cancleType
	
	flagFrom = "from"
	flagTo = "to"
	flagCoins = "coins"
	)
	
	// 
	func QueryApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "approve",
	Short: "Get approve info by from and to",
	RunE: func(cmd *cobra.Command, args []string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	fromAddr, err := qcliacc.GetAddrFromFlag(cliCtx, flagFrom)
	if err != nil {
	return err
	}
	toAddr, err := qcliacc.GetAddrFromFlag(cliCtx, flagTo)
	if err != nil {
	return err
	}
	
	approve, err := getApproveInfo(cliCtx, fromAddr, toAddr)
	if err != nil {
	return err
	}
	
	return cliCtx.PrintResult(approve)
	},
	}
	
	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	
	return cmd
	}
	
	func getApproveInfo(cliCtx context.CLIContext, approve, beneficiary qtypes.AccAddress) (approvetypes.Approve, error) {
	queryPath := "store/approve/key"
	output, err := cliCtx.Query(queryPath, approvetypes.BuildApproveKey(approve, beneficiary))
	if err != nil {
	return approvetypes.Approve{}, err
	}
	if output == nil {
	return approvetypes.Approve{}, context.RecordsNotFoundError
	}
	

	appr := approvetypes.Approve{}
	cliCtx.Codec.MustUnmarshalBinaryBare(output, &appr)
	
	return appr, nil
	}
	
	func QueryApprovesCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "approves <name or address>",
	Short: "Query approves by user",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	
	fromAddr, err := qcliacc.GetAddrFromValue(cliCtx, args[0])
	if err != nil {
	return err
	}
	
	approves, err := queryUserApproves(cliCtx, fromAddr)
	
	if len(approves) == 0 {
	return context.RecordsNotFoundError
	}
	
	if err != nil {
	return err
	}
	
	return cliCtx.PrintResult(approves)
	},
	}
	
	return cmd
	}
	
	func queryUserApproves(cliCtx context.CLIContext, approve qtypes.AccAddress) (approves []approvetypes.Approve, err error) {
	key := approvetypes.BuildApproveByFromKey(approve)
	path := "store/approve/subspace"
	output, err := cliCtx.Query(path, key)
	
	if err != nil {
	return nil, err
	}
	
	var pairs []qtypes.KVPair
	err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(output, &pairs)
	if err != nil {
	return nil, err
	}
	
	for _, v := range pairs {
	appr := approvetypes.Approve{}
	cliCtx.Codec.MustUnmarshalBinaryBare(v.Value, &appr)
	approves = append(approves, appr)
	}
	
	return
	}
	

	func CreateApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "create-approve",
	Short: "Create approve",
	RunE: func(cmd *cobra.Command, args []string) error {
	return handleApproveOperation(cdc, createType)
	},
	}
	
	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)
	
	return cmd
	}
	

	func IncreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "increase-approve",
	Short: "Increase approve",
	RunE: func(cmd *cobra.Command, args []string) error {
	return handleApproveOperation(cdc, increaseType)
	},
	}
	
	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)
	
	return cmd
	}
	

	func DecreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "decrease-approve",
	Short: "Decrease approve",
	RunE: func(cmd *cobra.Command, args []string) error {
	return handleApproveOperation(cdc, decreaseType)
	},
	}
	
	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)
	
	return cmd
	}

	func UseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "use-approve",
	Short: "Use approve",
	RunE: func(cmd *cobra.Command, args []string) error {
	return handleApproveOperation(cdc, useType)
	},
	}
	
	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)
	
	return cmd
	}
	
	
	func CancelApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
	Use: "cancel-approve",
	Short: "Cancel approve",
	RunE: func(cmd *cobra.Command, args []string) error {
	return handleApproveOperation(cdc, cancleType)
	},
	}
	
	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	
	return cmd
	}
	
	func handleApproveOperation(cdc *amino.Codec, operType operateType) error {
	iTxBuilder := func(ctx context.CLIContext) (txs.ITx, error) {

	fromAddr, err := qcliacc.GetAddrFromFlag(ctx, flagFrom)
	if err != nil {
	return nil, err
	}
	
	toAddr, err := qcliacc.GetAddrFromFlag(ctx, flagTo)
	if err != nil {
	return nil, err
	}
	
	
	if operType == cancleType {
	tx := atxs.TxCancelApprove{
	From: fromAddr,
	To: toAddr,
	}
	if err = tx.ValidateInputs(); err != nil {
	return nil, err
	}
	return tx, nil
	}
	
	// 
	qos, qscs, err := types.ParseCoins(viper.GetString(flagCoins))
	if err != nil {
	return nil, err
	}
	
	//
	tx := approvetypes.NewApprove(fromAddr, toAddr, qos, qscs)
	if err = tx.Valid(); err != nil {
	return nil, err
	}
	switch operType {
	case createType:
	return atxs.TxCreateApprove{Approve: tx}, nil
	case increaseType:
	return atxs.TxIncreaseApprove{Approve: tx}, nil
	case decreaseType:
	return atxs.TxDecreaseApprove{Approve: tx}, nil
	case useType:
	return atxs.TxUseApprove{Approve: tx}, nil
	default:
	return nil, errors.New("operation type invalid")
	}
	}
	
	return qclitx.BroadcastTxAndPrintResult(cdc, iTxBuilder)
	}
