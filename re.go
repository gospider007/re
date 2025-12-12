package re

import (
	"log"
	"regexp"
	"sync"
)

type ReData struct {
	data []string
}
type cancheData struct {
	regexs    map[string]*regexp.Regexp
	regexLock sync.RWMutex
}

func (obj *cancheData) get(reg string) (*regexp.Regexp, bool) {
	obj.regexLock.RLock()
	defer obj.regexLock.RUnlock()
	r, ok := obj.regexs[reg]
	return r, ok
}
func (obj *cancheData) set(reg string, r *regexp.Regexp) {
	obj.regexLock.Lock()
	defer obj.regexLock.Unlock()
	obj.regexs[reg] = r
}

var cacheMap = &cancheData{
	regexs: make(map[string]*regexp.Regexp),
}
var disCache bool

func Cache(val bool) {
	disCache = !val
}

// 返回分组的匹配
func (obj *ReData) Group(nums ...int) string {
	var num int
	if len(nums) > 0 {
		num = nums[0]
	}
	return obj.data[num]
}

func compile(reg any) (*regexp.Regexp, error) {
	switch val := reg.(type) {
	case string:
		if !disCache {
			r, ok := cacheMap.get(val)
			if ok {
				return r, nil
			}
			if r, err := regexp.Compile(val); err != nil {
				log.Printf("compile regexp error: %s", err.Error())
				return nil, err
			} else {
				cacheMap.set(val, r)
				return r, nil
			}
		}
		r, err := regexp.Compile(val)
		if err != nil {
			log.Printf("compile regexp error: %s", err.Error())
			return nil, err
		}
		return r, err
	case *regexp.Regexp:
		return val, nil
	default:
		panic("reg must be string or *regexp.Regexp")
	}
}

// 搜索
func Search(reg any, txt string) *ReData {
	regex, err := compile(reg)
	if err != nil {
		return nil
	}
	data := regex.FindStringSubmatch(txt)
	if len(data) == 0 {
		return nil
	}
	return &ReData{data: data}
}

// find 所有
func FindAll(reg any, txt string) []*ReData {
	datas := []*ReData{}
	regex, err := compile(reg)
	if err != nil {
		return nil
	}
	results := regex.FindAllStringSubmatch(txt, -1)
	for _, result := range results {
		datas = append(datas, &ReData{data: result})
	}
	return datas
}

// 替换匹配
func Sub(reg any, rep string, txt string) string {
	regex, err := compile(reg)
	if err != nil {
		return txt
	}
	return regex.ReplaceAllString(txt, rep)
}

// 使用方法替换匹配
func SubFunc(reg any, rep func(string) string, txt string) string {
	regex, err := compile(reg)
	if err != nil {
		return txt
	}
	return regex.ReplaceAllStringFunc(txt, rep)
}

// 分割
func Split(reg any, txt string) []string {
	regex, err := compile(reg)
	if err != nil {
		return []string{txt}
	}
	return regex.Split(txt, -1)
}

// 转义
func Quote(reg string) string {
	return regexp.QuoteMeta(reg)
}
