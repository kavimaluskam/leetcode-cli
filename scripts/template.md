---
id: {{.QuestionFrontendID}}
title: "{{.Title}}"
url: "https://leetcode.com/problems/{{.TitleSlug}}/description/"
tags:
{{- range .TopicTags }}
- "{{.Name}}"
{{- end }}
difficulty: "{{.Difficulty}}"
acceptance: "{{.ProblemStats.AcceptRate}}"
total-accepted: "{{.ProblemStats.TotalAcceptedRaw}}"
total-submissions: "{{.ProblemStats.TotalSubmissionRaw}}"
testcase-example: |
  {{.SampleTestCase}}
---

## Problem

{{.Content}}
## Discussion

### Solution

### Complexity Analysis

- Time Complexity:

- Space Complexity:
