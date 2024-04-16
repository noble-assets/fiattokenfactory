package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMintingDenomQuery(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	item := createTestMintingDenom(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetMintingDenomRequest
		response *types.QueryGetMintingDenomResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetMintingDenomRequest{},
			response: &types.QueryGetMintingDenomResponse{MintingDenom: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.MintingDenom(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					utils.Fill(tc.response),
					utils.Fill(response),
				)
			}
		})
	}
}
