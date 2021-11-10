window.setInterval(update,200)

let rows
const nums = 10

window.onload = function () {
    rows = document.getElementsByClassName("rows")[0]
    for(let i=0;i<nums;i++) {
        addrow(i)
    }
}

function update() {
    let a = ""
    for(let i=0;i<nums;i++) {

    }
    ajaxRequest()
}

function ajaxRequest() {
    httpRequest = new XMLHttpRequest();
    if (!httpRequest) {
        alert('Giving up :( Cannot create an XMLHTTP instance');
        return false;
    }
    httpRequest.timeout = 2000
    httpRequest.ontimeout = function (e) {
        document.getElementsByTagName("footer")[0].innerText="Таймаут соединения с сервером"
    }
    httpRequest.onreadystatechange = ajaxUpdate;

    let params = "?"
    let separator = ""
    for(let i=0;i<nums;i++) {
        params += separator
        separator = "&"
        let id = "dtype"+i
        params += id + "=" + document.getElementById(id).getElementsByTagName("input")[0].value.slice(2,4)
        id = "ipaddr"+i
        params += "&" + id + "=" + document.getElementById(id).getElementsByTagName("input")[0].value
        id = "rs485addr"+i
        params += "&" + id + "=" + document.getElementById(id).getElementsByTagName("input")[0].value
    }
    httpRequest.open('GET', 'ajax_update'+params);
    httpRequest.send();
}

function ajaxUpdate() {
    if (httpRequest.readyState === XMLHttpRequest.DONE) {
        if (httpRequest.status === 200) {
            document.getElementsByTagName("footer")[0].innerText=""
            let params = JSON.parse(httpRequest.responseText)
            console.log(params)
            for(let i=0;i<nums;i++){
                let DataPerfValue = params[i].DataPerfValue / 10
                let DataAccumValue = params[i].DataAccumValue
                if(DataPerfValue>=0) {
                    document.getElementById("DataPerfValue" + i).innerText = DataPerfValue
                }
                if(DataAccumValue>=0) {
                    document.getElementById("DataAccumValue" + i).innerText = DataAccumValue
                }
                document.getElementById("requests" + i).innerText = params[i].requests
                document.getElementById("responses" + i).innerText = params[i].responses
            }
        } else {
            document.getElementsByTagName("footer")[0].innerText="Нет соединения с сервером (статус: " +
                httpRequest.status + ")"
        }
    }
}

function addrow(i) {
    // name
    let d = document.createElement("div");
    d.setAttribute('id', 'name'+i)
    d.innerHTML = "<input type='text' maxlength='25' size='25' value='ДЦ-"+(i+1)+"'>"
    rows.appendChild(d)
    // DataAccumValue
    d = document.createElement("div");
    d.setAttribute('id', 'DataAccumValue'+i)
    d.innerHTML = "-2"
    rows.appendChild(d)
    // DataPerfValue
    d = document.createElement("div");
    d.setAttribute('id', 'DataPerfValue'+i)
    d.innerHTML = "-2"
    rows.appendChild(d)
    // type
    d = document.createElement("div");
    d.setAttribute('id', 'dtype'+i)
    d.className+='dropdown'
    d.innerHTML = "<input type=\"text\" />\n" +
                  "<select  onchange=\"this.previousElementSibling.value=this.value; this.previousElementSibling.focus()\">\n" +
                  "<option>0x?? Свое значение</option>\n" +
                  // "<option>0x60 Накопление</option>\n" +
                  "<option>0x5d Производительность v2</option>\n" +
                  "<option>0x3f Производительность v1</option>\n" +
                  "</select>\n"
    rows.appendChild(d)
    // ip address
    d = document.createElement("div");
    d.setAttribute('id', 'ipaddr'+i)
    d.innerHTML = "<input type=\"text\" minlength=\"7\" maxlength=\"15\" size=\"15\" pattern=\"^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$\">"
    rows.appendChild(d)
    // rs-485 address
    d = document.createElement("div");
    d.setAttribute('id', 'rs485addr'+i)
    d.innerHTML = "<input type=\"text\" minlength=\"1\" maxlength=\"3\" size=\"3\" pattern=\"^\\d{1,3}$\">"
    rows.appendChild(d)
    // Requests
    d = document.createElement("div");
    d.setAttribute('id', 'requests'+i)
    d.innerHTML = "-2"
    rows.appendChild(d)
    // Responses
    d = document.createElement("div");
    d.setAttribute('id', 'responses'+i)
    d.innerHTML = "-2"
    rows.appendChild(d)
    // Table
    d = document.createElement("div");
    d.setAttribute('id', 'table'+i)
    d.setAttribute('class', 'table')
    d.innerHTML = tableHtml
    rows.appendChild(d)
}

const tableHtml = "<table><tbody><tr><th rowspan='2'>№ ЛК</th><th rowspan='2'>Фракция</th><th colspan='6'>Количество продукции</th></tr>" +
    "<tr><th>1 смена</th><th>2 смена</th><th>3 смена</th><th>За сутки</th><th>С начала месяца</th><th>С начала года</th></tr>" +
    "<tr><td>-</td><td>20x40</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td></tr>" +
    "<tr><td>-</td><td>5x10</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td></tr>" +
    "<tr><td>-</td><td>10x20</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td></tr>" +
    "<tr><td>-</td><td>отсев</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td><td>-</td></tr>" +
    "</tbody></table>"
