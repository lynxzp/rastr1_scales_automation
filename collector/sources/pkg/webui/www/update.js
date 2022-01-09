window.setInterval(update,1000)

let rows
let names=[['ЛК-4', 'песок(0-5)'], ['ЛК-6', '0*40', '0*70', '20*40', '20*70'], ['ЛК-7', '5*10'], ['ЛК-8', '10*20', '5*20'],
    ['ЛК-9', 'песок(0-5)'], ['ЛК-14', '5*10'], ['ЛК-15', '10*20'], ['ЛК-17', '5*20']]
let firstLoad = true

window.onload = function () {
    rows = document.getElementsByClassName("rows")[0]
    names.forEach((element, index) => {addrow(element, index)})
    selectTab("tseha")
}

function update() {
    updateRequest()
    refreshReports()
}

function updateRequest() {
    httpUpdateRequest = new XMLHttpRequest();
    if (!httpUpdateRequest) {
        alert('Giving up :( Cannot create an XMLHTTP instance');
        return false;
    }
    httpUpdateRequest.timeout = 2000
    httpUpdateRequest.ontimeout = function (e) {
        document.getElementsByTagName("footer")[0].innerText="Таймаут соединения с сервером"
    }
    httpUpdateRequest.onreadystatechange = updateResponse;

    httpUpdateRequest.open('GET', 'ajax_update');
    httpUpdateRequest.send();
}

let stopChanging = false

function updateResponse() {
    if (httpUpdateRequest.readyState === XMLHttpRequest.DONE) {
        if (httpUpdateRequest.status === 200) {
            document.getElementsByTagName("footer")[0].innerText=""
            let params
            try {
                params = JSON.parse(httpUpdateRequest.responseText)
            } catch (e) {
                console.log(params)
                console.log(e)
                return
            }
            for(let i=0;i<names.length;i++){
                let DataPerfValue = params[i].DataPerfValue / 10
                let DataAccumValue = params[i].DataAccumValue
                if(DataPerfValue>=0) {
                    document.getElementById("DataPerfValue" + i).innerText = DataPerfValue
                }
                if(DataAccumValue>=0) {
                    document.getElementById("DataAccumValue" + i).innerText = DataAccumValue
                }
                if ((params[i].fraction.length > 0) && (!stopChanging)) {
                    document.getElementById("fraction" + i).childNodes[1].previousElementSibling.value = params[i].fraction
                }
                if(isLocalhost()){
                    document.getElementById("requests" + i).innerText = params[i].requests
                    document.getElementById("responses" + i).innerText = params[i].responses
                    if(firstLoad===true) {
                        document.getElementById("rs485addr" + i).getElementsByTagName("input")[0].value = params[i].rs485addr
                        document.getElementById("ipaddr" + i).getElementsByTagName("input")[0].value = params[i].ipaddr
                        document.getElementById("dtype" + i).getElementsByTagName("input")[0].value = "0x" + parseInt(params[i].dtype).toString(16)
                    }
                }
            }
            firstLoad = false
        } else {
            document.getElementsByTagName("footer")[0].innerText="Нет соединения с сервером (статус: " +
                httpUpdateRequest.status + ")"
        }
    }
}

function addrow(name, i) {
    // name
    let d = document.createElement("div");
    d.setAttribute('id', 'name'+i)
    d.innerHTML = name[0]
    rows.appendChild(d)
    // fraction
    d = document.createElement("div");
    d.setAttribute('id', 'fraction'+i)
    d.className+='dropdown'
    let fractions = name.slice(1, name.length)
    let fractions_select = ''
    fractions.forEach((fr) => {fractions_select+='<option>'+fr+'</option>'})
    d.innerHTML = "<input type='text'>" +
        "<select onchange='this.previousElementSibling.value=this.value; this.previousElementSibling.focus();stopChanging=true;' autofocus='true'>" +
        fractions_select + fractions[0] +
        "</select>\n"
    d.value=fractions[0]
    rows.appendChild(d)
    document.getElementById('fraction'+i).getElementsByTagName("input")[0].value=fractions[0]
    // Save Fraction
    d = document.createElement("div");
    d.innerHTML = '<input type="button" value="Сохранить" id="save'+i+'" onclick="saveFraction('+i+')">'
    rows.appendChild(d)
    // DataAccumValue
    d = document.createElement("div");
    d.setAttribute('id', 'DataAccumValue'+i)
    rows.appendChild(d)
    // DataPerfValue
    d = document.createElement("div");
    d.setAttribute('id', 'DataPerfValue'+i)
    rows.appendChild(d)
    if(isLocalhost()) {
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
        d.setAttribute('id', 'rs485addr'+i)
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
        // Save button
        d = document.createElement("div");
        d.innerHTML = '<input type="button" value="Сохранить" id="save'+i+'" onclick="saveScale('+i+')">'
        rows.appendChild(d)
        // Clear button
        d = document.createElement("div");
        d.innerHTML = '<input type="button" value="Очистить" id="clear' + i + '" onclick="clearClick(' + i + ')">'
        rows.appendChild(d)
    }
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

    if(tabName === "tseha") {
        refreshReports()
    }
}


function sendRequest(url, params, errorText) {
    let request = new XMLHttpRequest();
    if (!request) {
        alert('Giving up :( Cannot create an XMLHTTP instance');
        return false;
    }
    request.timeout = 2000
    request.ontimeout = function (e) {
        document.getElementsByTagName("footer")[0].innerText=errorText
    }
    // httpSaveRequest.onreadystatechange = ajaxUpdate;

    request.open('GET', url+params);
    request.send();

}

function saveScale(i) {
    let params = "?id="+i
    if (isLocalhost()){
    params +="&dtype=" + document.getElementById("dtype"+i).getElementsByTagName("input")[0].value.slice(2,4)
    params += "&ipaddr=" + document.getElementById("ipaddr"+i).getElementsByTagName("input")[0].value
    params += "&rs485addr=" + document.getElementById("rs485addr"+i).getElementsByTagName("input")[0].value}
    params += "&fraction=" + document.getElementById("fraction"+i).getElementsByTagName("input")[0].value
    sendRequest("save_scale", params, "Ошибка сохранения")
    stopChanging = false
}

function saveFraction(i) {
    let params = "?id="+i
    params += "&fraction=" + document.getElementById("fraction" + i).getElementsByTagName("input")[0].value
    sendRequest("save_fraction", params, "Ошибка сохранения")
    stopChanging = false
}

function clearClick(i) {
    sendRequest("clear", "?id="+i)
}

function logout() {
    document.cookie = "password=; Max-Age=0"
    document.cookie = "login=; Max-Age=0"
    window.location.href = "/login.html"
}

function isLocalhost() {
    return (location.hostname === "localhost" || location.hostname === "127.0.0.1")
}