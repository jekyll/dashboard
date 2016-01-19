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
<html>
<head>
  <meta charset="utf-8">
  <title>Dashboard</title>
</head>
<body>

<table>
  <caption>Jekyll At-a-Glance Dashboard</caption>
  <thead>
    <tr>
      <th>Repo</th>
      <th>Gem</th>
      <th>Travis</th>
      <th>Downloads</th>
      <th>Commits</th>
      <th>Pull Requests</th>
      <th>Issues</th>
      <th>Unreleased commits</th>
    </tr>
  </thead>
  <tbody>
{{range .Projects}}
<tr>
  <td>
	<a href="https://github.com/{{.Nwo}}" title="{{.Name}} on GitHub">{{.Name}}</a>
  </td>
  <td>{{if .Gem}}
    <a href="https://rubygems.org/gems/{{.Gem.Name}}" title="{{.Gem.Name}} homepage">{{.Gem.Name}}</a>
  {{end}}</td>
  <td>{{if .Travis}}
    <a href="https://travis-ci.org/{{.Travis.Nwo}}/builds/{{.Travis.Branch.Id}}">{{.Travis.Branch.State}}</a>
  {{end}}</td>
  <td>{{if .Gem}}{{.Gem.Downloads}}{{end}}</td>
  <td>{{if .GitHub}}{{.GitHub.CommitsThisWeek}}{{end}}</td>
  <td>{{if .GitHub}}{{.GitHub.OpenPRs}}{{end}}</td>
  <td>{{if .GitHub}}{{.GitHub.OpenIssues}}{{end}}</td>
  <td>{{if .GitHub}}{{.GitHub.CommitsSinceLastRelease}}{{end}}</td>
</tr>
{{end}}
  </tbody>
</table>

<div>
	<strong>Commits are as of this week. Issues and pull requests are total open.</strong>
</div>

</body>
</html>
`))
)
