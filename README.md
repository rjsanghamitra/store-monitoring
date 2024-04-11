## Store Monitoring
- This API is for generating reports for stores based on logs.
- The Report contains the uptime and downtime for the last hour, last day and the last week. It is generated in the form of a CSV file.
- First, I converted the given CSV files to an SQLite3 database. Before storing in the database, I sorted the 'polls' database by the store_id first and then in non-increasing order of timestamps,
for easy calculation of uptime and downtime.
- Since the data in the given datasets are recorded in 2023, I have used the last timestamp a poll was taken, as the reference time.
- Since most of the timestamps are at a difference of 1 hour(on an average), I have extrapolated it by considering the hour before that timestamp. In other words, I have truncated the timestamp to the previous hour.
It is mentioned in the GenerateReport() function. Please refer to that for further clarity.
- To calculate the uptime or downtime, I have taken the difference between the current timestamp and the next timestamp(which is earlier than the current one since the values are sorted in non-increasing order). I have also
applied the above extrapolation logic to calculate the uptime/downtime when the current status and the next status are not the same.

I still have to work on a few corner cases, which are giving undesired reports for some stores.

This implementation uses maps for storing the data in the database. It is not space-efficient, but it takes only about a second or two to complete the report generation. But performing all the operations directly on the database is space-efficient but it takes around 29 minutes(I tried that as well) to complete the report generation. Therefore I have not done that implementation.
