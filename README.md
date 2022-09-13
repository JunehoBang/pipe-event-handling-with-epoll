# pipe-event-handling-with-epoll

This git repository is created to test the handling of pipe events by epoll as an experient for development of virtualrouter in TmaxCloud
In this code, two different threads, main and worker, communicate each other using a pipe. 

While the main sends a struct data element using pipe, the worker receives it.
The reception operation is triggered by the epoll's signal notification. 

Please note that this code is for experiment. The boolean channel was applied to observe their properties.


In the advanced virtua
