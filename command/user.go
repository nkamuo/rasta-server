package command

import (
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nkamuo/rasta-server/data/pagination"
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/service"
	"github.com/spf13/cobra"
)

func buildListUserCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		var users []model.User
		var page pagination.Page

		firstName, _ := command.Flags().GetString("firstname")
		lastName, _ := command.Flags().GetString("lastname")
		email, _ := command.Flags().GetString("email")
		phone, _ := command.Flags().GetString("phone")
		// password, _ := command.Flags().GetString("password")
		// isAdmin, _ := command.Flags().GetBool("is-admin")
		// published, _ := command.Flags().Get("published")

		query := model.DB.Model(&model.User{})

		if firstName != "" {
			arg := fmt.Sprintf("%%%s%%", firstName)
			query = query.Where("first_name LIKE ?", arg)
		}

		if lastName != "" {
			arg := fmt.Sprintf("%%%s%%", lastName)
			query = query.Where("last_name LIKE ?", arg)
		}

		if email != "" {
			arg := fmt.Sprintf("%%%s%%", email)
			query = query.Where("email LIKE ?", arg)
		}

		if phone != "" {
			arg := fmt.Sprintf("%%%s%%", phone)
			query = query.Where("phone LIKE ?", arg)
		}

		// if isAdmin != "" {
		// 	arg := fmt.Sprintf("%%%s%%", phone)
		// 	query = query.Where("phone LIKE ?", arg)
		// }

		if err := query.Scopes(pagination.Paginate(users, &page, query)).Find(&users).Error; err != err {
			message := fmt.Sprintf("Error running Query: %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
		}

		// fmt.Printf("DATA: %#v", users)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"", "Total", "", "", "", page.TotalRows})
		t.AppendHeader(table.Row{"#", "ID", "First Name", "Last Name", "Email", "Phone", "Is Admin"})

		var count = 0
		for _, user := range users {
			count++
			t.AppendRows([]table.Row{
				{count, user.ID, user.FirstName, user.LastName, user.Email, user.Phone, *user.IsAdmin},
			})
			t.AppendSeparator()
		}
		// t.AppendFooter(table.Row{"", "", "Total", 10000})
		t.Render()

	}
}

func buildCreateUserCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// username, _ := rootCmd.Flags().GetString("username")
		// password, _ := rootCmd.Flags().GetString("password")

		userService := service.GetUserService()

		firstName, _ := command.Flags().GetString("firstname")
		lastName, _ := command.Flags().GetString("lastname")
		email, _ := command.Flags().GetString("email")
		phone, _ := command.Flags().GetString("phone")
		password, _ := command.Flags().GetString("password")
		isAdmin, _ := command.Flags().GetBool("is-admin")
		published, _ := command.Flags().GetBool("published")

		nameRegexp, err := regexp.Compile(`^[\w'\-,.][^0-9_!¡?÷?¿/\\+=@#$%ˆ&*(){}|~<>;:[\]]{2,}$`)
		if err != nil {
			message := fmt.Sprintf("Error Compiling Regexp %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
		}

		emailRegexp, err := regexp.Compile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$`)
		if err != nil {
			message := fmt.Sprintf("Error Compiling Regexp %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
		}

		if !nameRegexp.MatchString(firstName) {
			message := fmt.Sprintf("First name not Valid: %s", firstName)
			fmt.Println(message)
			os.Exit(1)
		}

		if !nameRegexp.MatchString(lastName) {
			message := fmt.Sprintf("Last name not Valid: %s", lastName)
			fmt.Println(message)
			os.Exit(1)
		}

		if !emailRegexp.MatchString(email) {
			message := fmt.Sprintf("email not Valid: %s", email)
			fmt.Println(message)
			os.Exit(1)
		}

		user := model.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Phone:     phone,
			IsAdmin:   &isAdmin,
			Published: published,
		}

		err = userService.HashUserPassword(&user, password)
		if err != nil {
			message := fmt.Sprintf("Error hashing user password: %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
			panic(err)
		}

		err = userService.Save(&user)
		if err != nil {
			message := fmt.Sprintf("Error saving user: %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
		}

	}
}

func buildRemoveUserCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		userService := service.GetUserService()
		// username, _ := rootCmd.Flags().GetString("username")
		id, _ := command.Flags().GetString("id")
		ID, err := uuid.Parse(id)
		if err != nil {
			message := fmt.Sprintf("Invalid UUID[%s]: %s", id, err.Error())
			fmt.Println(message)
			os.Exit(1)
			return
		}

		if user, err := userService.GetById(ID, "Password"); err != nil {
			message := fmt.Sprintf("Could not fetch user: %s", err.Error())
			fmt.Println(message)
			os.Exit(1)
			return
		} else {
			if err := userService.Delete(user); err != nil {
				message := fmt.Sprintf("Could not delete user: %s", err.Error())
				fmt.Println(message)
				os.Exit(1)
				return
			}
			message := fmt.Sprintf("Removed user \"%s\"!", user.FullName())
			fmt.Println(message)
		}

	}
}

// createUserCommand to start the rasta webserver
var userCommand = &cobra.Command{
	Use:   "user",
	Short: "Commands to manage the users",
	Long:  `Manage Users`,
}

var listUserCommand = &cobra.Command{
	Use:   "list",
	Short: "List Users",
	Long:  `List all user's in the system`,
}

// createUserCommand to start the rasta webserver
var createUserCommand = &cobra.Command{
	Use:   "create",
	Short: "Create new User",
	Long:  `Create a new user account record in the system`,
}

var removeUserCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes a User",
	Long:  `Removes a User from the user list`,
}

func init() {

	listUserCommand.Run = buildListUserCommand(listUserCommand)
	createUserCommand.Run = buildCreateUserCommand(createUserCommand)
	removeUserCommand.Run = buildRemoveUserCommand(removeUserCommand)

	rootCmd.AddCommand(userCommand)
	userCommand.AddCommand(listUserCommand)
	userCommand.AddCommand(createUserCommand)
	userCommand.AddCommand(removeUserCommand)

	listUserCommand.PersistentFlags().StringP("search", "s", "", "String to search against the user")
	listUserCommand.PersistentFlags().StringP("email", "m", "", "User's E-mail address")
	listUserCommand.PersistentFlags().StringP("phone", "t", "", "User's phone number")
	listUserCommand.PersistentFlags().StringP("firstname", "f", "", "User's First name")
	listUserCommand.PersistentFlags().StringP("lastname", "l", "", "User's Last name")
	listUserCommand.PersistentFlags().StringP("password", "p", "", "User's Password")
	listUserCommand.PersistentFlags().BoolP("is-admin", "a", false, "If this user is an admin or not")
	listUserCommand.PersistentFlags().BoolP("published", "o", false, "If this user is published")

	//

	createUserCommand.PersistentFlags().StringP("email", "m", "", "User's E-mail address")
	createUserCommand.PersistentFlags().StringP("phone", "t", "", "User's phone number")
	createUserCommand.PersistentFlags().StringP("firstname", "f", "", "User's First name")
	createUserCommand.PersistentFlags().StringP("lastname", "l", "", "User's Last name")
	createUserCommand.PersistentFlags().StringP("password", "p", "", "User's Password")
	createUserCommand.PersistentFlags().BoolP("is-admin", "a", false, "If this user is an admin or not")
	createUserCommand.PersistentFlags().BoolP("published", "o", false, "If this user is published")

	removeUserCommand.PersistentFlags().String("id", "i", "The ID of the user account to remove")
}
