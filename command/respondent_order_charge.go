package command

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/spf13/cobra"
)

func buildListRespondentChargesCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		var charges []model.RespondentOrderCharge
		var page pagination.Page

		label, _ := command.Flags().GetString("label")
		status, _ := command.Flags().GetStringArray("status")
		// password, _ := command.Flags().GetString("password")
		// isAdmin, _ := command.Flags().GetBool("is-admin")
		// published, _ := command.Flags().Get("published")

		query := model.DB.Model(&model.RespondentOrderCharge{})

		if label != "" {
			arg := fmt.Sprintf("%%%s%%", label)
			query = query.Where("label LIKE ?", arg)
		}

		if status != nil && len(status) > 0 {
			query = query.Where("status in ?", status)
		}

		if err := query.Scopes(pagination.Paginate(charges, &page, query)).Find(&charges).Error; err != err {
			message := fmt.Sprintf("Error running Query: %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
		}

		// fmt.Printf("DATA: %#v", charges)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"", "Total", "", "", page.TotalRows})
		t.AppendHeader(table.Row{"#", "ID", "Label", "Status", "Created At"})

		var count = 0
		for _, charge := range charges {
			count++
			t.AppendRows([]table.Row{
				{count, charge.ID, charge.Label, charge.Status, charge.CreatedAt},
			})
			t.AppendSeparator()
		}
		// t.AppendFooter(table.Row{"", "", "Total", 10000})
		t.Render()

	}
}

func buildUpdateAllRespondentOrderChargeCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		chargeService := service.GetRespondentOrderChargeService()
		// chargename, _ := rootCmd.Flags().GetString("chargename")
		// status, _ := command.Flags().GetStringArray("status")
		chargeService.UpdateAllCharges()

	}
}

func buildRemoveRespondentOrderChargeCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		chargeService := service.GetRespondentOrderChargeService()
		// chargename, _ := rootCmd.Flags().GetString("chargename")
		id, _ := command.Flags().GetString("id")
		ID, err := uuid.Parse(id)
		if err != nil {
			message := fmt.Sprintf("Invalid UUID[%s]: %s", id, err.Error())
			fmt.Println(message)
			os.Exit(1)
			return
		}

		if charge, err := chargeService.GetById(ID); err != nil {
			message := fmt.Sprintf("Could not fetch charge: %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
			return
		} else {
			if err := chargeService.Delete(charge); err != nil {
				message := fmt.Sprintf("Could not delete charge: %s", err.Error())
				fmt.Println(message)
				os.Exit(1)
				return
			}
			message := fmt.Sprintf("Removed charge \"%s\"!", charge.Label)
			fmt.Println(message)
		}

	}
}

// createRespondentOrderChargeCommand to start the rasta webserver
var chargeCommand = &cobra.Command{
	Use:   "charge",
	Short: "Commands to manage the charges",
	Long:  `Manage RespondentOrderCharges`,
}

var listRespondentOrderChargeCommand = &cobra.Command{
	Use:   "list",
	Short: "List RespondentOrderCharges",
	Long:  `List all charge's in the system`,
}

// // createRespondentOrderChargeCommand to start the rasta webserver
// var createRespondentOrderChargeCommand = &cobra.Command{
// 	Use:   "create",
// 	Short: "Create new RespondentOrderCharge",
// 	Long:  `Create a new charge account record in the system`,
// }

var removeRespondentOrderChargeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes a RespondentOrderCharge",
	Long:  `Removes a RespondentOrderCharge from the charge list`,
}

var updateAllRespondentOrderChargeCommand = &cobra.Command{
	Use:   "upate-all",
	Short: "Updates all RespondentOrderCharge",
	Long:  `Checks with the Stripe API for the current status of this charge`,
}

func init() {

	listRespondentOrderChargeCommand.Run = buildListRespondentChargesCommand(listRespondentOrderChargeCommand)
	removeRespondentOrderChargeCommand.Run = buildRemoveRespondentOrderChargeCommand(removeRespondentOrderChargeCommand)
	updateAllRespondentOrderChargeCommand.Run = buildUpdateAllRespondentOrderChargeCommand(updateAllRespondentOrderChargeCommand)

	rootCmd.AddCommand(chargeCommand)
	chargeCommand.AddCommand(listRespondentOrderChargeCommand)
	chargeCommand.AddCommand(removeRespondentOrderChargeCommand)
	chargeCommand.AddCommand(updateAllRespondentOrderChargeCommand)

	listRespondentOrderChargeCommand.PersistentFlags().StringP("search", "s", "", "String to search against the charge")
	listRespondentOrderChargeCommand.PersistentFlags().StringArrayP("label", "l", []string{"pending"}, "RespondentOrderCharge's status to include")
	//
	updateAllRespondentOrderChargeCommand.PersistentFlags().StringArrayP("status", "l", []string{"pending"}, "RespondentOrderCharge's status to include")

	//

	// createRespondentOrderChargeCommand.PersistentFlags().StringP("email", "m", "", "RespondentOrderCharge's E-mail address")
	// createRespondentOrderChargeCommand.PersistentFlags().StringP("phone", "t", "", "RespondentOrderCharge's phone number")
	// createRespondentOrderChargeCommand.PersistentFlags().StringP("firstname", "f", "", "RespondentOrderCharge's First name")
	// createRespondentOrderChargeCommand.PersistentFlags().StringP("lastname", "l", "", "RespondentOrderCharge's Last name")
	// createRespondentOrderChargeCommand.PersistentFlags().StringP("password", "p", "", "RespondentOrderCharge's Password")
	// createRespondentOrderChargeCommand.PersistentFlags().BoolP("is-admin", "a", false, "If this charge is an admin or not")
	// createRespondentOrderChargeCommand.PersistentFlags().BoolP("published", "o", false, "If this charge is published")

	removeRespondentOrderChargeCommand.PersistentFlags().String("id", "i", "The ID of the charge account to remove")
}
