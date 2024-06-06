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

func TestOwnerQuery(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()

	_, err := keeper.Owner(ctx, &types.QueryGetOwnerRequest{})
	require.Error(t, err, codes.NotFound)

	owner := types.Owner{Address: "test"}
	keeper.SetOwner(ctx, owner)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetOwnerRequest
		response *types.QueryGetOwnerResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetOwnerRequest{},
			response: &types.QueryGetOwnerResponse{Owner: owner},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Owner(ctx, tc.request)
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
