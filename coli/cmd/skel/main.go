package main

import (
	"github.com/spf13/cobra"
	"github.com/taskie/ose/coli"
	"go.uber.org/zap"
)

const CommandName = "skel"

var Command *cobra.Command

func init() {
	Command = NewCommand(coli.NewColiInThisWorld())
}

func main() {
	Command.Execute()
}

func NewCommand(cl *coli.Coli) *cobra.Command {
	cmd := &cobra.Command{
		Use: CommandName,
		Run: cl.WrapRun(run),
	}
	cl.Prepare(cmd)

	cmd.Flags().StringP("from-type", "f", "", "convert from")
	cmd.Flags().StringP("to-type", "t", "", "convert to")

	cl.BindFlags(cmd.Flags(), []string{"from-type", "to-type"})
	return cmd
}

func run(cl *coli.Coli, cmd *cobra.Command, args []string) {
	v := cl.Viper()
	m := make(map[string]interface{})
	err := v.Unmarshal(&m)
	if err != nil {
		zap.L().Fatal("can't unmarshal", zap.Error(err))
	}
	cmd.Printf("flags: %v\n", m)
	cmd.Printf("args: %v\n", args)
}
