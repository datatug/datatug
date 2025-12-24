package dbviewer

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/dal-go/dalgo/recordset"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ tview.TableContent = (*TableContentRecordset)(nil)

type TableContentRecordset struct {
	tview.TableContentReadOnly
	recordset recordset.Recordset
}

func (t TableContentRecordset) GetCell(row, column int) *tview.TableCell {
	col := t.recordset.GetColumnByIndex(column)
	if row == 0 { // Header
		return tview.NewTableCell(col.Name()).SetTextColor(tcell.ColorLightBlue)
	}
	row--
	v, err := col.GetValue(row)
	if err != nil {
		return tview.NewTableCell(fmt.Sprintf("ERROR: %v", err)).SetTextColor(tcell.ColorRed)
	}
	if v == nil {
		return tview.NewTableCell("")
	}
	vType := reflect.TypeOf(v)
	if vType.Kind() == reflect.Slice {
		vVal := reflect.ValueOf(v)
		itemType := vType.Elem().String()
		if itemType == "uint8" {
			itemType = "byte"
		}
		length := strconv.Itoa(vVal.Len())
		return tview.NewTableCell(fmt.Sprintf("[]%v - %s", itemType, length)).SetTextColor(tcell.ColorGray)
	}
	switch tVal := v.(type) {
	case string:
		return tview.NewTableCell(tVal).SetTextColor(tcell.ColorLightGreen)
	case bool:
		var cell *tview.TableCell
		if tVal {
			cell = tview.NewTableCell("YES")
		} else {
			cell = tview.NewTableCell("NO")
		}
		return cell.SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorCornflowerBlue)
	case float64, float32:
		return tview.NewTableCell(fmt.Sprintf("%v", tVal)).SetTextColor(tcell.ColorLightGoldenrodYellow).SetAlign(tview.AlignRight)
	case int, int64, int32, int16, int8:
		return tview.NewTableCell(fmt.Sprintf("%d", tVal)).SetTextColor(tcell.ColorLightSalmon).SetAlign(tview.AlignRight)
	case time.Time:
		const timeColor = tcell.ColorLightCoral
		if tVal.Hour() == 0 && tVal.Minute() == 0 && tVal.Second() == 0 && tVal.Nanosecond() == 0 {
			if tVal.Location().String() == "UTC" {
				return tview.NewTableCell(tVal.Format("2006-01-02")).SetTextColor(timeColor)
			}
			return tview.NewTableCell(tVal.Format("2006-01-02") + " " + tVal.Location().String()).SetTextColor(timeColor)
		}
		return tview.NewTableCell(fmt.Sprintf("%v", v)).SetTextColor(timeColor)
	default:
		return tview.NewTableCell(fmt.Sprintf("%T:%v", v, v)).SetTextColor(tcell.ColorLightGray)
	}
}

func (t TableContentRecordset) GetRowCount() int {
	return t.recordset.RowsCount() + 1
}

func (t TableContentRecordset) GetColumnCount() int {
	return t.recordset.ColumnsCount()
}
