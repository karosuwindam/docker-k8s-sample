function timech(timedata){
    // var date = Date.parse(timedata);
    var date = new Date(timedata);
    
  
    var year_str = date.getFullYear();
    //月だけ+1すること
    var month_str = 1 + date.getMonth();
    var day_str = date.getDate();
    var hour_str = date.getHours();
    var minute_str = date.getMinutes();
    var second_str = date.getSeconds();
    
    
    month_str = ('0' + month_str).slice(-2);
    day_str = ('0' + day_str).slice(-2);
    hour_str = ('0' + hour_str).slice(-2);
    minute_str = ('0' + minute_str).slice(-2);
    second_str = ('0' + second_str).slice(-2);
    
    format_str = 'YYYY-MM-DD hh:mm:ss';
    format_str = format_str.replace(/YYYY/g, year_str);
    format_str = format_str.replace(/MM/g, month_str);
    format_str = format_str.replace(/DD/g, day_str);
    format_str = format_str.replace(/hh/g, hour_str);
    format_str = format_str.replace(/mm/g, minute_str);
    format_str = format_str.replace(/ss/g, second_str);
    return format_str;
}
function statusckJSON(output){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = "Book:"+data.BookStatus+"("+timech(data.BookNowTIme)+")" +",Bookmark:"+ data.BookMarkStatus+"("+timech(data.BookMarkNowTime)+")";
        // document.getElementById(output).innerHTML = data;
      }
    };
    req.open("GET","/status",true);
    req.send();

}

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
        output += timech(tmp[i].Lastdate)
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
function table_conbd(data){
    var output = ""
    output += "<div class='block_n1'>"+data.Days+ "</div>";
    output += "<div class='block_n2'>"+data.Type+ "</div>";
    output += "<div class='block_title'>"+data.Title+ "</div>";
    output += "<div class='block_n2'>"+data.Writer+ "</div>";
    output += "<div class='block_n2'>"+data.Bround+ "</div>";
    output += "<div class='block_n2'>"+data.Ext+ "</div>";
    output += "<div class='block_img'>"+"<img src='"+ data.Img+"' alt='"+data.Img+"'>"+"</div>";
    return output
}
function tableb(data){
    var tmp = JSON.parse(data);
    var rajiob = document.getElementsByName("type");
    var output = tmp.Year +"年" +tmp.Month + "月<br>";
    output += "</div>"
    
    if (rajiob[1].checked ){
        for (var i=0;i<tmp.LiteNobel.length;i++){
            output += "<div class='table_line'>"
            output += table_conbd(tmp.LiteNobel[i])
            output += "</div>"
        }
    }else{
        for (var i=0;i<tmp.Comic.length;i++){
            output += "<div class='table_line'>"
            output += table_conbd(tmp.Comic[i])
            output += "</div>"
        }
    }
    output += "<div class='table'>"
//     output +="<table>";
//     if (rajiob[1].checked ){
//         for (var i=0;i<tmp.LiteNobel.length;i++){
//         output += "<tr>"+table_conb(tmp.LiteNobel[i])+"</tr>"
//         }
//     }else{
//         for (var i=0;i<tmp.Comic.length;i++){
//              output += "<tr>"+table_conb(tmp.Comic[i])+"</tr>";
//          }
//     }
// output += "</table>";
    return output
}