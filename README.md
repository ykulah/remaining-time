# remaining-time
Simple app with GCP Appengine and Datastore.
UI: ```https://frontend-dot-remaining-time-c9dd7.appspot.com```

Deployed on ```https://backend-dot-remaining-time-c9dd7.appspot.com```

### Add User
```
GET /api/addUser/<username>/<first-day-of-job>
```

### Add Trip
```
GET /api/addTrip/<username>/<departure-date>/<arrival-date>
```

### Remove Trip
```
GET /api/removeTrip/<username>/<departure-date>/<arrival-date>
```

### Get List of Trips
```
GET /api/getTrips/<username>
```

### Calculate Remanining Days
```
GET /api/getRemainingDays/<username>
```