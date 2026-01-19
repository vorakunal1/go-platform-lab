context-cancellation.go
    implemented context with DOne() to cancel the context.
channel will wait until consumer consumes all the data from buffer channel no deadlock it's just slow communication because consumers are slow. 
what if context has some process running and cancel signal is sent and what if we have more than 1 context in prod.