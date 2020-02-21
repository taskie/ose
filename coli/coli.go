package coli

import (
	"path/filepath"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/taskie/ose"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Coli struct {
	fs  afero.Fs
	io  ose.IO
	vpr *viper.Viper
}

func NewColi(fs afero.Fs, oio ose.IO, vpr *viper.Viper) *Coli {
	return &Coli{fs: fs, io: oio, vpr: vpr}
}

func NewColiInThisWorld() *Coli {
	w := ose.GetWorld()
	return NewColi(w.Fs(), w.IO(), viper.New())
}

func (c *Coli) Viper() *viper.Viper { return c.vpr }

func (c *Coli) Prepare(cmd *cobra.Command) {
	c.PrepareIO(cmd)
	c.PrepareFs(cmd)
	c.PrepareWellKnownFlags(cmd)
	c.PrepareConfig(cmd)
	c.PreparePreRun(cmd)
}

func (c *Coli) PrepareIO(cmd *cobra.Command) {
	cmd.SetIn(c.io.In())
	cmd.SetOut(c.io.Out())
	cmd.SetErr(c.io.Err())
}

func (c *Coli) PrepareFs(cmd *cobra.Command) {
	c.vpr.SetFs(c.fs)
}

func (c *Coli) PrepareWellKnownFlags(cmd *cobra.Command) {
	flg := cmd.PersistentFlags()
	flg.BoolP("verbose", "v", false, "verbose output")
	flg.Bool("debug", false, "debug output")
	flg.BoolP("version", "V", false, "show version")
	c.BindFlags(flg, []string{"verbose", "debug"})
}

func (c *Coli) PrepareConfig(cmd *cobra.Command) {
	v := c.vpr
	name := cmd.Use
	v.SetConfigName(name)
	v.AddConfigPath(".")
	configHome, err := ose.NewEnvPath(c.fs, ose.GetEnv()).GetXdgConfigHome()
	if err == nil {
		v.AddConfigPath(filepath.Join(configHome, name))
	}
	viper.SetEnvPrefix(name)
}

func (c *Coli) PreparePreRun(cmd *cobra.Command) {
	cmd.PreRun = c.PreRun
}

func (c *Coli) BindFlags(flg *pflag.FlagSet, names []string) {
	for _, s := range names {
		envKey := strcase.ToSnake(s)
		structKey := strcase.ToCamel(s)
		_ = c.vpr.BindPFlag(envKey, flg.Lookup(s))
		c.vpr.RegisterAlias(structKey, envKey)
	}
}

func newVerboseLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format(time.RFC3339))
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func newDebugLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cfg.DisableStacktrace = true
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func (c *Coli) PreRun(cmd *cobra.Command, args []string) {
	v := c.Viper()
	if v.GetBool("verbose") {
		zap.ReplaceGlobals(newVerboseLogger())
	} else if v.GetBool("debug") {
		zap.ReplaceGlobals(newDebugLogger())
	}
}

func (c *Coli) Execute(cmd *cobra.Command) error {
	v := c.Viper()
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		// nop
	}
	return cmd.Execute()
}

type ColiCommandProducer func(cl *Coli, name string, path []string) *cobra.Command
