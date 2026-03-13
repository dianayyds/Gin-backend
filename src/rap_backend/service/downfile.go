package service

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cihub/seelog"
	excelize "github.com/xuri/excelize/v2"
)

//标注内容导出Excel文件
func exportTaskFile(list []TaskCallLabelListDown, headers []string, filename string) error {
	f := excelize.NewFile()
	// style, err := f.NewStyle(`{"border":[{"type":"left","color":"0000FF","style":3},{"type":"top","color":"00FF00","style":4},{"type":"bottom","color":"FFFF00","style":5},{"type":"right","color":"FF0000","style":6}]}`)
	// if err != nil {
	// 	println(err.Error())
	// }
	style1, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#D0DAE4"],"pattern":1},"border":[{"type":"bottom","color":"BBBBBB","style":2},{"type":"right","color":"#BBBBBB","style":2}]}`)
	style2, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#FFFF00"],"pattern":1},"border":[{"type":"bottom","color":"BBBBBB","style":2},{"type":"right","color":"#BBBBBB","style":2}]}`)
	style3, err := f.NewStyle(`{"fill":{"type":"pattern","color":["#DAA569"],"pattern":1},"border":[{"type":"bottom","color":"BBBBBB","style":1},{"type":"right","color":"#BBBBBB","style":1}]}`)
	if err != nil {
		fmt.Println(err)
	}
	excelTitleArr := map[string]int{
		"RobotType":              style1,
		"RobotName":              style1,
		"StartTime":              style1,
		"CallID":                 style1,
		"BillSec":                style1,
		"Intention":              style1,
		"TalkRound":              style1,
		"Sentence1":              style1,
		"Final Intention":        style2,
		"Check":                  style2,
		"Problem":                style2,
		"Remark":                 style2,
		"Chinese":                style3,
		"Dialogue Fluency":       style3,
		"Real Intention":         style3,
		"Speech Recognition":     style3,
		"Semantic Comprehension": style3,
		"Robot language":         style3,
		"Robot Reaction Speed":   style3,
		"Noise situation":        style3,
		"Robot Recognition":      style3,
		"User cooperation":       style3,
		"Problem2":               style3,
		"Remark2":                style3,
	}
	sheet1 := "Sheet1"
	cellNames := []string{}
	for k, v := range headers {
		cellInt := 65 + k
		cellName := string(cellInt)
		if cellInt > 90 {
			cellName = "A" + string(cellInt-26)
		}
		cellNames = append(cellNames, cellName)

		style := style3
		if st, ok := excelTitleArr[v]; ok {
			style = st
		}
		v = strings.Replace(v, "Sentence1", "Sentence", -1)
		f.SetCellValue(sheet1, cellName+"1", v)
		f.SetCellStyle(sheet1, cellName+"1", cellName+"1", style)

	}
	// f.SetCellRichText("Sheet1", "A1", richTextRun)
	for k, val := range list {
		kstr := fmt.Sprintf("%d", k+2)
		// f.SetCellValue(sheet1, cellNames[0]+kstr, val.CallID)
		for lk, lab := range val.Labels {
			value := lab.LabelValue
			if lab.AuditorContent != "" {
				value = lab.AuditorContent
			}
			if lab.IsColor == 1 {
				//富文本
				richText := labelValueToRichText(value)
				f.SetCellRichText(sheet1, cellNames[lk]+kstr, richText)
				continue
			}
			f.SetCellValue(sheet1, cellNames[lk]+kstr, value)
		}
	}

	err = f.SaveAs(filename)
	if err != nil {
		seelog.Errorf("exportTaskFile SaveAs err:%s, filename:%s", err.Error(), filename)
	}
	err = os.Chmod(filename, 0644)
	if err != nil {
		seelog.Errorf("exportTaskFile Chmod err:%s, filename:%s", err.Error(), filename)
	}
	return err
}

func labelValueToRichText(text string) []excelize.RichTextRun {
	text = html.UnescapeString(text)
	richTexts := make([]excelize.RichTextRun, 0)
	regSpan := regexp.MustCompile(`<span[\S\s]+?>[\S\s]+?</span>`)
	spanTextArr := regSpan.Split(text, -1)
	spanIndexArr := regSpan.FindAllStringSubmatch(text, -1)

	if len(spanTextArr) < 2 {
		//无span标签
		rich := excelize.RichTextRun{
			Text: trimPLabel(text),
		}
		richTexts = append(richTexts, rich)
		return richTexts
	}
	siaLen := len(spanIndexArr)
	for i, sta := range spanTextArr {
		rich := excelize.RichTextRun{
			Text: trimPLabel(sta),
		}
		richTexts = append(richTexts, rich)
		if i < siaLen {
			//解析span里的内容
			spanRich := getSpanRichText(spanIndexArr[i][0])
			richTexts = append(richTexts, spanRich)
		}
	}

	return richTexts
}

func getSpanRichText(text string) excelize.RichTextRun {
	rText := ""
	rColor := "000000"
	regSpanText := regexp.MustCompile(`>[\S\s]+?<`)
	regSpanColor := regexp.MustCompile(`\d+`)
	//获取颜色
	colorArr := regSpanColor.FindAllStringSubmatch(text, -1)
	if len(colorArr) < 3 {
		seelog.Errorf("getSpanRichText err %s, arr:%v", text, colorArr)
	} else {
		rColor = rgbArr2Hex(colorArr[:3])
	}
	//获取span内容
	spanArr := regSpanText.FindAllStringSubmatch(text, -1)
	for _, sa := range spanArr {
		rText += strings.Replace(strings.Replace(sa[0], ">", "", -1), "<", "", -1)
	}

	richText := excelize.RichTextRun{
		Text: rText,
		Font: &excelize.Font{
			Color: rColor,
		},
	}
	return richText
}

func trimPLabel(text string) string {
	te := strings.ReplaceAll(text, "</p><p>", "\n")
	te = strings.ReplaceAll(te, "</p>", "\n")
	te = strings.ReplaceAll(te, "<p>", "\n")
	te = strings.ReplaceAll(te, "<br>", "\n")
	te = strings.ReplaceAll(te, "<br/>", "\n")
	te = strings.ReplaceAll(te, "</br>", "\n")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	te = re.ReplaceAllString(te, "\n")
	return te
}

func rgbArr2Hex(rgbArr [][]string) string {
	cl := ""
	for _, v := range rgbArr {
		val, _ := strconv.ParseInt(strings.TrimSpace(v[0]), 10, 64)
		if val > 255 {
			seelog.Errorf("getSpanRichText rgbArr2Hex, arr:%v", rgbArr)
		}
		cl += int216(val)
	}
	return cl
}

func labelValueToRichText1(text string) []excelize.RichTextRun {
	richTexts := make([]excelize.RichTextRun, 0)
	reg := regexp.MustCompile(`<span style=\"color: rgb\((.*?)\);">(?s:.*?)</span>`)
	b := reg.FindAllStringSubmatchIndex(text, -1)
	if len(b) > 0 {
		ib := reg.FindAllStringSubmatch(text, -1)
		l := len(text)
		s := 0
		for k, v := range b {
			if k == 0 {
				te := text[s:v[0]]
				te = strings.ReplaceAll(te, "</p><p>", "\n")
				te = strings.ReplaceAll(te, "</p>", "\n")
				te = strings.ReplaceAll(te, "<p>", "\n")
				rich := excelize.RichTextRun{
					Text: te,
				}
				richTexts = append(richTexts, rich)
				s = v[0]
				continue
			}
			te := text[s:v[0]]
			rtx1 := ib[k-1][0]
			rtx1 = strings.ReplaceAll(rtx1, "<span style=\"color: rgb("+ib[k-1][1]+");\">", "")
			rtx1 = strings.ReplaceAll(rtx1, "</span>", "")
			rich1 := excelize.RichTextRun{
				Text: rtx1,
				Font: &excelize.Font{
					Color: rgb2Hex(ib[k-1][1]),
				},
			}
			richTexts = append(richTexts, rich1)
			te = strings.ReplaceAll(te, ib[k-1][0], "")
			te = strings.ReplaceAll(te, "</p><p>", "\n")
			te = strings.ReplaceAll(te, "</p>", "\n")
			te = strings.ReplaceAll(te, "<p>", "\n")
			rich := excelize.RichTextRun{
				Text: te,
			}
			richTexts = append(richTexts, rich)
			s = v[0]
		}
		if s < l {
			te := text[s:]
			rtx1 := ib[len(b)-1][0]
			rtx1 = strings.ReplaceAll(rtx1, "<span style=\"color: rgb("+ib[len(b)-1][1]+");\">", "")
			rtx1 = strings.ReplaceAll(rtx1, "</span>", "")
			rich1 := excelize.RichTextRun{
				Text: rtx1,
				Font: &excelize.Font{
					Color: rgb2Hex(ib[len(b)-1][1]),
				},
			}
			richTexts = append(richTexts, rich1)
			te = strings.ReplaceAll(te, ib[len(b)-1][0], "")
			te = strings.ReplaceAll(te, "</p><p>", "\n")
			te = strings.ReplaceAll(te, "</p>", "\n")
			te = strings.ReplaceAll(te, "<p>", "\n")
			rich := excelize.RichTextRun{
				Text: te,
			}
			richTexts = append(richTexts, rich)
		}
	} else {
		te := strings.ReplaceAll(text, "</p><p>", "\n")
		te = strings.ReplaceAll(te, "</p>", "\n")
		te = strings.ReplaceAll(te, "<p>", "\n")
		rich := excelize.RichTextRun{
			Text: te,
		}
		richTexts = append(richTexts, rich)
	}

	return richTexts
}

func rgb2Hex(rgb string) string {
	ll := strings.Split(rgb, ",")
	cl := ""
	for _, v := range ll {
		val, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		cl += int216(val)
	}
	return cl
}

func int216(red int64) string {
	r := strconv.FormatInt(red, 16)
	if len(r) == 1 {
		r = "0" + r
	}
	return r
}
