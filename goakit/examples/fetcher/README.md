# Client Example

This example contains two dependent services: the upstream service fetcher
fetches HTTP documents given their URLs and stores them in a downstream service,
the archiver. 

The example illustrates how to write a service that is a client of another
service using the goa v2 goakit plugin.