package Geecache

import (
	"fmt"
	"go_demo/Multinodes"
	"go_demo/Singleflight"
	"log"
	"sync"
)

//主要内容接口型函数的定义以及使用方法

//定义一个接口，里面包含了一个回调函数
type Getter interface {
	Get(key string) ([]byte, error)
}

//定义函数类型
type GetterFunc func(key string) ([]byte, error)

//使用对应的函数类型实现接口
func (f GetterFunc) Get(key string) ([]byte, error){
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     Multinodes.PeerPicker
	loader *Singleflight.Group
}

var(
	mutex2 sync.Mutex
	groups = make(map[string]*Group)
)

//实例化group，每一个group都拥有唯一的名字，分别代表不同的缓存
func NewGroup(name string, cacheBytes int64, getter Getter) *Group{
	if getter == nil{
		panic("nil getter")
	}
	mutex2.Lock()
	defer mutex2.Unlock()
	g := &Group{
		name:	name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &Singleflight.Group{},
	}
	groups[name] = g
	return g
}

//获取一个group
func GetGroup(name string) *Group{
	mutex2.Lock()
	defer mutex2.Unlock()
	g := groups[name]
	return g
}

//group的get方法
func (g *Group) Get(key string) (ByteView, error){
	if key == ""{
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok :=g.mainCache.get(key); ok{
		log.Println("[Geecache] hit")
		return v, nil
	}

	return g.load(key)
}


//从本地资源获取
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)		//回调函数获取本地资源
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

//把从本地获取的元数据加入到缓存maincache中
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}



func (g *Group) RegisterPeers(peers Multinodes.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer Multinodes.PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}