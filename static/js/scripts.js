socket = new WebSocket("ws://" + window.location.host + "/ws")
var g_highscores = [];
var g_queue = [];
var currentView = "";
var fader = $("#fader");
apiView();

socket.onmessage = event => {
  let data = JSON.parse(event.data);
  let cnt = data.content
  switch (data.type) {
    case "highscores":
      g_highscores = cnt.highscores;
      if (currentView === "highscore"){
        highscoreView();
      }
      break;
    case "queue":
      g_queue = cnt.queue;
      if (currentView === "submit"){
        submitView();
      }
      break;
  }
  
}

function clearTable(){
  document.getElementById("header").innerHTML = "";
  document.getElementById("body").innerHTML = "";
  document.getElementById("topcontent").innerHTML = "";
}

function createHeaders(headerArr){
  let header = document.getElementById("header")
  headerArr.map((v ,i) => {
    let cell = document.createElement("th")
		cell.setAttribute("class", "cell100 column" + (i + 1))
    cell.append(v)
    header.append(cell)
  })
}

function createRow(row){
  let tr = document.createElement("tr")
  tr.setAttribute("class", "row100 body")
  row.map((v, i) => {
    let td = document.createElement("td")
    td.setAttribute("class", "cell100 column" + (i + 1))
    td.append(v);
    tr.append(td);
  })
  document.getElementById("body").append(tr)
}

function apiView(){
  currentView = "API";
  $("#Topbtn").attr("class", "");
  $("#APIbtn").attr("class", "active");
  $("#Solutionbtn").attr("class", "");
  $("#Submitbtn").attr("class", "");
  let table = $("#table")
  fader.fadeOut(200, () => {
    clearTable();
    createHeaders(["Namn", "Metod", "Header", "Body", "Response"]);
    $("#topcontent").append(`
    <div style="color:white;font-family:Lato-Regular;">
    <h style="font-family:Lato-Bold;">Instruktioner</h>
    <p>Parallellprogrammering handlar till mångt och mycket om att minimera dötid för processorn.<br> 
    Ett område där dötid naturligt uppstår är internetprogrammering. Varje gång ett anrop skickas
    måste vi vänta på att servern ska svara. </p>
    <p>Övningen går ut på att summera ett antal tal som skickas från en server.</p>
    <p>Servern går att ladda ned <a style="color:#00ad5f;" href="https://gist.github.com/Limpa/82c1b204e40892d4a57ec8d40bce94a6">här</a> för att köras lokalt. När du har skapat en lösning 
    kan du ladda upp den här för att vara med på topplistan.</p>
    <p>För att lösa problemet behöver du göra följande:</p>
    <p>1. Skicka ett anrop till /new för att få en "token" att identifiera dig med</p>
    <p>2. Med token i headern skicka ett antal anrop till /next (1000 i standardimplementationen).
    Servern kommer för varje anrop att svara med ett tal.<p/>
    <p>3. När alla anrop har genomförts ska de mottagna talen summeras och skickas till /answer tillsammans med ett namn (för topplistans skull)</p>
    </div><br><br>
    <p class="title-text">API</p>`)
    let rows = [
      ["/new", "GET", "", "", `token`],
      ["/next", "GET", "X-Token", "", `number`],
      ["/answer", "POST", "X-Token", `name, sum`, `status, time`]
    ]
    rows.map(v => {
      createRow(v);
    });
    fader.fadeIn(200);
  });
}

function highscoreView(){
  currentView = "highscore";
  $("#Topbtn").attr("class", "active");
  $("#APIbtn").attr("class", "");
  $("#Solutionbtn").attr("class", "");
  $("#Submitbtn").attr("class", "");
  let table = $("#table")
  fader.fadeOut(200, () => {
    clearTable();
    createHeaders(["Namn", "Tid"]);
    $("#topcontent").append('<p class="title-text">TOPPLISTA</p>')
    g_highscores.map(elem => {
      createRow([elem.name, elem.time])
    });
    fader.fadeIn(200);
  });
}

function submitView(){
  currentView = "submit"
  $("#Topbtn").attr("class", "");
  $("#APIbtn").attr("class", "");
  $("#Solutionbtn").attr("class", "");
  $("#Submitbtn").attr("class", "active");
  let table = $("#table")
  fader.fadeOut(200, () => {
    clearTable();
    g_queue.map(elem => {
      createRow([elem.name, elem.qtime, elem.status]);
    })
    $("#topcontent").html(`<div id="form-div">
                            <form id="uploadForm" enctype="multipart/form-data" action="javascript:;" onsubmit="uploadFile()">
                              <div style="display:flex;justify-content:center;align-items:center;flex-direction:column;">
                                <input type="file" name="file" class="inputfile" id="file"/>
                                <label for="file">Välj en fil</label>
                                <input type="submit" value="Ladda upp"/>
                              </div>
                             </form>
                           </div>
                           <p class="title-text">EXEKVERINGSKÖ</p>`)
    createHeaders(["Namn", "Tid i kö", "Status"]);
    var input = document.getElementById("file");
    input.addEventListener('change', (e) => {
      fileName = e.target.value.split( '\\' ).pop();
      if (fileName)
        input.nextElementSibling.innerHTML = fileName;
    })
    
    let dropArea = document.getElementById("form-div");
    ;['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
      dropArea.addEventListener(eventName, preventDefaults, false)
    })
    function preventDefaults (e) {
      e.preventDefault()
      e.stopPropagation()
    }

    ;['dragenter', 'dragover'].forEach(eventName => {
      dropArea.addEventListener(eventName, highlight, false)
    })
    
    ;['dragleave', 'drop'].forEach(eventName => {
      dropArea.addEventListener(eventName, unhighlight, false)
    })
    
    function highlight(e) {
      dropArea.classList.add('highlight')
    }
    
    function unhighlight(e) {
      dropArea.classList.remove('highlight')
    }

    dropArea.addEventListener('drop', handleDrop, false)

    function handleDrop(e) {
      let dt = e.dataTransfer
      let files = dt.files

      document.getElementById("file").files = files;
    }

    fader.fadeIn(200);
  })
}

function solutionView(){
  currentView = "solution";
  $("#Topbtn").attr("class", "");
  $("#APIbtn").attr("class", "");
  $("#Solutionbtn").attr("class", "active");
  $("#Submitbtn").attr("class", "");
  let data = fetch("http://" + window.location.host + "/solutions")
  fader.fadeOut(200, () => {
    clearTable();
    $("#topcontent").append('<p class="title-text">LÖSNINGSFÖRSLAG</p>')
    data
      .then(res => {
        if (res.status == 200)
          return res.json()
        throw "not allowed"
      })
      .then(json => {
        createHeaders(["Namn", "Länk"]);
        json.solutions.map(s => {
          let link = document.createElement("a")
          link.setAttribute("href", s.link);
          link.append("github");
          createRow([s.name, link])
        });
      })
      .catch(err => {
        console.log(err);
        $("#header").append('<div style="display: flex; padding: 5px; align-items:center; justify-content: center;"><img src="images/lock.svg"/></div>')
      })
    fader.fadeIn(200);
  }); 
}


function uploadFile(){
  var form = document.getElementById("uploadForm")
  var formData = new FormData(form);
  $.ajax({
      url:'/upload',
      type:'post',
      data:formData,
      contentType: false,
      processData: false,
      success:() => {
          alert("uppladdat!");
      },
      error:() => {
          alert("ogiltig fil!");
      }
  });
  form.reset();
}
