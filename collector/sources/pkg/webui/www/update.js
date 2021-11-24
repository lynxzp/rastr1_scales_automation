// window.setInterval(update,200)

let rows
let names=[['ЛК-4', 'песок(0*5)'], ['ЛК-6', '0*40', '0*70', '20*40', '20*70'], ['ЛК-7', '5*10'], ['ЛК-8', '10*20', '5*20'],
    ['ЛК-9', 'песок(0-5)'], ['ЛК-14', '5*10'], ['ЛК-15', '10*20'], ['ЛК-17', '5*20']]

window.onload = function () {
    rows = document.getElementsByClassName("rows")[0]
    names.forEach((element, index) => {addrow(element, index)})
}

function update() {
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
    for(let i=0;i<names.length;i++) {
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
            for(let i=0;i<names.length;i++){
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

function addrow(name, i) {
    // name
    let d = document.createElement("div");
    d.setAttribute('id', 'name'+i)
    d.innerHTML = "<input type='text' maxlength='10' size='1' value='"+name[0]+"'>"
    rows.appendChild(d)
    // fraction
    d = document.createElement("div");
    d.setAttribute('id', 'fraction'+i)
    d.className+='dropdown'
    let fractions = name.slice(1, name.length)
    let fractions_select = ''
    fractions.forEach((fr) => {fractions_select+='<option>'+fr+'</option>'})
    d.innerHTML = "<input type='text'>" +
        "<select onchange='this.previousElementSibling.value=this.value; this.previousElementSibling.focus()' autofocus='true'>" +
        fractions_select + fractions[0] +
        "</select>\n"
    d.value=fractions[0]
    rows.appendChild(d)
    document.getElementById('fraction'+i).getElementsByTagName("input")[0].value=fractions[0]
    // DataAccumValue
    d = document.createElement("div");
    d.setAttribute('id', 'DataAccumValue'+i)
    rows.appendChild(d)
    // DataPerfValue
    d = document.createElement("div");
    d.setAttribute('id', 'DataPerfValue'+i)
    rows.appendChild(d)
    // type
    d = document.createElement("div");
    d.setAttribute('id', 'dtype'+i)
    d.className+='dropdown'
    d.innerHTML = "<input type='text' />" +
                  "<select  onchange='this.previousElementSibling.value=this.value; this.previousElementSibling.focus()'>" +
                  "<option>0x?? Свое значение</option>" +
                  "<option>0x5d Производительность v2</option>" +
                  "<option>0x3f Производительность v1</option>" + "sdsdsd" +
                  "</select>"
    rows.appendChild(d)
    // ip address
    d = document.createElement("div");
    d.setAttribute('id', 'ipaddr'+i)
    d.innerHTML = "<input type=\"text\" minlength=\"7\" maxlength=\"15\" size=\"15\" pattern=\"^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$\">"
    rows.appendChild(d)
    // rs-485 address
    d = document.createElement("div");
    d.setAttribute('id', 'rs485adErrordr'+i)
    d.innerHTML = "<input type=\"text\" minlength=\"1\" maxlength=\"3\" size=\"3\" pattern=\"^\\d{1,3}$\">"
    rows.appendChild(d)
    // Requests
    d = document.createElement("div");
    d.setAttribute('id', 'requests'+i)
    rows.appendChild(d)
    // Responses
    d = document.createElement("div");
    d.setAttribute('id', 'responses'+i)
    rows.appendChild(d)
}

function selectTab(tabName) {
    // Declare all variables
    let i, tabcontent, tablinks;

    // Get all elements with class="tabcontent" and hide them
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    // Get all elements with class="tablinks" and remove the class "active"
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    // Show the current tab, and add an "active" class to the button that opened the tab
    document.getElementById(tabName).style.display = "block";
    event.currentTarget.className += " active";
}