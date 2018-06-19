package HtmlParser

import (
	"strings"
	"os"
	"bufio"
	"bytes"
	"fmt"
)

//标签属性
type HtmlTagAttr struct {
	AttrName string
	AttrVal  string
}

//html标签
type HtmlTag struct {
	TagName  string
	TagAttrs []HtmlTagAttr
	TagIndex int
	CodeStr  string
}

//栈内节点
type Node struct {
	val   HtmlTag
	pNode *Node
}
type Stack struct {
	pTop, pBottom *Node
	depth         int
}

//栈的初始化
func initStack(pStack *Stack) {
	pNew := new(Node)
	pNew.pNode = nil
	pStack.pTop = pNew
	pStack.pBottom = pNew
	pStack.depth = 0
	if pStack.pTop == nil || pStack.pBottom == nil {
		fmt.Println("分配头节点内存失败，程序退出")
		os.Exit(-1)
	}
}

//判断栈是否为空
func isEmpty(pStack *Stack) bool {
	if pStack.pTop == pStack.pBottom {
		pStack.depth = 0
		return true
	}
	return false
}

//入栈
func push(pStack *Stack, val HtmlTag) {
	pNew := new(Node)
	pNew.val = val
	if pNew == nil {
		fmt.Println("分配头节点内存失败，程序退出")
		os.Exit(-1)
	}
	pNew.pNode = pStack.pTop
	pStack.pTop = pNew
	pStack.depth++
}

//出栈
func pop(pStack *Stack) (bool, *HtmlTag) {
	if isEmpty(pStack) {
		return false, nil
	}
	reval := pStack.pTop.val
	pStack.pTop = pStack.pTop.pNode
	pStack.depth--
	return true, &reval
}

//字符串aa-bb—cc变aaBbCc
func strFirstToUpper(str string) string {
	strArr := strings.Split(str, "-")
	var upperStr string
	for index, tempStr := range strArr {
		tempStrLen := len([]rune(tempStr))
		for i := 0; i < tempStrLen; i++ {
			if i == 0 {
				if index == 0 {
					upperStr += string([]rune(tempStr)[i])
				} else {
					upperStr += string([]rune(tempStr)[i] - 32)
				}
			} else {
				upperStr += string([]rune(tempStr)[i])
			}
		}
	}
	return upperStr
}

//处理标签的class属性
func dealClass(template string) (newTemplate string, classAttr HtmlTagAttr) {
	template = strings.Trim(template, " ")
	var classSeparator = `class="`
	if strings.Index(template, classSeparator) > -1 {
		var templateRune = []rune(template)
		var classSeparatorIndex = strings.Index(template, classSeparator)
		/*for i, L := 0, len(templateRune); i < L; i++ {
			if string(templateRune[i:i+7]) == classSeparator {
				classSeparatorIndex = i
				break
			}
		}*/
		for i, L := classSeparatorIndex, len(templateRune); i < L; i++ {
			if templateRune[i] == '"' && templateRune[i-1] != '=' {
				classAttrStr := string(templateRune[classSeparatorIndex:i+1])
				classAttrArr := strings.Split(classAttrStr, "=")
				classAttr = HtmlTagAttr{classAttrArr[0], classAttrArr[1]}
				newTemplate = strings.Replace(template, classAttrStr, "", -1)
				break
			}
		}
	} else {
		newTemplate = template
		classAttr = HtmlTagAttr{"class", ""}
	}
	return newTemplate, classAttr
}

//处理标签的style属性
func dealStyle(template string) (newTemplate string, styleAttr HtmlTagAttr) {
	template = strings.Trim(template, " ")
	var styleSeparator = `style="`
	if strings.Index(template, styleSeparator) > -1 {
		var templateRune = []rune(template)
		var styleSeparatorIndex = strings.Index(template, styleSeparator)
		/*for i, L := 0, len(templateRune); i < L; i++ {
			if string(templateRune[i:i+7]) == styleSeparator {
				styleSeparatorIndex = i
				break
			}
		}*/
		for i, L := styleSeparatorIndex, len(templateRune); i < L; i++ {
			if templateRune[i] == '"' && templateRune[i-1] != '=' {
				styleAttrStr := string(templateRune[styleSeparatorIndex:i+1])
				styleAttrArr := strings.Split(styleAttrStr, "=")
				styleAttr = HtmlTagAttr{styleAttrArr[0], styleAttrArr[1]}
				newTemplate = strings.Replace(template, styleAttrStr, "", -1)
				break
			}
		}
	} else {
		newTemplate = template
		styleAttr = HtmlTagAttr{"style", ""}
	}
	return newTemplate, styleAttr
}

//处理中文和html标签同行的问题
func dealChinese(template string) (newTemplate string) {
	template = strings.Trim(template, " ")
	templateRune := []rune(template)
	if templateRune[0] != '<' {
		L := len(templateRune)
		for i := 0; i < L; i++ {
			if templateRune[i] == '<' {
				c := 0
				for j := i; j < L; j++ {
					if templateRune[j] == '>' {
						c++
						if c == 2 {
							newTemplate = string(templateRune[i:j+1])
							break
						}
					}
				}
				break
			}
		}
	} else {
		newTemplate = template
	}
	return newTemplate
}

//获取htmlTag的内容
func getHtmlTagContent(template string) (htmlTagContent string) {
	template = strings.Trim(template, " ")
	var leftSeparator, rightSeparator = '<', '>'
	templateRune := []rune(template)
	L := len(templateRune)
	for i := 0; i < L; i++ {
		if templateRune[i] == leftSeparator {
			for j := i; j < L; j++ {
				if templateRune[j] == rightSeparator {
					htmlTagContent = string(templateRune[i+1:j])
					break
				}
			}
			break
		}
	}
	return htmlTagContent
}

//生成整理html标签
func sortHtmlTags(path string) []HtmlTag {
	file, _ := os.Open(path)
	defer file.Close()
	var htmlTags []HtmlTag
	scanner := bufio.NewScanner(file)
	templateBuffer := bytes.Buffer{}
	for scanner.Scan() {
		if scanner.Text() == `<!DOCTYPE html>` {
			continue
		}
		templateBuffer.WriteString(" " + strings.TrimSpace(scanner.Text()))
		if !(strings.Contains(templateBuffer.String(), "<") && strings.Contains(templateBuffer.String(), ">")) {
			continue
		} else {
			var htmlTag = HtmlTag{TagName: "", TagAttrs: make([]HtmlTagAttr, 0)}
			var template = templateBuffer.String()
			templateBuffer.Reset()
			template = dealChinese(template)
			template, classAttr := dealClass(template)
			template, styleAttr := dealStyle(template)
			var htmlTagContent = getHtmlTagContent(template)
			var htmlTagContentArr = strings.Split(htmlTagContent, " ")
			for key, val := range htmlTagContentArr {
				if len(htmlTagContentArr) > 1 {
					if key == 0 {
						htmlTag.TagName = val
						if classAttr.AttrVal != "" && val != "html" && val != "head" && val != "meta" && val != "title" && val != "body" && !strings.Contains(val, "/") {
							htmlTag.TagAttrs = append(htmlTag.TagAttrs, classAttr)
						}
						if styleAttr.AttrVal != "" && val != "html" && val != "head" && val != "meta" && val != "title" && val != "body" && !strings.Contains(val, "/") {
							htmlTag.TagAttrs = append(htmlTag.TagAttrs, styleAttr)
						}
					} else {
						if val != "" {
							htmlTagAttrArr := strings.Split(val, "=")
							htmlTagAttr := HtmlTagAttr{htmlTagAttrArr[0], htmlTagAttrArr[1]}
							htmlTag.TagAttrs = append(htmlTag.TagAttrs, htmlTagAttr)
						}
					}
				} else {
					htmlTag.TagName = val
					if classAttr.AttrVal != "" && val != "html" && val != "head" && val != "meta" && val != "title" && val != "body" && !strings.Contains(val, "/") {
						htmlTag.TagAttrs = append(htmlTag.TagAttrs, classAttr)
					}
					if styleAttr.AttrVal != "" && val != "html" && val != "head" && val != "meta" && val != "title" && val != "body" && !strings.Contains(val, "/") {
						htmlTag.TagAttrs = append(htmlTag.TagAttrs, styleAttr)
					}
				}
			}
			htmlTags = append(htmlTags, htmlTag)
			if !strings.Contains(htmlTag.TagName, "/") && strings.Contains(template, "/") {
				var htmlEndTag = HtmlTag{TagName: "", TagAttrs: make([]HtmlTagAttr, 0)}
				htmlEndTag.TagName = "/" + htmlTag.TagName
				htmlTags = append(htmlTags, htmlEndTag)
			}
		}
	}
	var pStack Stack
	initStack(&pStack)
	//生成标签的层级
	for index, htmlTag := range htmlTags {
		push(&pStack, htmlTag)
		htmlTags[index].TagIndex = pStack.depth
		if strings.Contains(htmlTag.TagName, "/") {
			for {
				_, tag := pop(&pStack)
				if "/"+tag.TagName == htmlTag.TagName {
					/*if htmlTags[index-1].TagName == tag.TagName {
						htmlTags[index-1].TagIndex = pStack.depth + 1
					}*/
					break
				}
			}
		}
	}
	//处理不闭合标签
	for index, htmlTag := range htmlTags {
		if htmlTag.TagName == `head` {
			pIndex := index
			for {
				pIndex++
				htmlTags[pIndex].TagIndex = htmlTag.TagIndex + 1
				if htmlTags[pIndex].TagName == `/head` {
					break
				}
			}
		}
		if htmlTag.TagName == `colgroup` {
			pIndex := index
			for {
				pIndex++
				htmlTags[pIndex].TagIndex = htmlTag.TagIndex + 1
				if htmlTags[pIndex].TagName == `/colgroup` {
					break
				}
			}
		}
	}
	var sortHtmlTags []HtmlTag
	//去除闭合标签的结尾标签
	for _, htmlTag := range htmlTags {
		if strings.Contains(htmlTag.TagName, "/") {
			continue
		}
		htmlTag.CodeStr = `createElement("{!tagName!}",{{!attrs!},{!style!}},[{!child!}])`
		htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, "{!tagName!}", htmlTag.TagName, -1)
		if len(htmlTag.TagAttrs) > 0 {
			attrs := `attrs:{`
			attrsLen := len(htmlTag.TagAttrs)
			for index, attr := range htmlTag.TagAttrs {
				if index != attrsLen-1 {
					if attr.AttrName != "style" {
						if strings.Contains(attr.AttrName,"-"){
							attrs = attrs +`"`+ attr.AttrName + `":` + attr.AttrVal+","
						}else{
							attrs = attrs + attr.AttrName + `:` + attr.AttrVal+","
						}
					} else {
						styleAttrStr := "style:{"
						attr.AttrVal = strings.Replace(attr.AttrVal, `"`, ``, -1)
						styleAttrArr := strings.Split(attr.AttrVal, ";")
						for _, styleAttr := range styleAttrArr {
							if styleAttr != "" {
								style := strings.Split(styleAttr, ":")
								styleAttrStr = styleAttrStr + strFirstToUpper(style[0]) + ":" + `"` + style[1] + `"`
							}
						}
						styleAttrStr = styleAttrStr + "}"
						htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, "{!style!}", styleAttrStr, -1)
					}
				} else {
					if attr.AttrName != "style" {
						if strings.Contains(attr.AttrName,"-"){
							attrs = attrs +`"`+ attr.AttrName + `":` + attr.AttrVal
						}else{
							attrs = attrs + attr.AttrName + `:` + attr.AttrVal
						}
					} else {
						styleAttrStr := "style:{"
						attr.AttrVal = strings.Replace(attr.AttrVal, `"`, ``, -1)
						styleAttrArr := strings.Split(attr.AttrVal, ";")
						for _, styleAttr := range styleAttrArr {
							if styleAttr != "" {
								style := strings.Split(styleAttr, ":")
								styleAttrStr = styleAttrStr + strFirstToUpper(style[0]) + ":" + `"` + style[1] + `"`
							}
						}
						styleAttrStr = styleAttrStr + "}"
						htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, "{!style!}", styleAttrStr, -1)
					}
				}
			}
			attrs = attrs + `}`
			htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, "{!attrs!}", attrs, -1)
			if strings.Contains(htmlTag.CodeStr, ",{!style!}") {
				htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, "{!style!}", "style:{}", -1)
			}
		} else {
			htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, "{!attrs!}", "", -1)
			if strings.Contains(htmlTag.CodeStr, ",{!style!}") {
				htmlTag.CodeStr = strings.Replace(htmlTag.CodeStr, ",{!style!}", "style:{}", -1)
			}
		}
		sortHtmlTags = append(sortHtmlTags, htmlTag)
	}
	return sortHtmlTags
}

//将html解析成vue的render函数写法
func Parser(path string) string {
	sortHtmlTags := sortHtmlTags(path)
	maxLevel := 1
	for _, htmlTag := range sortHtmlTags {
		if maxLevel < htmlTag.TagIndex {
			maxLevel = htmlTag.TagIndex
		}
	}
	for i := maxLevel; i >= 1; i-- {
		level := i
		for index, htmlTag := range sortHtmlTags {
			if htmlTag.TagIndex == level {
				tempHtmlTag := htmlTag
				for i := index; i >= 0; i-- {
					if sortHtmlTags[i].TagIndex == level-1 {
						if tempHtmlTag.TagIndex == maxLevel || strings.Contains(tempHtmlTag.CodeStr, `{!child!}`) {
							tempHtmlTag.CodeStr = strings.Replace(tempHtmlTag.CodeStr, `{!child!}`, ``, -1)
						}
						if strings.Contains(sortHtmlTags[i].CodeStr, `{!child!}`) {
							sortHtmlTags[i].CodeStr = strings.Replace(sortHtmlTags[i].CodeStr, `{!child!}`, tempHtmlTag.CodeStr, -1)
						} else {
							/*lastIndex := strings.LastIndex(sortHtmlTags[i].CodeStr, `]`)*/
							lastIndex:=len([]rune(sortHtmlTags[i].CodeStr))-2
							sortHtmlTags[i].CodeStr = string([]rune(sortHtmlTags[i].CodeStr)[:lastIndex]) + `,` + tempHtmlTag.CodeStr + `])`
						}
						break
					}
				}
			}
		}
	}
	render := `render:function(createElement){ return ` + sortHtmlTags[0].CodeStr + `},`
	return render
}
