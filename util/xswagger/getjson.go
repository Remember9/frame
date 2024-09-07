package xswagger

import (
	"esfgit.leju.com/golang/frame/config"
	"log"
	"reflect"
	"strings"
)

//根据逻辑修改
type fieldi struct {
	Name, Type, Format, SubType, SubFormat string
}
type message struct {
	Name       string
	Properties []*fieldi
}

type service struct {
	Name, Host, Pkg, Scheme string
	Rpcs                    []*Rpc
	Messages                []*message
}

type Rpc struct {
	Name, Stype, Route, Rep, Req string
	MReq                         *message
}

var Services = service{}
var properties []interface{}

var tmpf = map[string]map[string]*fieldi{}
var fields = map[string][]*fieldi{}

func AddPropertie(p ...interface{}) {
	properties = append(properties, p...)
}
func init() {
	Services.Name = "mysrv"
	if appConfig := config.GetAppConfig(); appConfig.Name != "" {
		Services.Name = appConfig.Name
		if appConfig.Version != "" {
			Services.Name = Services.Name + ":" + appConfig.Version
		}
	}
}
func (s *service) AddHost(host, scheme string) {
	s.Host = host
	s.Host = strings.ReplaceAll(s.Host, "http://", "")
	s.Host = strings.ReplaceAll(s.Host, "https://", "")
	s.Host = strings.ReplaceAll(s.Host, "127.0.0.1", "localhost")
	s.Scheme = scheme
}

func (s *service) AddRpc(name, stype, route string, ireq, irep interface{}) {
	AddPropertie(ireq, irep)
	req := strings.ReplaceAll(reflect.ValueOf(ireq).Type().String(), "*", "")
	rep := strings.ReplaceAll(reflect.ValueOf(irep).Type().String(), "*", "")
	s.Rpcs = append(s.Rpcs, &Rpc{
		Name:  name,
		Stype: stype,
		Route: route,
		Req:   req,
		Rep:   rep,
	})
}

var isDone bool

func (s *service) GetSwagger() ([]byte, error) {
	if !isDone {
		for _, p := range properties {
			reflictValue2(p)
		}
		var fmap = map[string][]*fieldi{}
		for k, v := range tmpf {
			for _, sv := range v {
				fmap[k] = append(fmap[k], sv)
			}
		}
		for k, v := range fmap {
			s.Messages = append(s.Messages, &message{
				Name:       k,
				Properties: v,
			})
		}
		for k, v := range s.Rpcs {
			if m, ok := fmap[v.Req]; ok {
				s.Rpcs[k].MReq = &message{
					Name:       v.Req,
					Properties: m,
				}
			}
		}
		isDone = true
	}

	b, err := s.execute()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return b, nil
}

func reflictValue2(p interface{}) {
	object := reflect.ValueOf(p)
	myref := object.Elem()
	typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		args := typeOfType.Field(i).Name
		if args[0] >= 97 && args[0] <= 122 {
			continue
		}
		field := myref.Field(i)
		fname := strings.ToLower(args)
		if fn := typeOfType.Field(i).Tag.Get("json"); fn != "" {
			fname = strings.ReplaceAll(fn, ",omitempty", "")
		}
		srvName := myref.Type().String()
		tmpType := typeOfType.Field(i).Type.String()
		tmpFormat := getFormat(tmpType)
		tmpfieldi := &fieldi{Name: fname, Type: getType(tmpType), Format: tmpFormat}
		if field.Kind() == reflect.Slice {
			tmpfieldi.Type = "array"
			tmpfieldi.SubType = strings.ReplaceAll(typeOfType.Field(i).Type.String(), "[]", "")
			tmpfieldi.SubFormat = getFormat(tmpfieldi.SubType)
			tmpfieldi.SubType = getType(tmpfieldi.SubType)
		}
		if tmpf[srvName] == nil {
			tmpf[srvName] = map[string]*fieldi{fname: tmpfieldi}
		} else {
			tmpf[srvName][fname] = tmpfieldi
		}

		//fields[srvName] = append(fields[srvName],tmpfieldi)
		if field.Kind() == reflect.Slice && typeOfType.Field(i).Type.Elem().Kind() == reflect.Ptr {
			reflictValue2(reflect.New(field.Type().Elem().Elem()).Interface())
		}
		if field.Kind() == reflect.Ptr {
			reflictValue2(reflect.New(field.Type().Elem()).Interface())
		}
	}

}
