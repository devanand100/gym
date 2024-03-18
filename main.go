package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devanand100/gym/db"
	"github.com/devanand100/gym/server"
	_profile "github.com/devanand100/gym/server/profile"
	"github.com/devanand100/gym/store"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	mode    string
	addr    string
	profile *_profile.Profile
	port    int
	dbUri   string

	rootCmd = &cobra.Command{
		Use:   "gym",
		Short: "Gym managing app",
		Run: func(_cmd *cobra.Command, _args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			dbCtx, dbCancel := context.WithTimeout(ctx, 30*time.Second)

			defer dbCancel()

			DBInstance, err := db.NewDb(ctx, profile)

			defer DBInstance.Close(dbCtx)

			if err != nil {
				fmt.Println("error in  Database connection")
				cancel()
				return
			}

			storeInstance := store.New(DBInstance.GymDB, profile)

			s, err := server.NewServer(ctx, profile, storeInstance)

			if err != nil {
				cancel()
				log.Error("failed to create server", err)
				return
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			go func() {
				sig := <-c
				fmt.Sprintf("%s received.\n", sig.String())
				// s.Shutdown(ctx)
				cancel()
			}()

			if err := s.Start(ctx); err != nil {
				log.Error("failed to start server", err)
				cancel()
			}

			<-ctx.Done()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "demo", `mode of server, can be "prod" or "dev"`)
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "", "address of server")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8081, "port of server")
	rootCmd.PersistentFlags().StringVarP(&dbUri, "dbUri", "d", "mongodb://localhost:27017", "Database Connection Url")

	err := viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}

	viper.SetDefault("mode", "dev")
	viper.SetDefault("addr", "")
	viper.SetDefault("port", 8081)
	viper.SetDefault("dbUri", "mongodb://localhost:27017")
}

func main() {
	err := Execute()
	if err != nil {
		panic(err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	viper.AutomaticEnv()
	var err error
	profile, err = _profile.GetProfile()

	if err != nil {
		fmt.Printf("failed to get profile, error: %+v\n", err)
		return
	}

	println("---")
	println("Server profile")
	println("addr:", profile.Addr)
	println("port:", profile.Port)
	println("mode:", profile.Mode)
	println("db:", profile.DbUri)
	println("---")
}
