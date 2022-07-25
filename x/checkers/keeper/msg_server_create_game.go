package keeper

import (
	"context"
	"strconv"

	rules "github.com/alice/checkers/x/checkers/rules" //  готовый файл с правилами игры
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Handling the message (обработка сообщений создания игры)
func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// получение идентификатора новой игры
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	newIndex := strconv.FormatUint(nextGame.IdValue, 10)
	//объект для сохранения
	newGame := rules.New()
	storedGame := types.StoredGame{
		Creator: msg.Creator,
		Index:   newIndex,
		Game:    newGame.String(),
		Turn:    rules.PieceStrings[newGame.Turn],
		Red:     msg.Red,
		Black:   msg.Black,
	}
	// проверка адресов игроков
	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}
	// сохранение объекта
	k.Keeper.SetStoredGame(ctx, storedGame)

	// для следующей игры
	nextGame.IdValue++
	k.Keeper.SetNextGame(ctx, nextGame)

	// событие
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.StoredGameEventKey),
			sdk.NewAttribute(types.StoredGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.StoredGameEventIndex, newIndex),
			sdk.NewAttribute(types.StoredGameEventRed, msg.Red),
			sdk.NewAttribute(types.StoredGameEventBlack, msg.Black),
		),
	)

	// вернуть созданный идентификатор
	return &types.MsgCreateGameResponse{
		IdValue: newIndex,
	}, nil
}
