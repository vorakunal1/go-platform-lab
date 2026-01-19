routine.go
    all the routine didn't get execute in time so had to implement waitgroup.
    shared counter didn't update due to race condition - used mutex lock to when a routine/worker is updating the counter no other worker writes it.
channels.go
    consumer will wait until channel is open and producer keeps sending the data.
    if buffered producer won't be able to push the data until channel is empty or data is consumed.
