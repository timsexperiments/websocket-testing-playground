package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tim117/chatserver/chat"
	"gorm.io/gorm"
)

var user *chat.User
var username string
var verbose = false
var userRepo chat.UserRepo
var chatRepo chat.ChatRepo

func Execute(db *gorm.DB) {
	userRepo = chat.UserRepository(db)
	chatRepo = chat.ChatRepository(db)

	root := baseCmd()
	root.AddCommand(signInCommand())
	root.AddCommand(createCommand())
	cobra.CheckErr(root.Execute())
}

func baseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Chat CLI made in GO.",
		Long:  "Chat server CLI that is written in GO.",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if user != nil {
				return
			}
			if username != "" {
				fmt.Printf("Username [%s] was provided. Attempting to sign in.", username)
				findUser(username)
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&username, "username", "u", "", "The username of the user to complete operations with.")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Whether to be verbose about what the CLI is doing. Good for debugging.")

	return cmd
}

func signInCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sign-in [username]",
		Short: "This command signs in a user",
		Long:  "The sign-in command signs in a user with the given username.",
		Run: func(_ *cobra.Command, args []string) {
			username := args[0]
			findUser(username)
		},
	}
}

func createCommand() *cobra.Command {
	args := []string{"chat", "message", "user"}
	return &cobra.Command{
		Use:       "create chat | message | user []",
		ValidArgs: args,
		Short:     "This command creates an entity in the DB.",
		Long:      "The create command creates a new entity of the specified type in the database.",
		Run:       func(_ *cobra.Command, args []string) {},
	}
}

func findUser(username string) {
	if verbose {
		fmt.Printf("Signing in %s.\n", username)
	}
	if username == "" {
		panic("No username was provided but is required.")
	}
	found, err := userRepo.GetByUsername(username)
	if err != nil {
		panic(err.Error())
	}
	user = found
	if verbose {
		fmt.Printf("%s successfully signed in.\n", username)
	}
}
