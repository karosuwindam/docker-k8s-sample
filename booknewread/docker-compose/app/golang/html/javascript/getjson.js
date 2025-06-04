var JSON_DATA = ["", "", ""];
var timer = "";

var serche = function (keyword, output) {
  console.log(
    keyword,
    document.getElementById("page").value,
    document.getElementById("pagetype").value
  );
  var page = document.getElementById("page").value - 0;
  var pagetype = document.getElementById("pagetype").value - 0;
  if (JSON_DATA[page - 0] == "") {
    if (pagetype == 2) {
    } else {
      getnewJSON(output, page);
    }
  }
  if (page == 3) {
    return;
  }
  if (keyword == "") {
    // console.log(tableb(JSON_DATA[page-0]));
    document.getElementById(output).innerHTML = tableb(JSON_DATA[page - 0]);
  } else {
    // console.log(serche_table(JSON_DATA[page-0],keyword.toUpperCase() ))

    if (pagetype == 2) {
    } else {
      document.getElementById(output).innerHTML = serche_table(
        JSON_DATA[page - 0],
        keyword.toUpperCase()
      );
    }
    // document.getElementById(output).innerHTML = tableb(JSON_DATA[0]);
  }
};

function serche_table(data, keyword) {
  var output = "";
  var tmp = JSON.parse(data);
  var tmp_data = [];
  // ToDo
  // キーワードによる検索をまとめる
  // キーワード結果をHTMLに加工
  // for (var i=0;i<tmp.Comic.length;i++){
  //     if (data_ck(keyword,tmp.Comic[i])){
  //         tmp_data.push(tmp.Comic[i])
  //     }
  // }
  // for (var i=0;i<tmp.LiteNobel.length;i++){
  //     if (data_ck(keyword,tmp.LiteNobel[i])){
  //         tmp_data.push(tmp.LiteNobel[i])
  //     }
  // }
  // var output = tmp.Year +"年" +tmp.Month + "月<br>";
  // output += "</div>"
  // for (var i=0;i<tmp_data.length;i++){
  //     output += "<div class='table_line'>"
  //     output += table_conbd(tmp_data[i])
  //     output += "</div>"
  // }
  // output += "<div class='table'>"
  return output;
}
function hankaku2Zenkaku(str) {
  return str.replace(/[Ａ-Ｚａ-ｚ０-９]/g, function (s) {
    return String.fromCharCode(s.charCodeAt(0) - 0xfee0);
  });
}

function data_ck(keyword, data) {
  var output = false;
  var ckdata = [data.Bround, data.Ext, data.Title, data.Writer];
  for (var i = 0; i < ckdata.length; i++) {
    if (ckdata[i] == "") {
      continue;
    }
    str = hankaku2Zenkaku(ckdata[i]);
    if (str.indexOf(keyword) != -1) {
      output = true;
      break;
    }
  }
  return output;
}

function serche_keyword(keyword, output) {
  if (timer != "") {
    clearTimeout(timer);
  }
  timer = setTimeout(serche, 500, keyword, output);
}

function isSmartPhone() {
  if (navigator.userAgent.match(/iPhone|Android.+Mobile/)) {
    return true;
  } else {
    return false;
  }
}
function timech(timedata) {
  // var date = Date.parse(timedata);
  var date = new Date(timedata);

  var year_str = date.getFullYear();
  //月だけ+1すること
  var month_str = 1 + date.getMonth();
  var day_str = date.getDate();
  var hour_str = date.getHours();
  var minute_str = date.getMinutes();
  var second_str = date.getSeconds();

  month_str = ("0" + month_str).slice(-2);
  day_str = ("0" + day_str).slice(-2);
  hour_str = ("0" + hour_str).slice(-2);
  minute_str = ("0" + minute_str).slice(-2);
  second_str = ("0" + second_str).slice(-2);

  format_str = "YYYY-MM-DD hh:mm:ss";
  format_str = format_str.replace(/YYYY/g, year_str);
  format_str = format_str.replace(/MM/g, month_str);
  format_str = format_str.replace(/DD/g, day_str);
  format_str = format_str.replace(/hh/g, hour_str);
  format_str = format_str.replace(/mm/g, minute_str);
  format_str = format_str.replace(/ss/g, second_str);
  return format_str;
}

var statusRun = function(output) {
  statusckJSON(output)
}
var timer = nil
function statusckJSON(output) {
  clearInterval(timer);
 
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function () {
    if (req.readyState == 4 && req.status == 200) {
      var data = req.responseText;
      data = JSON.parse(data);
      console.log(data);
      document.getElementById(output).innerHTML =
        "Book:" +
        statusLoad(data.BookStatus) +
        "(" +
        timech(data.BookNowTIme) +
        ")" +
        "<br>Bookmark:" +
        statusLoad(data.BookMarkStatus) +
        "(" +
        timech(data.BookMarkNowTime) +
        ")";
        if (data.BookMarkStatus!="ok") {
          timer = setInterval(statusRun, 3000, output);
        }else {
          if (document.getElementById("pagetype").value == "2"){
            getnobleJSON("output");
          }
          timer = nil
        }
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("GET", "/v1/status", true);
  req.send();
}

function statusLoad(value) {
  if (value=="ok"){
    return value
  }
  var num = value.substr("Reload:".length,value.length-1-"Reload:".length) - 0
  return "<progress max='100'value='"+num+"' >"+num+"</progress>"
}
function serchgetJSON(output) {
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function () {
    if (req.readyState == 4 && req.status == 200) {
      var data = req.responseText;
      // data = JSON.parse(data);
      console.log(data);
      document.getElementById(output).innerHTML = table(data);
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("GET", "/v1/json", true);
  req.send();
}
function getnobleJSON(output) {
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function () {
    if (req.readyState == 4 && req.status == 200) {
      var data = req.responseText;
      // data = JSON.parse(data);
      console.log(JSON.parse(data));
      JSON_DATA[JSON_DATA.length - 1] = data;
      document.getElementById(output).innerHTML = table_noble(data);
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("GET", "/v1/json/nobel", true);
  req.send();
}
function getnewJSON(output, page) {
  var rajiom = document.getElementsByName("month");
  if (rajiom[rajiom.length - 1].checked) {
    if (JSON_DATA[0] != "") {
      document.getElementById(output).innerHTML = tableb(JSON_DATA[0]);
      return;
    }
  } else {
    if (JSON_DATA[page - 0] != "") {
      document.getElementById(output).innerHTML = tableb(JSON_DATA[page - 0]);
      return;
    }
  }
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function () {
    if (req.readyState == 4 && req.status == 200) {
      var data = req.responseText;
      // data = JSON.parse(data);
      if (rajiom[rajiom.length - 1].checked) {
        JSON_DATA[0] = data;
      } else {
        JSON_DATA[page - 0] = data;
      }
      console.log(JSON.parse(data));
      document.getElementById(output).innerHTML = tableb(data);
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("GET", "/v1/json/book/" + page, true);
  req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  req.send();
}

function serche_table_nobel(data, keyword) {
  var output = "";
  var tmp = JSON.parse(data);
  var tmp_data = [];

  for (var i = 0; i < tmp.Comic.length; i++) {
    if (data_ck(keyword, tmp.Comic[i])) {
      tmp_data.push(tmp.Comic[i]);
    }
  }
  for (var i = 0; i < tmp.LiteNobel.length; i++) {
    if (data_ck(keyword, tmp.LiteNobel[i])) {
      tmp_data.push(tmp.LiteNobel[i]);
    }
  }
  var output = tmp.Year + "年" + tmp.Month + "月<br>";
  output += "</div>";
  for (var i = 0; i < tmp_data.length; i++) {
    output += "<div class='table_line'>";
    output += table_conbd(tmp_data[i]);
    output += "</div>";
  }
  output += "<div class='table'>";
  return output;
}

function table_noble(data) {
  var output = "";
  var tmp = JSON.parse(data);
  console.log(tmp);
  output += "<div class='table'>";
  for (var i = 0; i < tmp.length; i++) {
    output += table_nobel_con(tmp[i]);
  }
  output += "</div>";
  return output;
}

function table_nobel_con(data) {
  var output = "";
  output += "<div class='table_line'>";
  output += "<div class='block'>";
  output += "<a href='" + data.Url + "'>" + data.Title + "</a>";
  output += "</div>";
  output += "<div class='block'>";
  output += "<a href='" + data.LastUrl + "'>" + data.LastStoryT + "</a>";
  output += "</div>";
  output += "<div class='block_time'>";
  output += timech(data.Lastdate);
  output += "</div>";
  output += "</div>";
  return output;
}

function table_con(data) {
  var output = "<tr>";
  for (var i = 0; i < data.length; i++) {
    output += "<td>" + data[i] + "</td>";
  }
  output += "</tr>";
  return output;
}
function table(data) {
  var tmp = JSON.parse(data);
  var output = "<table>";
  output += "<tr>" + table_con(tmp.column) + "</tr>";
  for (var i = 0; i < tmp.list.length; i++) {
    output += "<tr>" + table_con(tmp.list[i]) + "</tr>";
  }
  output += "</table>";
  return output;
}
// function table_conb(data){
//     var output = "<tr>";
//     output += "<td>"+data.Days+"</td>"
//     output += "<td>"+data.Type+"</td>"
//     output += "<td>"+data.Title+"</td>"
//     output += "<td>"+data.Writer+"</td>"
//     output += "<td>"+data.Bround+"</td>"
//     output += "<td>"+data.Ext+"</td>"
//     output += "<td>"+"<img src='"+ data.Img+"' alt='"+data.Img+"'>"+"</td>"

//     output += "</tr>";
//     return output
// }
function table_conbd(data) {
  var output = "";
  output += "<div class='block_n1'>" + data.Days + "</div>";
  output += "<div class='block_n2'>" + data.Type + "</div>";
  output += "<div class='block_title'>" + data.Title + "</div>";
  output += "<div class='block_n2'>" + data.Writer + "</div>";
  output += "<div class='block_n2'>" + data.Bround + "</div>";
  output += "<div class='block_n2'>" + data.Ext + "</div>";
  if (window.navigator.connection.type == "cellular") {
    output +=
      "<div class='block_img' onclick='testclick(this,\"" +
      data.Img +
      "\")'>" +
      "onclick<br>image<br>view</div>";
  } else {
    output +=
      "<div class='block_img'>" +
      "<img src='" +
      data.Img +
      "' alt='" +
      data.Img +
      "'>" +
      "</div>";
  }
  return output;
}

function testclick(data, iamgeurl) {
  // alert(data,iamgeurl);
  data.innerHTML = "<img src='" + iamgeurl + "' alt='" + iamgeurl + "'>";
}
function table_conbd_m(data) {
  var output = "";
  output += "<div class='block_n1'>" + data.Days + "</div>";
  output += "<div class='block_n2'>" + data.Type + "</div>";
  output +=
    "<div class='block_title'>" +
    data.Title +
    "<br>" +
    data.Writer +
    " " +
    data.Bround +
    " " +
    data.Ext +
    "</div>";
  // output += "<div class='block_img'>"+"<img src='"+ data.Img+"' alt='"+data.Img+"'>"+"</div>";
  if (window.navigator.connection.type == "cellular") {
    output +=
      "<div class='block_img' onclick='testclick(this,\"" +
      data.Img +
      "\")'>" +
      "onclick<br>image<br>view</div>";
  } else {
    output +=
      "<div class='block_img'>" +
      "<img src='" +
      data.Img +
      "' alt='" +
      data.Img +
      "'>" +
      "</div>";
  }

  return output;
}
function tableb(data) {
  var tmp = JSON.parse(data);
  var rajiob = document.getElementsByName("type");
  var rajiom = document.getElementsByName("month");
  var output = tmp.Year + "年" + tmp.Month + "月<br>";
  output += "</div>";
  if (rajiom[rajiom.length - 1].checked) {
    var date = new Date();
    var day = date.getDate();
    for (var i = 0; i < tmp.Comic.length; i++) {
      output += "<div class='table_line'>";
      if (isSmartPhone()) {
        if (tmp.Comic[i].Days == day) {
          output += table_conbd_m(tmp.Comic[i]);
        }
      } else {
        if (tmp.Comic[i].Days == day) {
          output += table_conbd(tmp.Comic[i]);
        }
      }
      output += "</div>";
    }
    for (var i = 0; i < tmp.LiteNobel.length; i++) {
      output += "<div class='table_line'>";
      if (isSmartPhone()) {
        if (tmp.LiteNobel[i].Days == day) {
          output += table_conbd_m(tmp.LiteNobel[i]);
        }
      } else {
        if (tmp.LiteNobel[i].Days == day) {
          output += table_conbd(tmp.LiteNobel[i]);
        }
      }
      output += "</div>";
    }
  } else {
    if (rajiob[1].checked) {
      for (var i = 0; i < tmp.LiteNobel.length; i++) {
        output += "<div class='table_line'>";
        if (isSmartPhone()) {
          output += table_conbd_m(tmp.LiteNobel[i]);
        } else {
          output += table_conbd(tmp.LiteNobel[i]);
        }
        output += "</div>";
      }
    } else {
      for (var i = 0; i < tmp.Comic.length; i++) {
        output += "<div class='table_line'>";
        if (isSmartPhone()) {
          output += table_conbd_m(tmp.Comic[i]);
        } else {
          output += table_conbd(tmp.Comic[i]);
        }
        output += "</div>";
      }
    }
  }
  output += "<div class='table'>";
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
  return output;
}
