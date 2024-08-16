// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
	"golang.org/x/term"
)

type RowData struct {
	Id         string
	Hash       string
	State      string
	PrebuildId string
	CreatedAt  string
}

func ListBuilds(buildList []apiclient.Build) {
	re := lipgloss.NewRenderer(os.Stdout)

	headers := []string{"ID", "Configuration hash", "State", "Prebuild ID", "Created"}

	data := [][]string{}

	for _, pc := range buildList {
		var rowData *RowData
		var row []string

		rowData = getTableRowData(pc)
		row = getRowFromRowData(*rowData)
		data = append(data, row)
	}

	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(data)
		return
	}

	breakpointWidth := views.GetContainerBreakpointWidth(terminalWidth)

	minWidth := views_util.GetTableMinimumWidth(data)

	if breakpointWidth == 0 || minWidth > breakpointWidth {
		renderUnstyledList(buildList)
		return
	}

	t := table.New().
		Headers(headers...).
		Rows(data...).
		BorderStyle(re.NewStyle().Foreground(views.LightGray)).
		BorderRow(false).BorderColumn(false).BorderLeft(false).BorderRight(false).BorderTop(false).BorderBottom(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return views.TableHeaderStyle
			}
			return views.BaseCellStyle
		}).Width(breakpointWidth - 2*views.BaseTableStyleHorizontalPadding - 1)

	fmt.Println(views.BaseTableStyle.Render(t.String()))
}

func renderUnstyledList(buildList []apiclient.Build) {
	// todo
}

func getRowFromRowData(rowData RowData) []string {

	row := []string{
		views.NameStyle.Render(rowData.Id),
		views.DefaultRowDataStyle.Render(rowData.Hash),
		views.DefaultRowDataStyle.Render(rowData.State),
		views.DefaultRowDataStyle.Render(rowData.PrebuildId),
		views.DefaultRowDataStyle.Render(rowData.CreatedAt),
	}

	return row
}

func getTableRowData(build apiclient.Build) *RowData {
	rowData := RowData{"", "", "", "", ""}

	rowData.Id = build.Id + views_util.AdditionalPropertyPadding
	rowData.Hash = build.Hash
	rowData.State = string(build.State)
	rowData.PrebuildId = build.PrebuildId
	rowData.CreatedAt = util.FormatCreatedTime(build.CreatedAt)

	return &rowData
}