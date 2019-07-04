package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/iij/legs-client/daemon/golum"
	"github.com/iij/legs-client/daemon/handler/socket/dto"
	"github.com/spf13/cobra"
)

type getOptions struct {
	verboseFlag bool
	dryRunFlag  bool
}

// newGetCmd creates the get command.
func newGetCmd() *cobra.Command {
	o := &getOptions{}

	cmd := &cobra.Command{
		Use:   "get <routing name or URL> <download path>",
		Short: "Get from server by routing name/URL.",
		Args:  cobra.ExactArgs(2),
		Run: func(c *cobra.Command, args []string) {
			o.getMessage(conf.GetString("sock_file"), args)
		},
	}

	cmd.Flags().BoolVarP(&o.verboseFlag, "verbose", "v", false, "verbose")
	cmd.Flags().BoolVarP(&o.dryRunFlag, "dryrun", "n", false, "dryrun")
	return cmd
}

func init() {
	rootCmd.AddCommand(newGetCmd())
}

func (o *getOptions) getMessage(sockFile string, args []string) {
	cli := golum.Client{
		SocketName: sockFile,
		Timeout:    30 * time.Second,
	}
	reqOpt := dto.GetOptions{IsDryRun: o.dryRunFlag}
	data := dto.NewGetRequest(args[0], args[1], reqOpt)

	ret, err := getRequest(&cli, data)
	if o.verboseFlag {
		fmt.Println(ret)
	}
	if err != nil {
		os.Exit(1)
	}
}

func getRequest(c *golum.Client, data *dto.GetRequest) (string, error) {
	b, err := data.ToBinary()
	if err != nil {
		return err.Error(), err
	}
	req := golum.NewRequest("get", b)
	res, err := c.Do(req)
	if err != nil {
		return err.Error(), err
	}
	ret, err := dto.ParseGetResponse(res.Body)
	if res.Code != golum.StatusSuccess {
		if err != nil {
			err = errors.New(string(res.Body))
		} else {
			err = errors.New(ret.String())
		}
		return err.Error(), err
	}
	return ret.String(), nil
}
