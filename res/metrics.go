package res

import (
	"bytes"
	"html/template"
	"time"
)

var RegisterAggregate = `
[
  {
    "$project": {
      "date": {
        "$dateToString": {
          "format": "%Y-%m-%d",
          "date": {
            "$toDate": "$register_date"
          }
        }
      }
    }
  },
  {
    "$group": {
      "_id": "$date",
      "value": {
        "$sum": 1
      }
    }
  },
  {
    "$project": {
      "_id": 0,
      "date": "$_id",
      "value": 1
    }
  }
]
`

func GetLoginStats() string {
	//Get time 5 days ago in %Y-%m-%d format
	filterDate := time.Now().AddDate(0, 0, -5).Format("2006-01-02")

	queryTemplate := `
[
  {
    "$match": {
      "login_date": {
        "$gt": "{{.filterDate}}"
      }
    }
  },
  {
    "$project": {
      "date": {
        "$dateToString": {
          "format": "%Y-%m-%d",
          "date": {
            "$toDate": "$login_date"
          }
        }
      }
    }
  },
  {
    "$group": {
      "_id": "$date",
      "value": {
        "$sum": 1
      }
    }
  },
  {
    "$project": {
      "_id": 0,
      "date": "$_id",
      "value": 1
    }
  }
]
`

	t, _ := template.New("text").Parse(queryTemplate)
	out := bytes.Buffer{}
	t.Execute(&out, map[string]interface{}{"filterDate": filterDate})
	return out.String()
}
