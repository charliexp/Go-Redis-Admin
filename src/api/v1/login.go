package v1

import (
	"Go-Redis-Admin/app/config"
	"Go-Redis-Admin/src/common/cookie"
	"Go-Redis-Admin/src/common/crypto"
	"Go-Redis-Admin/src/common/mysql"
	"Go-Redis-Admin/src/common/response"
	"log"
	"net/http"
	"strconv"
)

func (h *Handlers) LoginAction(w http.ResponseWriter, r *http.Request) {
	// init return
	var data = map[string]string{}
	jr := &response.JsonResponse{
		0,
		"faild",
		data,
	}
	if r.Method != "POST" {
		jr.Code = 2
		jr.Msg = "请求出错啦"
		response.OuputJson(w, jr)
		return
	}
	log.Println("API Login Action")
	r.ParseForm() //解析参数，默认是不会解析的
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }
	// fmt.Println("username", r.FormValue("username"))
	// fmt.Println("password", r.FormValue("password"))
	var username = r.FormValue("username")
	var password = r.FormValue("password")
	if username == "" {
		jr.Code = 3
		jr.Msg = "用户名不能为空"
		response.OuputJson(w, jr)
		return
	}
	if password == "" {
		jr.Code = 4
		jr.Msg = "密码不能为空"
		response.OuputJson(w, jr)
		return
	}
	// 校验密码
	userId := checkPass(username, password)
	if userId == -1 {
		jr.Code = 5
		jr.Msg = "用户名不存在"
		response.OuputJson(w, jr)
		return
	} else if userId == -2 {
		jr.Code = 6
		jr.Msg = "密码不正确"
		response.OuputJson(w, jr)
		return
	}
	// 更新登录IP
	var ip = r.RemoteAddr
	_, err := mysql.SetLastLoginIp(userId, ip)
	// 设置用户的cookie
	plainText := strconv.Itoa(int(userId)) + "|" + username
	keyText := config.UserAesKey
	cipherText, err := crypto.AesEncode(plainText, keyText)
	if err != nil {
		jr.Code = 7
		jr.Msg = "加密失败"
		response.OuputJson(w, jr)
		return
	}
	cookie.Set(w, "grd_username", username, "/", 8600)
	cookie.Set(w, "grd_auth", cipherText, "/", 8600)

	log.Println("cipherText", cipherText)
	plainTextCopy, err := crypto.AesDecode(cipherText, keyText)
	log.Println("plainTextCopy", plainTextCopy)

	jr.Code = 1
	jr.Msg = "登录成功"
	response.OuputJson(w, jr)
	return

}

func checkPass(username, password string) int32 {
	var user *mysql.UserEntity
	user, err := mysql.GetUserPass(username)
	if err != nil {
		log.Println(err)
		return -1
	}
	if user.Password != crypto.Md5Double(password, user.Salt) {
		log.Println("password error")
		return -2
	}
	return user.Id
}
