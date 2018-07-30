package lib

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type ServerConfigType struct {
	ServerNum int
	NodeIndex int
}

var AppId1 = "278ae090e15146d0932c45238b3d941a"

var AppId2 = "4a98a59f208448f5bff6dad92fdf8bdb"

var ServerConfig = ServerConfigType{}

var AccessToken = "2b8d85665941fba2e8a26f7ff7093f57683cebb245273915271d0dbed240e96658e5722c57bd7786ebea7d0797e9be6deacec755d8188a637344236c"

var Cookie = "ss-pid=ydoE0Vy8O6UXsERPqWsP; gr_user_id=4263d177-29bc-4dc9-aeb0-b4485b163e04; uniqueid=9a8d7885-1387-45e9-bd8d-ec19c542fd37; ppd_uname=pdu6867718253; ppdaiRole=4; __fp=fp; __vid=503543044.1521027645065; _ppdaiWaterMark=15238774245701; Hm_lvt_aab1030ecb68cd7b5c613bd7a5127a40=1524653748; _ga=GA1.2.1432998794.1524653748; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22pdu6867718253%22%2C%22%24device_id%22%3A%2216227e173a5877-03ee1ec19944a2-8343565-2073600-16227e173a6b28%22%2C%22first_id%22%3A%2216227e173a5877-03ee1ec19944a2-8343565-2073600-16227e173a6b28%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC(%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80)%22%7D%7D; regSourceId=0; referID=0; fromUrl=; referDate=2018-6-25%208%3A30%3A43; openid=5d598aa843d5071c20c8db4a2a4e0255; __vsr=1528767208188.src%3Ddirect%7Cmd%3Ddirect%7Ccn%3Ddirect%3B1528849526683.refSite%3Dhttps%3A//tz.ppdai.com/account/indexV2%7Cmd%3Dreferral%7Ccn%3Dreferral%3B1530590807215.refSite%3Dhttps%3A//tz.ppdai.com/menu/statistic-scatter%7Cmd%3Dreferral%7Ccn%3Dreferral%3B1530665196638.refSite%3Dhttps%3A//tz.ppdai.com/account/indexV2%7Cmd%3Dreferral%7Ccn%3Dreferral%3B1531137660312.src%3Ddirect%7Cmd%3Ddirect%7Ccn%3Ddirect; aliyungf_tc=AQAAANxeQSlwfgkAmaA6ewSBVIwRA8wY; Hm_lvt_f87746aec9be6bea7b822885a351b00f=1531270497,1531271054,1531281708,1531301374; token=2bd987685941fba2e8a26f7ff7093f576506f0e42b5c1e055c6a6685d603783495a352f87981ca3d41; __eui=t5yUdQAEOE11XwLJ1ZeKUQ%3D%3D; __tsid=119078225; ss-id=WsITymT6hBru0umGY9uT; gr_cs1_eeee18c2-105a-4307-ab45-87d095697d9f=user_name%3Afell_2015; gr_session_id_b9598a05ad0393b9_eeee18c2-105a-4307-ab45-87d095697d9f=true; gr_session_id_b9598a05ad0393b9=84866c3a-4026-4899-9700-2a891b933f0e; gr_session_id_b9598a05ad0393b9_84866c3a-4026-4899-9700-2a891b933f0e=true; currentUrl=https%3A%2F%2Finvdebt.ppdai.com%2Ftransferring%2Fcancel; __sid=1531308219053.22.1531309333301; Hm_lpvt_f87746aec9be6bea7b822885a351b00f=1531309333"


func init() {
	file, _ := os.OpenFile("cookie.txt", os.O_RDONLY, 0)
	data, _ := ioutil.ReadAll(file)
	Cookie = string(data)


	configData, _ := ioutil.ReadFile("server.config")
	json.Unmarshal(configData, &ServerConfig)
}
