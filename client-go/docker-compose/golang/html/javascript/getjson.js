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
    req.open("GET","/json",true);
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
        for(var j=0;j<tmp[i].Data.length;j++){
            output += "<div>"
            if (tmp[i].Domain != ""){
                if (tmp[i].Data[j].Podinfo != ""){
                    output += tmp[i].Data[j].Podinfo
                }else{
                    output += tmp[i].Data[j].PName
                }
                output += " " + "<a href='"+"http://"+tmp[i].Data[j].URL+"'>"+"http</a>"
                output += " " + "<a href='"+"https://"+tmp[i].Data[j].URL+"'>"+"https</a>"
                output += " " + tmp[i].Data[j].URL
            }else{
                if (tmp[i].Data[j].Podinfo != ""){
                    output += tmp[i].Data[j].PName
                    output += " "+tmp[i].Data[j].Podinfo
                }else{
                    output += tmp[i].Data[j].PName
                }
                if (tmp[i].Data[j].HostNetwork){
                    if (tmp[i].Data[j].Port.length != 0){
                    for (var k=0;k<tmp[i].Data[j].Port.length;k++){
                        output += "<a href='http://"+tmp[i].Data[j].Ip+":" + tmp[i].Data[j].Port[k]+"'>"+tmp[i].Data[j].Ip+":" + tmp[i].Data[j].Port[k]+"</a>"
                    }}
                }

            }
            output += "</div>"        }
    }
    
    return output
}