package testing

import (
	"fmt"

	sdkchannel "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

func NewClient(cfgPath, user, secret, org, channel string) error {
	sdk, err := fabsdk.New(config.FromFile(cfgPath))
	if err != nil {
		return err
	}

	mspc, err := msp.New(sdk.Context(), msp.WithOrg(org))
	if err != nil {
		return err
	}

	if err = mspc.Enroll(user, msp.WithSecret(secret)); err != nil {
		return err
	}

	if _, err = mspc.GetSigningIdentity(user); err != nil {
		return errors.WithMessage(err, "unable to get signing identity")
	}

	clientContext := sdk.ChannelContext(
		channel,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org),
	)

	channelClient, err := sdkchannel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "unable to create client from channel context")
	}

	
	resp, err := channelClient.Query(
		sdkchannel.Request{
			ChaincodeID: "resources",
			Fcn:         "Index",
			Args:        [][]byte{[]byte(``)},
		},
		sdkchannel.WithTargetEndpoints("peer0-"+org),
		sdkchannel.WithRetry(retry.DefaultChannelOpts),
	)
	if err != nil {
		return err
	}

	fmt.Println(string(resp.Payload))

	return nil
}
