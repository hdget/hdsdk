package excel

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"path"
	"reflect"
)

type ec struct {
	sheetName string
	header    reflect.Type
	f         *excelize.File
	style     int
	rows      []interface{}
}

type Excel interface {
	SaveFile(dir, filename string) error
	GetBytes() ([]byte, error)
}

const (
	numFmtTextPlaceHolder = 49 // '@'文本占位符。单个@的作用是引用单元格内输入的原始内容，将其以文本格式显示出来
	defaultSheetName      = "Sheet1"
)

func New(rows []interface{}, args ...string) (Excel, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty rows")
	}

	// 创建表格
	f := excelize.NewFile()

	sheetName := defaultSheetName
	if len(args) > 0 {
		sheetName = args[0]
	}

	sheet, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(sheet)

	// 通过反射获取表头type定义
	headerType, err := getExcelHeaderType(rows)
	if err != nil {
		return nil, errors.Wrap(err, "get excel header type")
	}

	// 创建一个默认文本类型的style
	style, _ := f.NewStyle(&excelize.Style{
		NumFmt: numFmtTextPlaceHolder,
	})

	return &ec{
		sheetName: sheetName,
		header:    headerType,
		f:         f,
		style:     style,
		rows:      rows,
	}, nil
}

func (ec *ec) Close() error {
	ec.rows = nil
	return ec.f.Close()
}

// SaveFile 保存文件
func (ec *ec) SaveFile(dir, filename string) error {
	err := ec.process(ec.rows)
	if err != nil {
		return errors.Wrap(err, "process rows")
	}

	err = ec.f.SaveAs(path.Join(dir, filename))
	if err != nil {
		return errors.Wrap(err, "save excel file")
	}

	return nil
}

func (ec *ec) GetBytes() ([]byte, error) {
	err := ec.process(ec.rows)
	if err != nil {
		return nil, errors.Wrap(err, "process rows")
	}

	// 获取文件写入buffer
	buf, err := ec.f.WriteToBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "excel write to buffer")
	}

	return buf.Bytes(), nil
}

// process 处理数据
func (ec *ec) process(rows []interface{}) error {
	// 输出数据
	for i := 0; i < ec.header.NumField(); i++ {
		// 输出表头
		colName := ec.header.Field(i).Tag.Get("col_name")
		colAxis := ec.header.Field(i).Tag.Get("col_axis")

		if colName != "" {
			axis := fmt.Sprintf("%s%d", colAxis, 1)

			err := ec.f.SetCellValue(ec.sheetName, axis, colName)
			if err != nil {
				return errors.Wrap(err, "generate header")
			}
		}

		// 设置所有有效列的格式为文本类型
		if colAxis != "" {
			err := ec.f.SetColStyle(ec.sheetName, colAxis, ec.style)
			if err != nil {
				return errors.Wrap(err, "set col style")
			}
		}

		// 如果有列坐标的输出行数据
		if colAxis != "" {
			for j, r := range rows {
				axis := fmt.Sprintf("%s%d", colAxis, j+2)
				value := reflect.ValueOf(r).Elem().FieldByName(ec.header.Field(i).Name)
				err := ec.f.SetCellValue(ec.sheetName, axis, value)

				if err != nil {
					return errors.Wrap(err, "generate row")
				}
			}
		}
	}
	return nil
}

// getExcelHeader 通过反射获取表格的标题属性, 生成表格表头
func getExcelHeaderType(rows []interface{}) (reflect.Type, error) {
	// 取第一行并获取Elem()
	t := reflect.TypeOf(rows[0])
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid struct: %s", t.String())
	}
	return t, nil
}
