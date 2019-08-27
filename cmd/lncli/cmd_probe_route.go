// +build routerrpc

package main

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"

	"github.com/urfave/cli"
)

var probeRouteCommand = cli.Command{
	Name:     "proberoute",
	Category: "Payments",
	Usage:    "Probe a route made up of channel ids.",
	Action:   actionDecorator(probeRoute),
	Flags: []cli.Flag{
		cli.Int64Flag{
			Name:  "amt",
			Usage: "the amount to send expressed in satoshis",
		},
		cli.Int64Flag{
			Name: "final_cltv_delta",
			Usage: "number of blocks the last hop has to reveal " +
				"the preimage",
			Value: 40,
		},
		cli.StringFlag{
			Name:  "chanids",
			Usage: "comma separated channel ids",
		},
		cli.Int64Flag{
			Name:  "timeout",
			Usage: "per-hop timeout in seconds",
			Value: 30,
		},
	},
}

func probeRoute(ctx *cli.Context) error {
	conn := getClientConn(ctx, false)
	defer conn.Close()

	client := routerrpc.NewRouterClient(conn)

	if !ctx.IsSet("amt") {
		return errors.New("amt required")
	}
	if !ctx.IsSet("chanids") {
		return errors.New("chanids required")
	}

	chanIDs := strings.Split(ctx.String("chanids"), ",")
	rpcChanIDs := make([]uint64, 0, len(chanIDs))
	for _, k := range chanIDs {
		chanID, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			return err
		}
		rpcChanIDs = append(rpcChanIDs, chanID)
	}

	req := &routerrpc.ProbeRouteRequest{
		TotalAmtMsat:   ctx.Int64("amt") * 1000,
		FinalCltvDelta: int32(ctx.Int64("final_cltv_delta")),
		RouteChannels:  rpcChanIDs,
		PerHopTimeout:  int32(ctx.Int64("timeout")),
	}

	rpcCtx := context.Background()
	failure, err := client.ProbeRoute(rpcCtx, req)
	if err != nil {
		return err
	}

	printJSON(failure)

	return nil
}
