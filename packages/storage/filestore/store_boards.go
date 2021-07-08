package filestore

import (
	"context"
	"fmt"
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/storage"
	"os"
	"path"
)

var _ storage.BoardsStore = (*fsBoardsStore)(nil)

type fsBoardsStore struct {
	fsProjectStore
	boardsDirPath string
}

func (store fsBoardsStore) Project() storage.ProjectStore {
	return store.fsProjectStore
}

func newFsBoardsStore(fsProjectStore fsProjectStore) fsBoardsStore {
	return fsBoardsStore{
		fsProjectStore: fsProjectStore,
		boardsDirPath:  path.Join(fsProjectStore.projectPath, BoardsFolder),
	}
}

func (store fsBoardsStore) saveBoards(ctx context.Context, boards models.Boards) (err error) {
	return saveItems(BoardsFolder, len(boards), func(i int) func() error {
		return func() error {
			board := boards[i]
			_, err := store.SaveBoard(ctx, *board)
			return err
		}
	})
}

func (store fsBoardsStore) DeleteBoard(_ context.Context, id string) (err error) {
	filePath := path.Join(store.boardsDirPath, jsonFileName(id, boardFileSuffix))
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(filePath)
}

func (store fsBoardsStore) CreateBoard(ctx context.Context, board models.Board) (*models.Board, error) {
	panic("implement me")
}

func (store fsBoardsStore) SaveBoard(_ context.Context, board models.Board) (*models.Board, error) {
	if err := store.updateProjectFileWithBoard(board); err != nil {
		return &board, fmt.Errorf("failed to update project file with board: %w", err)
	}
	fileName := jsonFileName(board.ID, boardFileSuffix)
	board.ID = ""
	if err := saveJSONFile(
		store.boardsDirPath,
		fileName,
		board,
	); err != nil {
		return &board, fmt.Errorf("failed to save board file: %w", err)
	}
	return &board, nil
}

// GetBoard loads board
func (store fsBoardsStore) GetBoard(_ context.Context, id string) (*models.Board, error) {
	fileName := path.Join(store.boardsDirPath, fmt.Sprintf("%v.json", id))
	var board models.Board
	if err := readJSONFile(fileName, true, &board); err != nil {
		err = fmt.Errorf("faile to load board [%v] from project [%v]: %w", id, store.projectID, err)
		return nil, err
	}
	return &board, nil
}
