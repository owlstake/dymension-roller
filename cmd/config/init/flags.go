package initconfig

import (
	"fmt"
	"regexp"

	"github.com/dymensionxyz/roller/cmd/consts"
	"github.com/dymensionxyz/roller/cmd/utils"
	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(FlagNames.HubID, "", TestnetHubID, fmt.Sprintf("The ID of the Dymension hub. %s", getAvailableHubsMessage()))
	cmd.Flags().StringP(FlagNames.RollappBinary, "", "", "The rollapp binary. Should be passed only if you built a custom rollapp")
	cmd.Flags().Uint64P(FlagNames.Decimals, "", 18, "The number of decimal places a rollapp token supports")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		hubID, err := cmd.Flags().GetString(FlagNames.HubID)
		if err != nil {
			return err
		}
		if _, ok := Hubs[hubID]; !ok {
			return fmt.Errorf("invalid hub ID: %s. %s", hubID, getAvailableHubsMessage())
		}
		rollappID := args[0]
		if !validateRollAppID(rollappID) {
			return fmt.Errorf("invalid RollApp ID '%s'. %s", rollappID, getValidRollappIdMessage())
		}
		return nil
	}
}

func getValidRollappIdMessage() string {
	return "A valid RollApp ID should follow the format 'rollapp-name_EIP155_version', where 'rollapp-name' is made up of" +
		" lowercase English letters, 'EIP155_version' is a 1 to 5 digit number representing the EIP155 rollapp ID, and '" +
		"version' is a 1 to 5 digit number representing the version. For example: 'mars_9721_1'"
}

func getDecimals(cmd *cobra.Command) uint64 {
	decimals, err := cmd.Flags().GetUint64(FlagNames.Decimals)
	if err != nil {
		panic(err)
	}
	return decimals
}

func getRollappBinaryPath(cmd *cobra.Command) string {
	rollappBinaryPath := cmd.Flag(FlagNames.RollappBinary).Value.String()
	if rollappBinaryPath == "" {
		return consts.Executables.RollappEVM
	}
	return rollappBinaryPath
}

func GetInitConfig(initCmd *cobra.Command, args []string) InitConfig {
	rollappId := args[0]
	denom := args[1]
	home := initCmd.Flag(utils.FlagNames.Home).Value.String()
	rollappBinaryPath := getRollappBinaryPath(initCmd)
	decimals := getDecimals(initCmd)
	hubID := initCmd.Flag(FlagNames.HubID).Value.String()
	return InitConfig{
		Home:          home,
		RollappID:     rollappId,
		RollappBinary: rollappBinaryPath,
		Denom:         denom,
		Decimals:      decimals,
		HubData:       Hubs[hubID],
	}
}

func getAvailableHubsMessage() string {
	return fmt.Sprintf("Acceptable values are '%s', '%s' or '%s'", TestnetHubID, StagingHubID, LocalHubID)
}

func validateRollAppID(id string) bool {
	pattern := `^[a-z]+_[0-9]{1,5}_[0-9]{1,5}$`
	r, _ := regexp.Compile(pattern)
	return r.MatchString(id)
}