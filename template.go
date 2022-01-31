package dashboard

import (
	"html/template"
)

type templateInfo struct {
	Projects []*Project
}

var (
	indexTmpl = template.Must(template.New("index.html").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
  <meta content="origin-when-cross-origin" name="referrer" />
  <link crossorigin="anonymous" media="all" integrity="sha512-R+Vpkv86him5JZcqAEuQRUGOKqH897w6q7uJ1P65tQR+9Hxar5vU4wpEd4uvcXT8ooRZ7zsNftrjnCemEt2u2Q==" rel="stylesheet" href="https://github.githubassets.com/assets/frameworks-f4557b27209914aa4705202b188165b5.css" />
  <link crossorigin="anonymous" media="all" integrity="sha512-KY4UdRxVuu2IR/1+5BfVjOAci/NTaEiyg0NelgnXevWPbhUU4YLcpNeWyVh6kVHzLwOF75XT+iW9BwA4zURVaw==" rel="stylesheet" href="https://github.githubassets.com/assets/github-af00f93db15422dc4aa5207a2f2ee52c.css" />
  <title>Dashboard</title>
  <style type="text/css">
  .pagehead {
      margin-bottom: 60px;
      padding: 30px 0;
  }
  .container-sm {
      width: 300px;
  }
  .avatar-container {
      margin-right: 18px;
  }
  .markdown-body {
      width: 95%;
      margin: 0 auto;
  }
  .markdown-body .h2 {
      padding-bottom: 0;
      border-bottom: none;
  }
  .markdown-body table {
      display: table;
      width: auto;
      margin: 0 auto;
  }
  .markdown-body table td {
      text-align: center;
  }
  #jekyll_org th {
      cursor: pointer;
  }
  #jekyll_org th:after {
      margin-left: 9px;
      vertical-align: 1px;
      font-size: 0.6em;
      color: #bfbfbf;
  }
  #jekyll_org th.sorted-ascending:after {
      content: "\25B2";
  }
  #jekyll_org th.sorted-descending:after {
      content: "\25BC";
  }
  .ci-status {
      padding: 6px 0 !important;
  }
  .ci-status a {
      display: block;
      padding: 6px 12px;
  }
  .ci-status a:first-child {
      margin-top: -6px;
  }
  .ci-status a:last-child {
      margin-bottom: -6px;
  }
  .status-good, .success {
      background-color: rgba(0, 255, 0, 0.1);
  }
  .ci-status.success a:hover {
      background-color: rgba(24, 200, 24, 0.1);
  }
  .status-tbd, .pending {
      background-color: rgba(255, 255, 0, 0.1);
  }
  .ci-status.pending a:hover {
      background-color: rgba(210, 210, 180, 0.1);
  }
  .status-bad, .failure {
      background-color: rgba(255, 0, 0, 0.1);
  }
  .ci-status.failure a:hover {
      background-color: rgba(240, 5, 5, 0.1);
  }
  .status-unknown, .error {
      background-color: rgba(0, 0, 0, 0.1);
  }
  .ci-status.error a:hover {
      background-color: rgba(140, 140, 140, 0.1);
  }
  footer {
      margin-top: 30px;
      padding: 15px 0 30px;
      border-top: 1px solid #eee;
  }
  .inline-block {
      display: inline-block;
  }
  .footer-note, .reset-form-container {
      margin: 0 auto;
      padding-bottom: 12px;
  }
  .footer-note {
      max-width: 418px;
  }
  .footer-note div:first-child {
      padding: 12px 15px;
      text-align: center;
      line-height: 1.4;
      font-weight: 600;
      border-right: 1px solid #eee;
  }
  .icon {
      padding: 12px 0 0 12px;
      min-width: 48px;
  }
  .reset-form-container {
      max-width: 250px;
  }
  .reset-form-container [type=submit] {
      margin-left: 4px;
  }
  </style>
  <script type="application/javascript">
  function reqListener () {
    console.log(this.responseText);

    if (this.responseText === null || this.responseText === "") {
        console.log("nada");
        return;
    }

    var info = JSON.parse(this.responseText);
    var tr = document.getElementById(info.name);

    // Name
    var nameTD = document.createElement("td");
    var nameA = document.createElement("a");
    nameA.href = "https://github.com/" + info.nwo;
    nameA.title = info.name + " on GitHub";
    nameA.innerText = info.name;
    nameTD.appendChild(nameA);
    tr.appendChild(nameTD);

    // Gem
    var rubygemsTD = document.createElement("td");
    if (info.gem) {
        var rubygemsA = document.createElement("a");
        rubygemsA.href = "https://rubygems.org/gems/" + info.gem.name;
        rubygemsA.title = info.gem.name + " on RubyGems.org";
        rubygemsA.innerText = "v" + info.gem.version;
        rubygemsTD.appendChild(rubygemsA);
    } else {
        rubygemsTD.innerText = "no info";
    }
    tr.appendChild(rubygemsTD);

    // CI
    var ciTD = document.createElement("td");
    if (info.github.latest_commit_ci_data) {
        ciTD.className = "ci-status"
        let worstCIState = "success";
        for (context of info.github.latest_commit_ci_data) {
            var ciA = document.createElement("a");
            ciA.href = context.url;
            ciA.title = context.name;
            ciA.innerText = ""+context.name+" ("+context.state.toLowerCase()+")";
            ciTD.appendChild(ciA);
            if (worstCIState === "success" && context.state.toLowerCase() !== "neutral") {
                worstCIState = context.state.toLowerCase();
            }
        }
        ciTD.classList.add(worstCIState);
    } else {
        ciTD.innerText = "no info";
    }
    tr.appendChild(ciTD);

    // Downloads
    var downloadsTD = document.createElement("td");
    if (info.gem && info.gem.downloads) {
        downloadsTD.innerText = info.gem.downloads;
    } else {
        downloadsTD.innerText = "no info";
    }
    tr.appendChild(downloadsTD);

    if (info.github === undefined || info.github === null) {
        for (i = 0; i < 4; i++) {
            var emptyTD = document.createElement("td");
            emptyTD.innerText = "no info";
            tr.appendChild(emptyTD);
        }
        return;
    }

    // Pull Requests
    var pullrequestsTD = document.createElement("td");
    if (info.github.open_prs > 0) {
        var pullrequestsA = document.createElement("a");
        pullrequestsA.href = "https://github.com/" + info.nwo + "/pulls";
        pullrequestsA.title = info.name + " pull requests on GitHub";
        pullrequestsA.innerText = info.github.open_prs;
        pullrequestsTD.appendChild(pullrequestsA);
    } else if (info.github.open_prs < 0) {
        pullrequestsTD.innerText = "no info";
    } else {
        pullrequestsTD.innerText = info.github.open_prs;
    }
    tr.appendChild(pullrequestsTD);

    // Issues
    var issuesTD = document.createElement("td");
    if (info.github.open_issues > 0) {
        var issuesA = document.createElement("a");
        issuesA.href = "https://github.com/" + info.nwo + "/issues";
        issuesA.title = info.name + " issues on GitHub";
        issuesA.innerText = info.github.open_issues;
        issuesTD.appendChild(issuesA);
    } else if (info.github.open_issues < 0) {
        issuesTD.innerText = "no info";
    } else {
        issuesTD.innerText = info.github.open_issues;
    }
    tr.appendChild(issuesTD);

    // Unreleased commits
    var unreleasedcommitsTD = document.createElement("td");
    if (info.github.commits_since_latest_release > 0) {
        var unreleasedcommitsA = document.createElement("a");
        unreleasedcommitsA.href = "https://github.com/" + info.nwo + "/compare/" + info.github.latest_release_tag + "...master";
        unreleasedcommitsA.title = info.name + " commits since latest release on GitHub";
        unreleasedcommitsA.innerText = info.github.commits_since_latest_release;
        unreleasedcommitsTD.appendChild(unreleasedcommitsA);
    } else if (info.github.commits_since_latest_release < 0) {
        unreleasedcommitsTD.innerText = "no info";
    } else {
        unreleasedcommitsTD.innerText = info.github.commits_since_latest_release;
    }

    // Start warning us when we get more than 50 commits since the latest release.
    if (info.github.commits_since_latest_release < 10) {
        unreleasedcommitsTD.classList.add("status-good");
    } else if (info.github.commits_since_latest_release < 30) {
        unreleasedcommitsTD.classList.add("status-tbd");
    } else {
        unreleasedcommitsTD.classList.add("status-bad");
    }

    tr.appendChild(unreleasedcommitsTD);
  }

  {{range .Projects}}
  var oReq = new XMLHttpRequest();
  oReq.addEventListener("load", reqListener);
  oReq.open("GET", "/show.json?name={{urlquery .Name}}");
  oReq.send();
  {{end}}
  </script>
</head>
<body>
<div class="markdown-body">
<header class="pagehead">
  <div class="container-sm clearfix">
    <div class="d-flex flex-wrap flex-md-items-center">
      <div class="avatar-container">
        <img itemprop="image" class="avatar" src="https://avatars.githubusercontent.com/u/3083652?s=200&amp;v=4" width="60" height="60" alt="jekyll@github.com">
      </div>
      <div class="flex-1">
        <h1 class="h2 lh-condensed">Dashboard</h1>
        <div><small>Jekyll Organization at a glance</small></div>
      </div>
    </div>
  </div>
</header>
<table id="jekyll_org">
  <thead>
    <tr>
      <th>Repository</th>
      <th>Gem Version</th>
      <th>Continuous Integration</th>
      <th>Downloads</th>
      <th>Pull Requests</th>
      <th>Issues</th>
      <th>Unreleased commits</th>
    </tr>
  </thead>
  <tbody>
    {{range .Projects}}
    <tr id="{{.Name}}"></tr>
    {{end}}
  </tbody>
</table>
<footer>
  <div class="footer-note">
    <div class="inline-block">
      Commits are as of this week.<br>Issues and pull requests are total open.
    </div>
    <div class="inline-block icon">
      <a href="https://github.com/jekyll/dashboard" title="Source Code at GitHub">
        <svg viewBox="0 0 16 16">
          <path
            fill="#828282"
            d="M7.999,0.431c-4.285,0-7.76,3.474-7.76,7.761 c0,3.428,2.223,6.337,5.307,
               7.363c0.388,0.071,0.53-0.168,0.53-0.374c0-0.184-0.007-0.672-0.01-1.32 c-2.159,
               0.469-2.614-1.04-2.614-1.04c-0.353-0.896-0.862-1.135-0.862-1.135c-0.705-0.481,
               0.053-0.472,0.053-0.472 c0.779,0.055,1.189,0.8,1.189,0.8c0.692,1.186,1.816,
               0.843,2.258,0.645c0.071-0.502,0.271-0.843,0.493-1.037 C4.86,11.425,3.049,
               10.76,3.049,7.786c0-0.847,0.302-1.54,0.799-2.082C3.768,5.507,3.501,4.718,
               3.924,3.65 c0,0,0.652-0.209,2.134,0.796C6.677,4.273,7.34,4.187,8,4.184c0.659,
               0.003,1.323,0.089,1.943,0.261 c1.482-1.004,2.132-0.796,2.132-0.796c0.423,1.068,
               0.157,1.857,0.077,2.054c0.497,0.542,0.798,1.235,0.798,2.082 c0,2.981-1.814,
               3.637-3.543,3.829c0.279,0.24,0.527,0.713,0.527,1.437c0,1.037-0.01,1.874-0.01,
               2.129 c0,0.208,0.14,0.449,0.534,0.373c3.081-1.028,5.302-3.935,5.302-7.362C15.76,
               3.906,12.285,0.431,7.999,0.431z"
          />
        </svg>
      </a>
    </div>
  </div>
  <div class="reset-form-container">
    Look wrong?
    <form class="inline-block" action="/reset.json" method="post">
      <input type="Submit" value="Reset the cache.">
    </form>
  </div>
</footer>
<script>
  const table = document.getElementById("jekyll_org");
  const rows = table.rows;
  const header_cells = Array.from(rows[0].cells);
  const numericCells = ["Downloads", "Pull Requests", "Issues", "Unreleased commits"];

  header_cells.forEach((cell, idx) => {
    const type = numericCells.includes(cell.innerText) ? "Numeric" : "String";
    cell.addEventListener("click", function () {
      let switching, i, x, y, shouldSwitch, dir, switchcount = 0;
      switching = true;

      dir = "ascending";

      while (switching) {
        switching = false;

        // Loop through all table rows (except the first, which contains table headers).
        for (i = 1; i < (rows.length - 1); i++) {
          shouldSwitch = false;

          // Get two elements for comparison, one from current row and one from the next.
          x = rows[i].getElementsByTagName("td")[idx];
          y = rows[i + 1].getElementsByTagName("td")[idx];

          // Check if the two rows should switch place, based on the direction, "ascending"
          // or "descending".
          if (dir == "ascending") {
            if (type == "Numeric") {
              if (parseInt(x.innerText) > parseInt(y.innerText)) {
                // If so, mark as a switch and break the loop.
                shouldSwitch = true;
                break;
              }
            } else {
              if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
                // If so, mark as a switch and break the loop.
                shouldSwitch = true;
                break;
              }
            }
          } else if (dir == "descending") {
            if (type == "Numeric") {
              if (parseInt(x.innerText) < parseInt(y.innerText)) {
                // If so, mark as a switch and break the loop.
                shouldSwitch = true;
                break;
              }
            } else {
              if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                // If so, mark as a switch and break the loop.
                shouldSwitch = true;
                break;
              }
            }
          }
        }

        if (shouldSwitch) {
          // If a switch has been marked, make the switch and mark that a switch has been done.
          rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
          switching = true;
          // Each time a switch is done, increase this count by 1.
          switchcount ++;
        } else {
          // If no switching has been done AND the direction is "ascending", set the direction
          // to "descending" and run the while loop again.
          if (switchcount == 0 && dir == "ascending") {
            dir = "descending";
            switching = true;
          }
        }
      }

      // Ensure only the current header cell is assigned designed classname AND reset assigned
      // classname attribute to header-cell based on direction.
      header_cells.forEach(cell => cell.classList.contains("sorted-ascending") ?
        cell.classList.remove("sorted-ascending") : cell.classList.remove("sorted-descending")
      );
      cell.classList.add("sorted-" + dir);
    });
  })
</script>
</body>
</html>
`))
)
