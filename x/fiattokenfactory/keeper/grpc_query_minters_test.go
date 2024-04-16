package keeper_test

import (
	"strconv"
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestMintersQuerySingle(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	msgs := createNMinters(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetMintersRequest
		response *types.QueryGetMintersResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetMintersRequest{
				Address: msgs[0].Address,
			},
			response: &types.QueryGetMintersResponse{Minters: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetMintersRequest{
				Address: msgs[1].Address,
			},
			response: &types.QueryGetMintersResponse{Minters: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetMintersRequest{
				Address: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Minters(ctx, tc.request)
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

func TestMintersQueryPaginated(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	msgs := createNMinters(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllMintersRequest {
		return &types.QueryAllMintersRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.MintersAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Minters), step)
			require.Subset(t,
				utils.Fill(msgs),
				utils.Fill(resp.Minters),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.MintersAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Minters), step)
			require.Subset(t,
				utils.Fill(msgs),
				utils.Fill(resp.Minters),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.MintersAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			utils.Fill(msgs),
			utils.Fill(resp.Minters),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.MintersAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
