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

func TestPauserQuery(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	item := createTestPauser(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPauserRequest
		response *types.QueryGetPauserResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetPauserRequest{},
			response: &types.QueryGetPauserResponse{Pauser: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Pauser(ctx, tc.request)
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
