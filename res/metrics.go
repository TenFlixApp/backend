package res

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
