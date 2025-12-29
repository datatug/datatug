package api

import (
	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug-core/pkg/dto"
	"github.com/datatug/datatug/pkg/sqlexecute"
	"github.com/strongo/validation"
)

// GetServerDatabases returns list of databases hosted at a server
func GetServerDatabases(request dto.GetServerDatabasesRequest) (databases []*datatug.EnvDbCatalog, err error) {
	if err = request.Validate(); err != nil {
		return nil, validation.NewBadRequestError(err)
	}

	executor := sqlexecute.NewExecutor(nil, nil)

	command := sqlexecute.RequestCommand{
		ServerReference: request.ServerReference,
		Text:            "select name from sys.databases where owner_sid > 0x01",
	}
	var response sqlexecute.Response
	if response, err = executor.ExecuteSingle(command); err != nil {
		return nil, err
	}
	recordset := response.Commands[0].Items[0].Value.(datatug.Recordset)
	databases = make([]*datatug.EnvDbCatalog, len(recordset.Rows))
	for i, row := range recordset.Rows {
		databases[i] = new(datatug.EnvDbCatalog)
		databases[i].ID = row[0].(string)
	}
	return
}
