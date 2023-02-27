package web

import (
	"net/http"
	"fmt"
	"FabricEduction/web/controller"
)

func WebStart(app controller.Application)  {

	fs:= http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)
	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/help", app.Help)
	http.HandleFunc("/addEduInfo", app.AddEduShow)	
	http.HandleFunc("/addEdu", app.AddEdu)	
	http.HandleFunc("/queryPage", app.QueryPage)	
	http.HandleFunc("/query", app.FindCertByNoAndName)
	http.HandleFunc("/queryPage2", app.QueryPage2)	
	http.HandleFunc("/query2", app.FindByID)	
	http.HandleFunc("/modifyPage", app.ModifyShow)	
	http.HandleFunc("/modify", app.Modify)	
	http.HandleFunc("/upload", app.UploadFile)
	fmt.Println("启动Web服务, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}
}



