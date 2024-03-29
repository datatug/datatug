package models

import (
	"encoding/json"
	"fmt"
	"github.com/strongo/validation"
)

// Boards is a slice of *Board
type Boards []*Board

// Validate returns error if failed
func (v Boards) Validate() error {
	for i, item := range v {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("invalid board at index %v: %w", i, err)
		}
	}
	return nil
}

// Board is holding all details about board
type Board struct {
	ProjBoardBrief
	Rows BoardRows `json:"rows,omitempty" firestore:"rows,omitempty"`
}

// Validate returns error if failed
func (v Board) Validate() error {
	if err := v.ProjBoardBrief.Validate(); err != nil {
		return err
	}
	if err := v.Rows.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("rows", err.Error())
	}
	return nil
}

// ProjBoardBrief defines brief information of Board
type ProjBoardBrief struct {
	ProjItemBrief
	Parameters Parameters `json:"parameters,omitempty" firestore:"parameters,omitempty"`
	// TODO: Document why and how to use, add tests
	RequiredParams [][]string `json:"requiredParams,omitempty" firestore:"requiredParams,omitempty"`
}

func (v ProjBoardBrief) Validate() error {
	if err := v.ProjItemBrief.Validate(true); err != nil {
		return err
	}
	for i, p := range v.Parameters {
		if err := p.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("parameters", fmt.Sprintf("invalid parameter #%v: %v", i+1, err))
		}
	}
	return nil
}

//func validateBoardBriefsMappedByID(boards map[string]*ProjBoardBrief) error { // TODO: generic
//	for id, board := range boards {
//		if err := validateItemMappedByID(id, board.ID, board); err != nil {
//			return err
//		}
//	}
//	return nil
//}

// ParameterLookup holds definition for parameter lookup
type ParameterLookup struct {
	DB        string   `json:"db" firestore:"db"`
	SQL       string   `json:"sql" firestore:"sql"`
	KeyFields []string `json:"keyFields" firestore:"keyFields"`
}

// BoardRows is slice of `BoardRow`
type BoardRows []BoardRow

// Validate returns error if failed
func (v BoardRows) Validate() error {
	for i, row := range v {
		if err := row.Validate(); err != nil {
			return fmt.Errorf("invalid row at index %v: %w", i, err)
		}
	}
	return nil
}

// BoardRow holds all details about a row in a board
type BoardRow struct {
	MinHeight string     `json:"minHeight,omitempty"`
	MaxHeight string     `json:"maxHeight,omitempty"`
	Cards     BoardCards `json:"cards,omitempty"`
}

// Validate returns error if failed
func (v BoardRow) Validate() error {
	// TODO: test MinHeight & MaxHeight
	if err := v.Cards.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("cards", err.Error())
	}
	return nil
}

// BoardCards is slice of BoardCard
type BoardCards []BoardCard

// Validate returns error if not valid
func (v BoardCards) Validate() error {
	for i, card := range v {
		if err := card.Validate(); err != nil {
			return fmt.Errorf("invalid card at index %v: %w", i, err)
		}
	}
	return nil
}

// BoardCard describes board card
type BoardCard struct {
	ID     string       `json:"id"`
	Title  string       `json:"title"`
	Cols   int          `json:"cols,omitempty"`
	Widget *BoardWidget `json:"widget,omitempty"`
}

// Validate returns error if failed
func (v BoardCard) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if v.Cols < 0 {
		return validation.NewErrBadRecordFieldValue("cols", "should be positive")
	}
	if err := v.Widget.Validate(); err != nil {
		return fmt.Errorf("invalid card widget with id=%v: %w", v.ID, err)
	}
	return nil
}

// BoardWidget specifies widget. Some widgets can contain otter widgets.
type BoardWidget struct {
	Name string      `json:"name"`
	Data interface{} `json:"data,omitempty"`
}

// Validate returns error if failed
func (v BoardWidget) Validate() (err error) {
	var widget interface{ Validate() error }
	var ok bool
	switch v.Name {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("name")
	case "SQL":
		if widget, ok = v.Data.(*SQLWidgetDef); !ok {
			if widget, err = newWidgetDef(&SQLWidgetDef{}, v.Data); err != nil {
				return err
			}
		}
	case "HTTP":
		if widget, ok = v.Data.(*HTTPWidgetDef); !ok {
			if widget, err = newWidgetDef(&HTTPWidgetDef{}, v.Data); err != nil {
				return err
			}
		}
	case "tabs":
		if widget, ok = v.Data.(*TabsWidgetDef); !ok {
			if widget, err = newWidgetDef(&TabsWidgetDef{}, v.Data); err != nil {
				return err
			}
		}
	default:
		return validation.NewErrBadRecordFieldValue("id", "unknown widget id")
	}
	if err = widget.Validate(); err != nil {
		return fmt.Errorf("failed to test widget of type %T: %w", widget, err)
	}
	return nil
}

func newWidgetDef(widgetDef validatable, data interface{}) (validatable, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return widgetDef, json.Unmarshal(d, widgetDef)
}

// WidgetBase is base struct for widgets
type WidgetBase struct {
	Title      string     `json:"title,omitempty"`
	Parameters Parameters `json:"parameters,omitempty"`
}

// Validate returns error if not valid
func (v WidgetBase) Validate(isTitleRequired bool) error {
	if err := validateStringField("title", v.Title, isTitleRequired, MaxTitleLength); err != nil {
		return err
	}
	if err := v.Parameters.Validate(); err != nil {
		return err
	}
	return nil
}

// SQLWidgetDef holds info about a widget that makes SQL queries
type SQLWidgetDef struct {
	WidgetBase
	SQL SQLWidgetSettings `json:"sql"`
}

// SQLWidgetSettings holds settings for an SQL widget
type SQLWidgetSettings struct {
	Query string `json:"query"`
}

// Validate returns error if not valid
func (v SQLWidgetSettings) Validate() error {
	if v.Query == "" {
		return validation.NewErrRecordIsMissingRequiredField("query")
	}
	return nil
}

// Validate returns error if not valid
func (v *SQLWidgetDef) Validate() error {
	if err := v.WidgetBase.Validate(false); err != nil {
		return err
	}
	if err := v.SQL.Validate(); err != nil {
		return err
	}
	return nil
}

// HTTPHeaders is a []HTTPHeaderItem
type HTTPHeaders []HTTPHeaderItem

// Validate returns error if not valid
func (v HTTPHeaders) Validate() error {
	for i, item := range v {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("invalid header item at index %v: %w", i, err)
		}
	}
	return nil
}

// HTTPHeaderItem describes an HTTP header item
type HTTPHeaderItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Validate returns error if not valid
func (v HTTPHeaderItem) Validate() error {
	if v.Name == "" {
		return validation.NewErrRecordIsMissingRequiredField("name")
	}
	if v.Value == "" {
		return validation.NewErrRecordIsMissingRequiredField("value")
	}
	return nil
}

// HTTPRequest describes an HTTP request
type HTTPRequest struct {
	Method             string      `json:"method"`
	URL                string      `json:"url"`
	Protocol           string      `json:"protocol,omitempty"`
	Headers            HTTPHeaders `json:"headers,omitempty"`
	TimeoutThresholdMs int         `json:"timeoutThresholdMs,omitempty"` // in milliseconds
	Parameters         Parameters  `json:"parameters,omitempty"`
	Content            string      `json:"content,omitempty"`
}

// Validate returns error if not valid
func (v HTTPRequest) Validate() error {
	if v.URL == "" {
		return validation.NewErrRecordIsMissingRequiredField("url")
	}
	if v.TimeoutThresholdMs < 0 {
		return validation.NewErrBadRecordFieldValue("timeoutThresholdMs", "should be 0 or positive")
	}
	switch v.Method {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("method")
	case "GET", "HEAD", "OPTIONS":
		if v.Content != "" {
			return validation.NewErrBadRecordFieldValue("content", v.Method+"%v request can't have content")
		}
	case "POST", "PUT", "DELETE", "CONNECT", "TRACE", "PATCH":
		// According to https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
	default:
		return validation.NewErrBadRecordFieldValue("method", "unknown method: "+v.Method)
	}
	if err := v.Headers.Validate(); err != nil {
		return fmt.Errorf("invalid header: %w", err)
	}
	return nil
}

// HTTPWidgetDef holds info about a widget that makes HTTP requests
type HTTPWidgetDef struct {
	WidgetBase
	Request HTTPRequest `json:"request"`
}

// Validate returns error if not valid
func (v HTTPWidgetDef) Validate() error {
	if err := v.WidgetBase.Validate(false); err != nil {
		return err
	}
	if err := v.Request.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("request", err.Error())
	}
	return nil
}

// TabWidget describes a tab widget
type TabWidget struct {
	Title  string      `json:"title"`
	Widget BoardWidget `json:"widget,omitempty"`
}

// TabsWidgetDef describes set of tab widgets
type TabsWidgetDef struct {
	WidgetBase
	Tabs []TabWidget `json:"tabs"`
}

// Validate returns error if not valid
func (v TabsWidgetDef) Validate() error {
	if err := v.WidgetBase.Validate(false); err != nil {
		return err
	}
	return nil
}
