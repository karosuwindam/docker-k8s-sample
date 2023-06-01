package getkube

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"ingresslist/config"
	"ingresslist/getkube/client"
	"ingresslist/getkube/getdata"
	"log"
	"sort"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var setupflag bool = false

type OutputData struct {
	PName       string   `json:"Podname"`
	Namespace   string   `json:"Namspace"`
	Podinfo     string   `json:"Podinfo"`
	Port        []string `json:"Port"`
	Ip          string   `json:"Ip"`
	ClusterIP   string   `json:"ClusterIP"`
	HostNetwork bool     `json:"Hostnetwork"`
	SName       string   `json:"Servicename"`
	Selector    string   `json:"Selector"`
	URL         string   `json:"URL"`
}

type DomainData struct {
	Domain string       `json:"Domain"`
	Data   []OutputData `json:"Data"`
}

type OuteJsonData struct {
	Data []DomainData `json:"Data"`
}

// ドメイン別にデータを分ける
func GetDomainData(d []getdata.OutputData) (OuteJsonData, error) {
	//domainの一覧を取得
	var tmpdomainlist []string
	for _, v := range d {
		for _, vv := range v.URL {
			tmpdomainlist = append(tmpdomainlist, vv)
		}
	}
	var domainlist []string
	//からデータの追加
	domainlist = append(domainlist, "")
	for _, v := range tmpdomainlist {
		//ピリオド区切りのドメインを分割
		tmp := strings.Split(v, ".")
		if len(tmp) > 2 {
			ts := tmp[len(tmp)-2] + "." + tmp[len(tmp)-1]
			domainlist = append(domainlist, ts)
		} else {
			domainlist = append(domainlist, v)
		}
	}
	//重複を削除
	domainlist = removeDuplicate(domainlist)
	//ドメイン別にベースデータを作る
	var output OuteJsonData
	for _, v := range domainlist {
		var tmpdata DomainData
		tmpdata.Domain = v
		output.Data = append(output.Data, tmpdata)
	}
	//ドメイン別にデータを振り分ける
	for _, v := range d {
		for i, vv := range output.Data {
			if (len(v.URL) == 0) && (vv.Domain == "") {
				tmp := OutputData{
					PName:       v.PName,
					Namespace:   v.Namespace,
					Podinfo:     v.Podinfo,
					Port:        v.Port,
					Ip:          v.Ip,
					ClusterIP:   v.ClusterIP,
					HostNetwork: v.HostNetwork,
					SName:       v.SName,
					Selector:    v.Selector,
				}
				output.Data[i].Data = append(output.Data[i].Data, tmp)
				continue
			}
			for _, vvv := range v.URL {
				if (strings.Index(vvv, vv.Domain) > -1) && (vvv != "") && (vv.Domain != "") {
					tmp := OutputData{
						PName:       v.PName,
						Namespace:   v.Namespace,
						Podinfo:     v.Podinfo,
						Port:        v.Port,
						Ip:          v.Ip,
						ClusterIP:   v.ClusterIP,
						HostNetwork: v.HostNetwork,
						SName:       v.SName,
						Selector:    v.Selector,
						URL:         vvv,
					}
					output.Data[i].Data = append(output.Data[i].Data, tmp)
				}
			}
		}
	}
	//ソート
	sort.Slice(output.Data, func(i, j int) bool {
		return output.Data[i].Domain > output.Data[j].Domain
	})
	for i := 0; i < len(output.Data); i++ {
		tmp := output.Data[i].Data
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].PName < tmp[j].PName
		})
		output.Data[i].Data = tmp
	}
	for i, v := range output.Data {
		if v.Domain == "" || len(v.Domain) < 2 {
			output.Data[i].Domain = v.Domain
		} else if v.Domain[len(v.Domain)-1:len(v.Domain)] == "/" {
			output.Data[i].Domain = "*." + v.Domain[:len(v.Domain)-1]
		} else {
			output.Data[i].Domain = "*." + v.Domain
		}
	}

	return output, nil
}

// 重複を削除
func removeDuplicate(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for _, v := range args {
		if !encountered[v] {
			encountered[v] = true
			results = append(results, v)
		}
	}
	return results
}

func GetJsonData() string {
	if setupflag == false {
		return ""
	}
	var output string
	for i := 0; i < 3; i++ {
		if tmp, err := getdata.CreateData(); err != nil {
			log.Println(err)
			output = ""
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			//ドメイン別に分ける
			if tmp, err := GetDomainData(tmp); err != nil {
				log.Println(err)
				output = ""
			} else {
				//配列からjsonに変換
				tmpjson, _ := json.Marshal(tmp)
				var buf bytes.Buffer
				json.Indent(&buf, tmpjson, "", "  ")
				output = string(buf.Bytes())
			}
		}
		break
	}

	return output
}

func Run(ctx context.Context, ch chan<- error) {
	if setupflag == false {
		ch <- errors.New("Setup is not set")
		return
	}
	var err error = nil
	go func() {
		for {
			select {
			case <-ctx.Done():
				ch <- err
				return
			default:
				err = getdata.GetKubeData()
				time.Sleep(10 * time.Second)
			}
		}
	}()
	<-ctx.Done()
}

func Setup(cfg *config.Config) client.ClientGo {
	output := client.ClientGo{}
	if cfg.Kube.KubePathFlag {
		tmp, err := newClient_cli(cfg.Kube.KubePath)
		if err == nil {
			output.Client = tmp
			output.Flag = true
		} else {
			log.Println(err.Error())
		}

	} else {
		tmp, err := newClient_con()
		if err == nil {
			output.Client = tmp
			output.Flag = true
		} else {
			log.Println(err.Error())
		}
	}
	if output.Flag {
		setupflag = true
		getdata.Setup(&output)
	}
	return output
}

func newClient_con() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func newClient_cli(path string) (kubernetes.Interface, error) {
	// kubeConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", path, "(optional) absolute path to the kubeconfig file")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}
