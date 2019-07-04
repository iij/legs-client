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

type sendOptions struct {
	verboseFlag bool
}

// newSendCmd creates the send command.
func newSendCmd() *cobra.Command {
	o := &sendOptions{}

	cmd := &cobra.Command{
		Use:   "send <routing name or URL> <string for send server>",
		Short: "Send to server by routing name and message strings.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			o.sendMessage(conf.GetString("sock_file"), args)
		},
	}

	cmd.Flags().BoolVarP(&o.verboseFlag, "verbose", "v", false, "verbose")
	return cmd
}

func init() {
	rootCmd.AddCommand(newSendCmd())
}

func (o *sendOptions) sendMessage(sockFile string, args []string) {
	cli := golum.Client{
		SocketName: sockFile,
		Timeout:    30 * time.Second,
	}
	data := dto.NewSendRequest(args[0], args[1:])

	ret, err := sendRequest(&cli, data)
	if o.verboseFlag {
		fmt.Println(ret)
	}
	if err != nil {
		os.Exit(1)
	}
}

func sendRequest(c *golum.Client, data *dto.SendRequest) (string, error) {
	b, err := data.ToBinary()
	if err != nil {
		return err.Error(), err
	}
	req := golum.NewRequest("send", b)
	res, err := c.Do(req)
	if err != nil {
		return err.Error(), err
	}
	ret, err := dto.ParseSendResponse(res.Body)
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
