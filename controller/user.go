package controller

import (
	"encoding/json"
	"net/http"
	"updater-server/pkg/app"
	"updater-server/service"
)

type UserController struct {
	UserService *service.UserService
}

var userStr = `
{
	"name": "Serati Ma",
	"avatar": "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
	"userid": "00000001",
	"email": "antdesign@alipay.com",
	"signature": "海纳百川，有容乃大",
	"title": "交互专家",
	"group": "蚂蚁金服－某某某事业群－某某平台部－某某技术部－UED",
	"tags": [
	  {
		"key": "0",
		"label": "很有想法的"
	  },
	  {
		"key": "1",
		"label": "专注设计"
	  },
	  {
		"key": "2",
		"label": "辣~"
	  },
	  {
		"key": "3",
		"label": "大长腿"
	  },
	  {
		"key": "4",
		"label": "川妹子"
	  },
	  {
		"key": "5",
		"label": "海纳百川"
	  }
	],
	"notifyCount": 12,
	"unreadCount": 11,
	"country": "China",
	"geographic": {
	  "province": {
		"label": "浙江省",
		"key": "330000"
	  },
	  "city": {
		"label": "杭州市",
		"key": "330100"
	  }
	},
	"address": "西湖区工专路 77 号",
	"phone": "0752-268888888"
  }  
`

func (uc *UserController) UserInfo(c *app.Context) {

	mdata := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(userStr), &mdata)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(mdata)
}
