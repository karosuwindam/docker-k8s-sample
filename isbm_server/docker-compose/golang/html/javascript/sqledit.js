var meta_suburl;

function getJSONlist(output) {
  var req = new XMLHttpRequest();
  req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
          var data=req.responseText;
          var tmp = JSON.parse(data)
          console.log(tmp)
          document.getElementById(output).innerHTML="<div>処理時間:"+tmp.Time+"</div>"+outputdiv_sqllist(tmp.Data);
      }
  };
  req.open("GET","./v1/sql/view/",false);
  req.send(null);  
}
function outputdiv_sqllist(json){
  var output ="<div clase=\"table\">"
  var title = ["ID","ISBN","Title","Writer","Brand","Ext","Image"]
  output += "<div class=\"tableRow\">"
  for (var i=0;i<title.length;i++){
      output += "<div class=\"tableHead\">"+title[i]+"</div>"
  }
  output += "</div>"
  for (var i=0;i<json.length;i++){
      output += "<div class=\"tableRow\">"
      output += "<div class=\"tableCell\">"+json[i].Id+"</div>"
      output += "<div class=\"tableCell\">"+json[i].Isbn+"</div>"
      output += "<div class=\"tableCell\">"+json[i].Title+"</div>"
      output += "<div class=\"tableCell\">"+json[i].Writer+"</div>"
      output += "<div class=\"tableCell\">"+json[i].Brand+"</div>"
      output += "<div class=\"tableCell\">"+json[i].Ext+"</div>"
      output += "<div class=\"tableCell\">"+"<img class='imagelist' src='"+ json[i].Image+"' alt='"+json[i].Image+"'>"+"</div>"
      output += "<div class=\"tableCellLess\">"+editurltext(json[i].Id)+"</div>"
      output += "<div class=\"tableCellLess\">"+destoryurltext(json[i].Id)+"</div>"
      output += "</div>"
  }
  return output
}

function editurltext(id){
  var output;
  output ="<a href=\"./sql/edit/"+meta_suburl+"/"+id+"\">";
  output += "edit";
  output += "</a>";
  return output;
};


function destoryurltext(id){
  var output;
  output = "<a href='javascript:destory("+id+");'>"
  output += "destory";
  output += "</a>";
  return output;
};

function destory(id){
    myRet = confirm("destory id="+id+" OK??");
    if (myRet){
      var xhr = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
      xhr.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
        if(xhr.readyState == 4 && xhr.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
            var data = xhr.responseText;
            console.log(data);		          // 取得した ファイルの中身を表示
            document.getElementById("answer").innerHTML = "destory id=" + id + " OK"
            getJSONlist("output");
        }
      };
  
      var url = "/sql/destory/"
      if (meta_suburl != ""){
        url += meta_suburl + "/"
      }
      url += id;
      xhr.open('POST', url, true);
      xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded;charset=UTF-8');
      xhr.send(null);
    }

}