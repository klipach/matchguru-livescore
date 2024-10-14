### Write event-driven functions
https://cloud.google.com/functions/docs/writing/write-event-driven-functions

### Create scheduler
https://cloud.google.com/sdk/gcloud/reference/scheduler/jobs/create/pubsub

Cloud Run function is triggered each **n** seconds by Cloud Scheduler job
https://cloud.google.com/scheduler/docs/tut-gcf-pub-sub


### Livescores API doc
https://docs.sportmonks.com/football/endpoints-and-entities/endpoints/livescores/get-all-livescores

### No Livescores events response example
```
HTTP/1.1 200 OK

{
  "message": "No result(s) found matching your request. Either the query did not return any results or you don't have access to it via your current subscription.",
  "subscription": [
    {
      "meta": {
        "trial_ends_at": "2024-10-25 07:51:57",
        "ends_at": "2024-11-11 07:01:10",
        "current_timestamp": 1728901134
      },
      "plans": [
        {
          "plan": "Worldwide Plan",
          "sport": "Football",
          "category": "Advanced"
        }
      ],
      "add_ons": [],
      "widgets": []
    }
  ],
  "rate_limit": {
    "resets_in_seconds": 911,
    "remaining": 2910,
    "requested_entity": "Fixture"
  },
  "timezone": "UTC"
}
```

### Livescore success response
```
{
    "data": [
        {
            "id": 19146701,
            // ...
        }
    ],
    "subscription": [
        // ...
    ],
    "rate_limit": {
        // ...
    },
    "timezone": "UTC"
}
```

