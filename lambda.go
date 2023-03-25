package main

import (
	"context"
	"flag"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-acme/lego/v4/cmd"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func main() {
	lambda.Start(issue)
}

func issue(ctx context.Context, in any) (out any, err error) {
	PrintAsJSON(in)
	cmds := cmd.CreateCommands()
	issue := lo.Filter(cmds, func(c *cli.Command, _ int) bool {
		return c.Name == "run"
	})[0]

	set := &flag.FlagSet{}
	for _, v := range append(cmd.CreateFlags(os.Getenv("LEGO_PATH")), issue.Flags...) {
		if err = v.Apply(set); err != nil {
			return
		}
	}
	if err = set.Set("email", os.Getenv("EMAIL")); err != nil {
		return
	}
	if err = set.Set("domains", os.Getenv("DOMAINS")); err != nil {
		return
	}
	if err = set.Set("dns", os.Getenv("DNS_SERVER")); err != nil {
		return
	}
	if err = set.Set("server", os.Getenv("CA_DIRECTORY")); err != nil {
		return
	}
	if err = set.Set("eab", "true"); err != nil {
		return
	}

	cmdCtx := cli.NewContext(&cli.App{}, set, nil)
	cmdCtx.Context = ctx

	if err = issue.Run(cmdCtx); err != nil {
		return
	}
	out = map[string]any{
		"msg": "Success",
	}
	return
}
