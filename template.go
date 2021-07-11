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
  .markdown-body {
      width: 95%;
      margin: 0 auto;
  }
  .status-good, .ci-status-success {
      background-color: rgba(0, 255, 0, 0.1);
  }
  .status-tbd, .ci-status-pending {
      background-color: rgba(255, 255, 0, 0.1);
  }
  .status-bad, .ci-status-failure {
      background-color: rgba(255, 0, 0, 0.1);
  }
  .status-unknown, .ci-status-error {
      background-color: rgba(0, 0, 0, 0.1);
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
        let worstCIState = "success";
        for (context of info.github.latest_commit_ci_data) {
            var ciA = document.createElement("a");
            ciA.href = context.url;
            ciA.title = context.name;
            ciA.innerText = ""+context.name+" ("+context.state.toLowerCase()+")";
            ciTD.appendChild(ciA);
            ciTD.appendChild(document.createElement("br"));
            if (worstCIState === "success" && context.state.toLowerCase() !== "neutral") {
                worstCIState = context.state.toLowerCase();
            }
        }
        ciTD.classList.add("ci-status-"+worstCIState);
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

<table>
  <caption>Jekyll At-a-Glance Dashboard</caption>
  <thead>
    <tr>
      <th>Repo</th>
      <th>Gem</th>
      <th>CI</th>
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

<div>
	<strong>Commits are as of this week. Issues and pull requests are total open.</strong>
	<a href="https://github.com/jekyll/dashboard">Source Code</a>.
</div>

<p>
	Look wrong? <form action="/reset.json" method="post"><input type="Submit" value="Reset the cache."></form>
</p>

</div>
</body>
</html>
`))
)
