package Http

import (
	"fmt"
	"go_demo/Geecache"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	self     string		//用来记录自己的地址，包括了主机名和端口号
	basePath string		//自己添加的一段base api
}

// 实例化
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//实现日志的输出
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//实现server接口，完成服务端
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//判断前缀是否符合预先设定的
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// 检查url输入格式是否正确
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	//获取第一个访问地址
	groupName := parts[0]
	key := parts[1]

	group := Geecache.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}