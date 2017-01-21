package product

import (
	"encoding/json"
	"fmt"
	"goERP/controllers/base"
	md "goERP/models"

	"strconv"
	"strings"
)

type ProductTemplateController struct {
	base.BaseController
}

func (ctl *ProductTemplateController) Post() {
	ctl.URL = "/product/template/"
	ctl.Data["URL"] = ctl.URL
	action := ctl.Input().Get("action")
	fmt.Println(action)
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "attribute":
		ctl.ProductTemplateAttributes()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}
func (ctl *ProductTemplateController) Get() {
	ctl.URL = "/product/template/"
	ctl.PageName = "产品款式管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()
	}
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuProductTemplateActive"] = "active"
}

// Put 修改产品款式
func (ctl *ProductTemplateController) Put() {
	fmt.Println("enter put")
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	template := new(md.ProductTemplate)
	var (
		err    error
		id     int64
		errs   []error
		debugs []string
	)
	if err = json.Unmarshal([]byte(postData), template); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(template)).Type().Name()
		if id, errs = md.AddProductTemplate(template, &ctl.User); len(errs) == 0 {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "数据创建失败"

			for _, item := range errs {
				debugs = append(debugs, item.Error())
			}
			result["debug"] = debugs
		}
	}
	if err != nil {
		result["code"] = "failed"
		debugs = append(debugs, err.Error())
		result["debug"] = debugs
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}
func (ctl *ProductTemplateController) ProductTemplateAttributes() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")

	result := make(map[string]interface{})
	if paginator, arrs, err := md.GetAllProductAttributeLine(query, fields, sortby, order, offset, limit); err == nil {
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			attributes := make(map[string]string)
			if line.Attribute != nil {
				attributes["id"] = strconv.FormatInt(line.Attribute.ID, 10)
				attributes["name"] = line.Attribute.Name
			}
			tmpValues := make(map[string]string)
			if line.ProductTemplate != nil {
				tmpValues["id"] = strconv.FormatInt(line.ProductTemplate.ID, 10)
				tmpValues["name"] = line.ProductTemplate.Name
			}
			attributeValuesLines := make([]interface{}, 0, 4)
			attributeValues := line.AttributeValues
			if attributeValues != nil {
				for _, line := range attributeValues {
					mapAttributeValues := make(map[string]string)
					mapAttributeValues["id"] = strconv.FormatInt(line.ID, 10)
					mapAttributeValues["name"] = line.Name
					attributeValuesLines = append(attributeValuesLines, oneLine)
				}

			}
			oneLine["Attribute"] = attributes
			oneLine["ProductTemplate"] = tmpValues
			oneLine["AttributeValues"] = attributeValuesLines

			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()

}
func (ctl *ProductTemplateController) PostCreate() {
	result := make(map[string]interface{})
	postData := ctl.GetString("postData")
	template := new(md.ProductTemplate)
	var (
		err  error
		id   int64
		errs []error
	)
	if err = json.Unmarshal([]byte(postData), template); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(template)).Type().Name()
		if id, errs = md.AddProductTemplate(template, &ctl.User); len(errs) == 0 {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "数据创建失败"
			var debugs []string
			for _, item := range errs {
				debugs = append(debugs, item.Error())
			}
			result["debug"] = debugs
		}
	} else {
		result["code"] = "failed"
		result["message"] = "请求数据解析失败"
		result["debug"] = err.Error()
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}
func (ctl *ProductTemplateController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if template, err := md.GetProductTemplateByID(idInt64); err == nil {
				ctl.PageAction = template.Name
				ctl.Data["Tp"] = template
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["FormField"] = "form-edit"
	ctl.Layout = "base/base.html"
	ctl.TplName = "product/product_template_form.html"
}
func (ctl *ProductTemplateController) Detail() {
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["FormTreeField"] = "form-tree-edit"
	ctl.Data["Action"] = "detail"
}
func (ctl *ProductTemplateController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.Data["FormField"] = "form-create"
	ctl.Data["FormTreeField"] = "form-tree-create"
	ctl.TplName = "product/product_template_form.html"
}

func (ctl *ProductTemplateController) Validator() {
	name := strings.TrimSpace(ctl.GetString("Name"))
	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	obj, err := md.GetProductTemplateByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.ID {

				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// 获得符合要求的款式数据
func (ctl *ProductTemplateController) productTemplateList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.ProductTemplate
	paginator, arrs, err := md.GetAllProductTemplate(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Name
			oneLine["sequence"] = line.Sequence
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			oneLine["defaultCode"] = line.DefaultCode
			category := line.Category
			if category != nil {
				oneLine["category"] = category.Name
			}
			oneLine["variantCount"] = line.VariantCount
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}
func (ctl *ProductTemplateController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	if result, err := ctl.productTemplateList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *ProductTemplateController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-product-template"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "product/product_template_list_search.html"
}
