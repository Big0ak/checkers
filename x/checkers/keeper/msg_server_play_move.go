package keeper

import (
	"context"
	"strconv"
	"strings"

	rules "github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Handling the message обработка сообщений возможность сделать ход
func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// получить сохраненную игровую информацию
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "game not found %s", msg.IdValue)
	}

	// проверяем правильность игрока
	isRed := strings.Compare(storedGame.Red, msg.Creator) == 0
	isBlack := strings.Compare(storedGame.Black, msg.Creator) == 0
	var player rules.Player
	if !isRed && !isBlack {
		return nil, types.ErrCreatorNotPlayer
	} else if isRed && isBlack {
		player = rules.StringPieces[storedGame.Turn].Player
	} else if isRed {
		player = rules.RED_PLAYER
	} else {
		player = rules.BLACK_PLAYER
	}

	// создание экземпляра доски, чтобы реализовать правила
	game, err := storedGame.ParseGame()
	if err != nil {
		panic(err.Error())
	}
	//Is it the player's turn?
	if !game.TurnIs(player) {
		return nil, types.ErrNotPlayerTurn
	}

	// правильно выполнить ход используя правила Move функции
	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	// сохранение изменений
	storedGame.Game = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]
	k.Keeper.SetStoredGame(ctx, storedGame)

	// событие
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.PlayMoveEventKey),
			sdk.NewAttribute(types.PlayMoveEventCreator, msg.Creator),
			sdk.NewAttribute(types.PlayMoveEventIdValue, msg.IdValue),
			sdk.NewAttribute(types.PlayMoveEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.PlayMoveEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.PlayMoveEventWinner, rules.PieceStrings[game.Winner()]),
		),
	)

	// возвращает информацию о результате хода
	return &types.MsgPlayMoveResponse{
		IdValue:   msg.IdValue,
		CapturedX: int64(captured.X),
		CapturedY: int64(captured.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil

}
