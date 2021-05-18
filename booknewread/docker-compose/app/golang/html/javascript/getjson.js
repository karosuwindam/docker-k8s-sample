function serchgetJSON(output){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        // data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = table(data);
        // document.getElementById(output).innerHTML = data;
      }
    };
    req.open("GET","/json",true);
    req.send();
}
function getnobleJSON(output){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        // data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = table_noble(data);
        // document.getElementById(output).innerHTML = data;
      }
    };
    req.open("GET","/jsonnobel",true);
    req.send();
}
function getnewJSON(output,page){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        // data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = tableb(data);
        // document.getElementById(output).innerHTML = data;
      }
    };
    req.open("POST","/jsonb",true);
    req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    var str = "page="+page
    req.send(str);
}

function table_noble(data){
    var output = "";
    var tmp = JSON.parse(data);
    console.log(tmp);
    output += "<div class='table'>"
    for (var i=0;i<tmp.length;i++){
        output += "<div class='table_line'>"
        output += "<div class='block'>"
        output += "<a href='"+tmp[i].Url+"'>"+tmp[i].Title+"</a>"
        output += "</div>"
        output += "<div class='block'>"
        output += "<a href='"+tmp[i].LastUrl+"'>"+tmp[i].LastStoryT+"</a>"
        output += "</div>"
        output += "<div class='block'>"
        output += tmp[i].Lastdate
        output += "</div>"
        output += "</div>"
    }
    output += "</div>"
    return output;
}

function table_con(data){
    var output = "<tr>";
    for (var i=0;i<data.length;i++){
        output += "<td>"+data[i]+"</td>"
    }
    output += "</tr>";
    return output
}
function table(data){
    var tmp = JSON.parse(data);
    var output ="<table>";
    output += "<tr>"+table_con(tmp.column)+"</tr>"
    for (var i=0;i<tmp.list.length;i++){
        output += "<tr>"+table_con(tmp.list[i])+"</tr>"
    }
    output += "</table>";
    return output
}
function table_conb(data){
    var output = "<tr>";
    output += "<td>"+data.Days+"</td>"
    output += "<td>"+data.Type+"</td>"
    output += "<td>"+data.Title+"</td>"
    output += "<td>"+data.Writer+"</td>"
    output += "<td>"+data.Bround+"</td>"
    output += "<td>"+data.Ext+"</td>"
    output += "<td>"+"<img src='"+ data.Img+"' alt='"+data.Img+"'>"+"</td>"

    output += "</tr>";
    return output
}
function tableb(data){
    var tmp = JSON.parse(data);
    var rajiob = document.getElementsByName("type");
    var output = tmp.Year +"年" +tmp.Month + "月<br>";
    output +="<table>";
    if (rajiob[1].checked ){
        for (var i=0;i<tmp.LiteNobel.length;i++){
        output += "<tr>"+table_conb(tmp.LiteNobel[i])+"</tr>"
        }
    }else{
        for (var i=0;i<tmp.Comic.length;i++){
             output += "<tr>"+table_conb(tmp.Comic[i])+"</tr>";
         }
    }
output += "</table>";
    return output
}