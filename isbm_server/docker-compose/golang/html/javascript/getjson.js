var keyword_back

function getJSON(output,flag){
    var req = new XMLHttpRequest();
    var keyword = document.getElementById("keyword").value
    if ((keyword == keyword_back)&&(flag==false)){
        return
    }
    if (keyword == ""){
        document.getElementById(output).innerHTML = ""
        keyword_back = keyword
        return
    }
    req.onreadystatechange = function(){
        if(req.readyState == 4 && req.status == 200){
            nowserchpage = 1
            var data=req.responseText;
            var tmp = JSON.parse(data)
            console.log(tmp)
            if (tmp.Data.length == 1){
                document.getElementById(output).innerHTML = "<div>処理時間:"+tmp.Time+"</div>"+outputdiv(tmp.Data[0]);
            }else{
                document.getElementById(output).innerHTML = "<div>処理時間:"+tmp.Time+"</div>"+outputdiv_list(tmp.Data);
            }
            keyword_back = keyword
        }
    };
    req.open("GET","./v1/isbn/"+keyword,false);
    req.send(null);
}

function outputdiv_list(json){
    var output ="<div clase=\"table\">"
    var title = ["ID","ISBN","Title"]
    output += "<div class=\"tableRow\">"
    for (var i=0;i<title.length;i++){
        output += "<div class=\"tableHead\">"+title[i]+"</div>"
    }
    output += "</div>"
    for (var i=0;i<json.length;i++){
        output += "<div class=\"tableRow\">"
        output += "<div class=\"tableCell\">"+json[i].Id+"</div>"
        output += "<div class=\"tableCell\">"+"<a href=\"javascript:void(0);\" onclick=\"outputkey(this)\">"+json[i].Isbn+"</a>"+"</div>"
        output += "<div class=\"tableCell\">"+json[i].Title+"</div>"
        output += "</div>"
    }
    return output
}
function outputkey(str){
    document.getElementById("keyword").value = str.innerHTML
    getJSON(`output`,true)
    // console.log(str.innerHTML)
}
function outputdiv(json){
    var output = "<div clase=\"table\">"
    var title = ["ISBN","Title","作者","出版","あらすじ","その他"]
    var bodydata = [json.Isbn,json.Title,json.Writer,json.Brand,json.Synopsis,json.Ext]
    output += "<div class=\"tableRow\">"
    output += "<div class=\"tableHead\">Name</div>"
    output += "<div class=\"tableHead\">Data</div>"
    output += "</div>"
    for (var i=0;i<title.length;i++){
        output += "<div class=\"tableRow\">"
        output += "<div class=\"tableCell\">"+title[i]+"</div>"
        output += "<div class=\"tableCell\">"+bodydata[i]+"</div>"
        output += "</div>"
    }
    output += "<div class=\"tableRow\">"
    output += "<div class=\"tableCell\">"+"画像"+"</div>"
    output += "<div class=\"tableCell\">"+"<img src='"+ json.Image+"' alt='"+json.Image+"'>"+"</div>"
    output += "</div></div>"
    return output
}