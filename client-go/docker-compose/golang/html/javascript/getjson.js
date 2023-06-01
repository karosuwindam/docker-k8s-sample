var HOSTURL = ""

function getjson(output){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = table(data);
        // document.getElementById(output).innerHTML = data;
      }
    };
    var url = HOSTURL + "/v1/json"
    req.open("GET",url,true);
    req.send();
}

function table(json){
    var output = ""
    var tmp = json.Data
    for(var i=0;i<tmp.length;i++){
        if (tmp[i].Domain != ""){
            output += "<p><b>"+tmp[i].Domain+"</b></p>"
        }else{
            output += "<p></p>"
        }
        if (tmp[i].Data == null){
            continue
        }
        for(var j=0;j<tmp[i].Data.length;j++){
            output += "<div>"
            if (tmp[i].Domain != ""){
                if (tmp[i].Data[j].Podinfo != ""){
                    output += tmp[i].Data[j].Podinfo
                }else{
                    output += tmp[i].Data[j].Podname
                }
                if (tmp[i].Data[j].URL != ""){
                    output += " " + "<a href='"+"http://"+tmp[i].Data[j].URL+"'>"+"http</a>"
                    output += " " + "<a href='"+"https://"+tmp[i].Data[j].URL+"'>"+"https</a>"
                    output += " " + tmp[i].Data[j].URL
                }
            }else{
                if (tmp[i].Data[j].Podinfo != ""){
                    output += tmp[i].Data[j].Podname
                    output += " "+tmp[i].Data[j].Podinfo
                }else{
                    output += tmp[i].Data[j].Podname
                }
                if ((tmp[i].Data[j].Hostnetwork)&&(tmp[i].Data[j].Port != null)){
                    if (tmp[i].Data[j].Port.length != 0){
                    for (var k=0;k<tmp[i].Data[j].Port.length;k++){
                        output += "<a href='http://"+tmp[i].Data[j].Ip+":" + tmp[i].Data[j].Port[k]+"'>"+tmp[i].Data[j].Ip+":" + tmp[i].Data[j].Port[k]+"</a> "
                    }}
                }

            }
            if (tmp[i].Data[j].ClusterIP != "") {
                output += " ClusterIP :"+tmp[i].Data[j].ClusterIP
            }
            output += "</div>"        }
    }
    
    return output
}