socket = new WebSocket("ws://" + window.location.host + "/ws")
var highscores = [];
var currentView = "highscore";
highscoreView();

socket.onmessage = event => {
  let h = JSON.parse(event.data).highscores;
  highscores = h;
  if (currentView === "highscore"){
    clearTable();
    createHeaders(["Namn", "Tid"]);
    h.map(elem => {
      createRow([elem.name, elem.time])
    })
  }
}

function clearTable(){
  document.getElementById("header").innerHTML = ""
  document.getElementById("body").innerHTML = ""
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
  let table = $("#table")
  table.fadeOut(200, () => {
    clearTable();
    createHeaders(["Namn", "Metod", "Header", "Body", "Response"]);
    let rows = [
      ["/new", "GET", "", "", `token`],
      ["/next", "GET", "X-Token", "", `number`],
      ["/answer", "POST", "X-Token", `name, sum`, `status, time`]
    ]
    rows.map(v => {
      createRow(v);
    });
    table.fadeIn(200);
  });
}

function highscoreView(){
  currentView = "highscore";
  $("#Topbtn").attr("class", "active");
  $("#APIbtn").attr("class", "");
  $("#Solutionbtn").attr("class", "");
  let table = $("#table")
  table.fadeOut(200, () => {
    clearTable();
    createHeaders(["Namn", "Tid"]);
    highscores.map(elem => {
      createRow([elem.name, elem.time])
    });
    table.fadeIn(200);
  });
}

function solutionView(){
  currentView = "solution";
  $("#Topbtn").attr("class", "");
  $("#APIbtn").attr("class", "");
  $("#Solutionbtn").attr("class", "active");
  let data = fetch("http://" + window.location.host + "/solutions")
  let table = $("#table")
  table.fadeOut(200, () => {
    clearTable();
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
    table.fadeIn(200);
  }); 
}