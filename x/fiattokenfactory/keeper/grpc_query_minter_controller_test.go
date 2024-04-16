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

func TestMinterControllerQuerySingle(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	msgs := createNMinterController(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetMinterControllerRequest
		response *types.QueryGetMinterControllerResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetMinterControllerRequest{
				ControllerAddress: msgs[0].Controller,
			},
			response: &types.QueryGetMinterControllerResponse{MinterController: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetMinterControllerRequest{
				ControllerAddress: msgs[1].Controller,
			},
			response: &types.QueryGetMinterControllerResponse{MinterController: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetMinterControllerRequest{
				ControllerAddress: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.MinterController(ctx, tc.request)
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

func TestMinterControllerQueryPaginated(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	msgs := createNMinterController(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllMinterControllerRequest {
		return &types.QueryAllMinterControllerRequest{
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
			resp, err := keeper.MinterControllerAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MinterController), step)
			require.Subset(t,
				utils.Fill(msgs),
				utils.Fill(resp.MinterController),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.MinterControllerAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MinterController), step)
			require.Subset(t,
				utils.Fill(msgs),
				utils.Fill(resp.MinterController),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.MinterControllerAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			utils.Fill(msgs),
			utils.Fill(resp.MinterController),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.MinterControllerAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
